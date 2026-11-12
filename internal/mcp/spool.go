package mcp

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode"
)

// maxImageBytes — лимит раскодированной картинки (20 МБ: потолок InSales ~22 МБ,
// дальше их nginx отдаёт 413 — research/probe_article_image.md).
const maxImageBytes = 20 << 20

// spoolPrefix — маска временных файлов картинок (spoolImageData, spliceImageBody).
const spoolPrefix = "mcp-is-img-"

// imageMask — маска значения image_data в args и в теле запроса. Формат общий для
// спула в handler'е и вырезки на HTTP-уровне (body_splice.go): от него зависят
// diffFromArgs (не двоит маску) и устойчивость хэша аудита/идемпотентности.
func imageMask(n int64, sum []byte) string {
	return fmt.Sprintf("<base64, %d КБ, sha256:%s>",
		base64.StdEncoding.EncodedLen(int(n))/1024, hex.EncodeToString(sum)[:12])
}

// spoolPathFromContext — путь спул-файла, вырезанного из тела запроса на
// HTTP-уровне (spliceImageBody). Пусто — значения не было или тело не тронуто.
func spoolPathFromContext(ctx context.Context) string {
	p, _ := ctx.Value(mcpSpoolPathKey).(string)
	return p
}

// decodeImageStream декодирует base64-поток r во временный файл mcp-is-img-*.
// Потоково, без полных копий в памяти: байты уходят на диск и оттуда стримятся
// в InSales (insales.spoolBody). data:-обёртку снимает вызывающий. При ошибке
// временный файл удаляется до возврата.
func decodeImageStream(r io.Reader) (path string, n int64, sum []byte, err error) {
	f, err := os.CreateTemp("", spoolPrefix+"*")
	if err != nil {
		return "", 0, nil, fmt.Errorf("image_data: временный файл: %w", err)
	}
	h := sha256.New()
	dec := base64.NewDecoder(base64.StdEncoding, &spaceSkipper{r: r})
	// +1 байт над лимитом: отличить «ровно 20 МБ» от «больше лимита».
	n, err = io.Copy(io.MultiWriter(f, h), io.LimitReader(dec, maxImageBytes+1))
	if cerr := f.Close(); cerr != nil && err == nil {
		err = cerr
	}
	if err != nil {
		_ = os.Remove(f.Name())
		var corr base64.CorruptInputError
		if errors.As(err, &corr) {
			return "", 0, nil, fmt.Errorf("image_data: некорректный base64: %w", err)
		}
		return "", 0, nil, fmt.Errorf("image_data: запись временного файла: %w", err)
	}
	if n == 0 {
		_ = os.Remove(f.Name())
		return "", 0, nil, fmt.Errorf("image_data: пустое изображение")
	}
	if n > maxImageBytes {
		_ = os.Remove(f.Name())
		return "", 0, nil, fmt.Errorf("изображение слишком большое: максимум %d МБ", maxImageBytes>>20)
	}
	return f.Name(), n, h.Sum(nil), nil
}

// spoolImageData декодирует args["image_data"] (base64 или data:URI) во временный
// файл и подменяет значение маской imageMask. Потоково, без полных копий строки
// в памяти. В args["image_file"] кладётся путь — его обязан вырезать маппер
// пэлоада. cleanup удаляет файл (вызывать defer'ом).
// Если тело запроса уже распаковано в спул на HTTP-уровне (spliceImageBody),
// декодировать нечего: в image_data уже маска, остаётся подставить путь из ctx.
func spoolImageData(ctx context.Context, args map[string]any) (cleanup func(), err error) {
	noop := func() {}
	if p := spoolPathFromContext(ctx); p != "" {
		args["image_file"] = p
		return noop, nil
	}
	raw, _ := args["image_data"].(string)
	if strings.HasPrefix(raw, "data:") {
		if i := strings.Index(raw, "base64,"); i >= 0 {
			raw = raw[i+len("base64,"):]
		}
	}
	path, n, sum, err := decodeImageStream(strings.NewReader(raw))
	if err != nil {
		return noop, err
	}
	args["image_file"] = path
	// Маска вместо мегабайт base64: diff/аудит/хэши не тащат сырые байты.
	// sha256 содержимого в маске — иначе разные картинки одного размера дают
	// одинаковый ArgsHash аудита.
	args["image_data"] = imageMask(n, sum)
	return func() { _ = os.Remove(path) }, nil
}

// spaceSkipper пропускает пробельные символы в base64-потоке (эквивалент
// strings.Fields до декода: переносы строк в PEM-подобном вводе).
type spaceSkipper struct {
	r   io.Reader
	buf []byte
}

func (s *spaceSkipper) Read(p []byte) (int, error) {
	if s.buf == nil {
		s.buf = make([]byte, len(p))
	}
	n, err := s.r.Read(s.buf)
	out := p[:0]
	for _, c := range s.buf[:n] {
		if !unicode.IsSpace(rune(c)) {
			out = append(out, c)
		}
	}
	return len(out), err
}

// SweepSpoolFiles удаляет забытые спул-файлы картинок старше часа
// (мусор на диске после сбоя процесса между спулом и чисткой).
func SweepSpoolFiles() {
	matches, err := filepath.Glob(filepath.Join(os.TempDir(), spoolPrefix+"*"))
	if err != nil {
		return
	}
	for _, m := range matches {
		if fi, err := os.Stat(m); err == nil && time.Since(fi.ModTime()) > time.Hour {
			_ = os.Remove(m)
		}
	}
}
