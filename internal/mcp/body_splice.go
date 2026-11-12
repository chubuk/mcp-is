package mcp

// body_splice.go — вырезка params.arguments.image_data из JSON-RPC тела запроса
// ДО того, как его прочитает mcp-go (streamable_http.go: io.ReadAll + копия
// RawArguments держат 3–4 копии base64 до конца вызова). Значение потоково
// декодируется в спул-файл, а в тело встаёт маска imageMask — пик ОЗУ не
// зависит ни от размера картинки, ни от числа параллельных загрузок.
//
// Любая неоднозначность (batch-массив, чужой тул, дубль ключа, ошибка base64,
// невалидный JSON) = тело не трогается: запрос уходит в mcp-go байт-в-байт, а
// tool-ошибку вернёт штатный handler (контракт не меняется).

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

// imageTools — тулы, у которых params.arguments.image_data — base64 картинка.
var imageTools = map[string]bool{
	"add_image":        true,
	"create_blog_post": true,
	"update_blog_post": true,
}

// errAbort — «не поняли/не доверяем тело»: запрос уходит в mcp-go без подмены.
var errAbort = errors.New("body splice: abort")

// readCloser сохраняет Close исходного r.Body при подмене Reader'а.
type readCloser struct {
	io.Reader
	io.Closer
}

// bodySplicer сканирует JSON-RPC тело, копируя байты в buf. Значение целевого
// image_data НЕ попадает в buf (на его месте «гэп» позиции gapPos) — сырые
// байты значения (с кавычками и JSON-escape) идут в rawFile, откуда
// декодируются в спул. Восстановление тела = buf[:gap]+rawFile+buf[gap:]+br —
// байт-в-байт на любой точке отказа.
type bodySplicer struct {
	br      *bufio.Reader
	orig    io.Closer
	buf     bytes.Buffer
	gapPos  int // позиция гэпа в buf; <0 — значения ещё не было
	rawFile *os.File
	rawOut  *bufio.Writer
	rawErr  error // первая ошибка записи в rawFile (узнаём на Flush)
	rawLen  int64
	// capturing — значение image_data сейчас вырезается (байты → rawFile).
	capturing bool
	toolName  string
	captures  int
	spoolFile string
	// committed — подменённое тело (значение image_data заменено маской).
	committed []byte
}

// spliceImageBody сканирует тело POST-запроса и, если находит
// params.arguments.image_data у image-тула, декодирует его в спул-файл и
// подменяет тело маленьким JSON с маской. Возвращает путь спула ("" — тело не
// тронуто) и cleanup, убирающий временные файлы; cleanup вызывать defer'ом в
// ServeHTTP — спул-файл должен жить до конца вызова (InSales-загрузка).
func spliceImageBody(r *http.Request) (spoolPath string, cleanup func()) {
	if r == nil || r.Body == nil || r.Method != http.MethodPost {
		return "", func() {}
	}
	s := &bodySplicer{br: bufio.NewReader(r.Body), orig: r.Body, gapPos: -1}
	scanErr := s.run()
	s.flushRaw()
	if scanErr == nil {
		scanErr = s.rawErr
	}
	if scanErr == nil && s.captures == 1 && imageTools[s.toolName] {
		path, n, sum, err := s.decodeValue()
		if err == nil {
			s.spoolFile = path
			if err := s.commit(imageMask(n, sum)); err == nil {
				s.applyBody(r)
				return path, s.closeAll
			}
		}
		// Значение не раскодировалось/не подставилось — тело уйдёт как есть,
		// «некорректный base64» вернёт штатный handler.
	}
	s.restore(r)
	return "", s.closeAll
}

// run — основной цикл сканера. errAbort/ошибка = тело восстановит restore().
func (s *bodySplicer) run() error {
	// Корень — только JSON-объект; batch-массив (и всё прочее) не трогаем.
	c, err := s.nextNonSpace()
	if err != nil {
		return err
	}
	if c != '{' {
		return errAbort
	}
	stack := []frame{{isObject: true, expectKey: true}}
	for {
		c, err := s.nextNonSpace()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return err
		}
		switch c {
		case '{':
			stack = append(stack, frame{isObject: true, expectKey: true})
		case '[':
			stack = append(stack, frame{})
		case '}', ']':
			if len(stack) == 0 {
				return errAbort
			}
			stack = stack[:len(stack)-1]
		case ',':
			if len(stack) > 0 && stack[len(stack)-1].isObject {
				top := &stack[len(stack)-1]
				top.expectKey = true
				top.key = ""
			}
		case ':':
			// значение придёт следующим токеном
		case '"':
			if len(stack) == 0 {
				return errAbort
			}
			top := &stack[len(stack)-1]
			if top.isObject && top.expectKey {
				key, err := s.readString()
				if err != nil {
					return err
				}
				top.key = key
				top.expectKey = false
				continue
			}
			if isImageValue(stack) {
				if s.gapPos >= 0 {
					return errAbort // дубль image_data — тело не трогаем
				}
				if err := s.beginCapture(); err != nil {
					return err
				}
				// Открывающая кавычка (последний байт в buf) — часть значения.
				s.buf.Truncate(s.buf.Len() - 1)
				s.gapPos = s.buf.Len()
				if _, err := s.readString(); err != nil {
					return err
				}
				s.endCapture()
				s.captures++
				continue
			}
			small, err := s.readString()
			if err != nil {
				return err
			}
			if isNameValue(stack) {
				s.toolName = small
				if !imageTools[small] {
					return errAbort // чужой тул — дальше не сканируем
				}
			}
		default:
			if err := s.readLiteral(); err != nil {
				return err
			}
		}
	}
	if len(stack) != 0 {
		return errAbort // обрезанное/невалидное тело
	}
	return nil
}

// frame — уровень вложенности JSON и текущий ключ объекта.
type frame struct {
	isObject  bool
	expectKey bool
	key       string
}

// isImageValue — значение это params.arguments.image_data (ровно этот путь).
func isImageValue(stack []frame) bool {
	return len(stack) == 3 &&
		stack[0].isObject && stack[0].key == "params" &&
		stack[1].isObject && stack[1].key == "arguments" &&
		stack[2].isObject && stack[2].key == "image_data"
}

// isNameValue — значение это params.name (имя вызываемого тула).
func isNameValue(stack []frame) bool {
	return len(stack) == 2 &&
		stack[0].isObject && stack[0].key == "params" &&
		stack[1].isObject && stack[1].key == "name"
}

// nextByte читает байт и сохраняет его: в buf или (при capture) в rawFile.
func (s *bodySplicer) nextByte() (byte, error) {
	c, err := s.br.ReadByte()
	if err != nil {
		return 0, err
	}
	if s.capturing {
		if err := s.rawOut.WriteByte(c); err != nil && s.rawErr == nil {
			s.rawErr = err
		}
		s.rawLen++
	} else {
		s.buf.WriteByte(c)
	}
	return c, nil
}

func (s *bodySplicer) nextNonSpace() (byte, error) {
	for {
		c, err := s.nextByte()
		if err != nil {
			return 0, err
		}
		switch c {
		case ' ', '\t', '\n', '\r':
			continue
		}
		return c, nil
	}
}

func (s *bodySplicer) peekByte() (byte, error) {
	p, err := s.br.Peek(1)
	if err != nil {
		return 0, err
	}
	return p[0], nil
}

// readString дочитывает JSON-строку (открывающая кавычка уже съедена).
// Возвращается раскодированное содержимое только для коротких строк
// (ключи, params.name) — до 256 байт; длинные значения не материализуются.
func (s *bodySplicer) readString() (string, error) {
	var small []byte
	smallOK := true
	for {
		c, err := s.nextByte()
		if err != nil {
			return "", err
		}
		if c == '"' {
			return string(small), nil
		}
		if c == '\\' {
			e, err := s.nextByte()
			if err != nil {
				return "", err
			}
			if !smallOK {
				if e == 'u' {
					for i := 0; i < 4; i++ {
						if _, err := s.nextByte(); err != nil {
							return "", err
						}
					}
				}
				continue
			}
			switch e {
			case '"', '\\', '/':
				small = append(small, e)
			case 'b':
				small = append(small, '\b')
			case 'f':
				small = append(small, '\f')
			case 'n':
				small = append(small, '\n')
			case 'r':
				small = append(small, '\r')
			case 't':
				small = append(small, '\t')
			case 'u':
				var hexs [4]byte
				for i := range hexs {
					if hexs[i], err = s.nextByte(); err != nil {
						return "", err
					}
				}
				n, err := strconv.ParseUint(string(hexs[:]), 16, 16)
				if err != nil {
					return "", errAbort
				}
				small = append(small, []byte(string(rune(n)))...)
			default:
				return "", errAbort
			}
			continue
		}
		if smallOK {
			small = append(small, c)
			if len(small) > 256 {
				smallOK = false
				small = nil
			}
		}
	}
}

// readLiteral дочитывает литерал (число/true/false/null): первый байт съеден.
func (s *bodySplicer) readLiteral() error {
	for {
		c, err := s.peekByte()
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return err
		}
		switch c {
		case ',', '}', ']', ' ', '\t', '\n', '\r':
			return nil
		}
		if _, err := s.nextByte(); err != nil {
			return err
		}
	}
}

// beginCapture переключает поток байт в rawFile (значение image_data).
// Порядок важен: файл создаётся ДО вырезания кавычки из buf, иначе ошибка
// CreateTemp испортит восстановление тела.
func (s *bodySplicer) beginCapture() error {
	f, err := os.CreateTemp("", spoolPrefix+"*")
	if err != nil {
		return fmt.Errorf("image_data: временный файл: %w", err)
	}
	s.rawFile = f
	s.rawOut = bufio.NewWriter(f)
	s.capturing = true
	// Открывающая кавычка — уже вне buf (см. run).
	s.rawLen = 1
	s.rawErr = s.rawOut.WriteByte('"')
	return nil
}

func (s *bodySplicer) endCapture() {
	s.capturing = false
}

func (s *bodySplicer) flushRaw() {
	if s.rawOut == nil {
		return
	}
	if err := s.rawOut.Flush(); err != nil && s.rawErr == nil {
		s.rawErr = err
	}
}

// decodeValue декодирует вырезанное значение (без кавычек, с JSON-escape)
// base64 → спул-файл. Путь возвращается уже без data:-обёртки.
func (s *bodySplicer) decodeValue() (path string, n int64, sum []byte, err error) {
	if s.rawFile == nil || s.rawLen < 2 {
		return "", 0, nil, errAbort
	}
	if _, err := s.rawFile.Seek(0, io.SeekStart); err != nil {
		return "", 0, nil, err
	}
	sec := io.NewSectionReader(s.rawFile, 1, s.rawLen-2) // без кавычек
	unescaped := &jsonUnescapeReader{r: bufio.NewReader(sec)}
	return decodeImageStream(stripDataURL(unescaped))
}

// commit подготавливает маленькое тело: значение image_data → маска.
func (s *bodySplicer) commit(mask string) error {
	maskJSON, err := json.Marshal(mask)
	if err != nil {
		return err
	}
	b := s.buf.Bytes()
	out := make([]byte, 0, len(b)+len(maskJSON))
	out = append(out, b[:s.gapPos]...)
	out = append(out, maskJSON...)
	out = append(out, b[s.gapPos:]...)
	s.committed = out
	return nil
}

// applyBody ставит подменённое тело в запрос.
func (s *bodySplicer) applyBody(r *http.Request) {
	r.Body = io.NopCloser(bytes.NewReader(s.committed))
	r.ContentLength = int64(len(s.committed))
}

// restore возвращает тело байт-в-байт: сырые байты вырезанного значения
// встают на своё место (rawFile входит в цепочку до конца вызова — cleanup
// закроет его после чтения тела mcp-go).
func (s *bodySplicer) restore(r *http.Request) {
	var parts []io.Reader
	b := s.buf.Bytes()
	if s.gapPos < 0 {
		parts = append(parts, bytes.NewReader(b))
	} else {
		parts = append(parts, bytes.NewReader(b[:s.gapPos]))
		if s.rawFile != nil {
			_, _ = s.rawFile.Seek(0, io.SeekStart)
			parts = append(parts, io.LimitReader(s.rawFile, s.rawLen))
		}
		parts = append(parts, bytes.NewReader(b[s.gapPos:]))
	}
	parts = append(parts, s.br)
	r.Body = readCloser{io.MultiReader(parts...), s.orig}
}

// closeAll убирает временные файлы (вызывать после завершения запроса).
func (s *bodySplicer) closeAll() {
	if s.rawFile != nil {
		name := s.rawFile.Name()
		_ = s.rawFile.Close()
		_ = os.Remove(name)
	}
	if s.spoolFile != "" {
		_ = os.Remove(s.spoolFile)
	}
}

// stripDataURL снимает обёртку data:…;base64, из начала потока
// (эквивалент strings.HasPrefix/Index в spoolImageData).
func stripDataURL(r io.Reader) io.Reader {
	br := bufio.NewReader(r)
	if head, _ := br.Peek(len("data:")); string(head) != "data:" {
		return br
	}
	// data:-префикс короткий: base64, гарантированно в первых 512 байтах.
	chunk, _ := br.Peek(512)
	if i := bytes.Index(chunk, []byte("base64,")); i >= 0 {
		_, _ = br.Discard(i + len("base64,"))
	}
	return br
}

// jsonUnescapeReader раскрывает JSON-escape строки (\n, \/, \uXXXX):
// значение image_data в теле лежит экранированным, а base64-декодеру нужны
// исходные байты. Содержимое строки ограничено её кавычками — их секция
// не передаётся (SectionReader без краёв). Быстрый путь (base64 без escape):
// байты идут из входа в выход через bytes.IndexByte, без посимвольных копий.
type jsonUnescapeReader struct {
	r    *bufio.Reader
	pend []byte // сырой хвост: escape в начале или кусок до следующего escape
	esc  [4]byte
	escN int // готовых к отдаче байт в esc
}

func (u *jsonUnescapeReader) Read(p []byte) (int, error) {
	for {
		if u.escN > 0 {
			n := copy(p, u.esc[:u.escN])
			copy(u.esc[:], u.esc[n:u.escN])
			u.escN -= n
			return n, nil
		}
		if len(u.pend) > 0 {
			if u.pend[0] == '\\' {
				if err := u.unescape(); err != nil {
					return 0, err
				}
				continue
			}
			i := bytes.IndexByte(u.pend, '\\')
			if i < 0 {
				i = len(u.pend)
			}
			n := copy(p, u.pend[:i])
			u.pend = u.pend[n:]
			return n, nil
		}
		n, err := u.r.Read(p)
		if n == 0 {
			return 0, err
		}
		i := bytes.IndexByte(p[:n], '\\')
		if i < 0 {
			return n, err
		}
		// Нашли escape: кусок до него уже в p (отдаём), хвост — в pend.
		u.pend = append(u.pend[:0], p[i:n]...)
		if i > 0 {
			return i, nil
		}
	}
}

// unescape раскрывает один escape в начале pend в u.esc.
func (u *jsonUnescapeReader) unescape() error {
	if len(u.pend) < 2 {
		return fmt.Errorf("image_data: обрыв escape")
	}
	e := u.pend[1]
	u.pend = u.pend[2:]
	switch e {
	case '"', '\\', '/':
		u.esc[0], u.escN = e, 1
	case 'b':
		u.esc[0], u.escN = '\b', 1
	case 'f':
		u.esc[0], u.escN = '\f', 1
	case 'n':
		u.esc[0], u.escN = '\n', 1
	case 'r':
		u.esc[0], u.escN = '\r', 1
	case 't':
		u.esc[0], u.escN = '\t', 1
	case 'u':
		if len(u.pend) < 4 {
			return fmt.Errorf("image_data: обрыв \\u")
		}
		n, err := strconv.ParseUint(string(u.pend[:4]), 16, 16)
		if err != nil {
			return fmt.Errorf("image_data: некорректный \\u: %w", err)
		}
		u.pend = u.pend[4:]
		u.escN = copy(u.esc[:], []byte(string(rune(n))))
	default:
		return fmt.Errorf("image_data: неизвестный escape \\%c", e)
	}
	return nil
}
