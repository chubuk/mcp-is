package core

import (
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// Slugify converts a human-readable name into a URL/JSON-safe slug.
// Rules: lowercase, latin + digits + hyphens only, no leading/trailing hyphens.
// Cyrillic is transliterated. Empty input returns "mcp-is".
func Slugify(name string) string {
	if strings.TrimSpace(name) == "" {
		return "mcp-is"
	}

	// Transliterate Cyrillic FIRST (before NFD decomposes й→и+breve)
	lower := strings.ToLower(name)
	normalized := transliterate(lower)

	// Normalize + strip diacritics (ü→u, etc.)
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	normalized, _, _ = transform.String(t, normalized)

	// Keep only latin, digits, hyphens
	var b strings.Builder
	prevHyphen := false
	for _, r := range normalized {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
			prevHyphen = false
		case r == ' ' || r == '_' || r == '-':
			if !prevHyphen && b.Len() > 0 {
				b.WriteRune('-')
				prevHyphen = true
			}
		}
	}

	result := strings.TrimRight(b.String(), "-")
	if result == "" {
		return "mcp-is"
	}
	return result
}

var cyrillicMap = map[rune]string{
	'а': "a", 'б': "b", 'в': "v", 'г': "g", 'д': "d",
	'е': "e", 'ё': "e", 'ж': "zh", 'з': "z", 'и': "i",
	'й': "y", 'к': "k", 'л': "l", 'м': "m", 'н': "n",
	'о': "o", 'п': "p", 'р': "r", 'с': "s", 'т': "t",
	'у': "u", 'ф': "f", 'х': "kh", 'ц': "ts", 'ч': "ch",
	'ш': "sh", 'щ': "shch", 'ъ': "", 'ы': "y", 'ь': "",
	'э': "e", 'ю': "yu", 'я': "ya",
}

func transliterate(s string) string {
	var b strings.Builder
	for _, r := range s {
		if latin, ok := cyrillicMap[r]; ok {
			b.WriteString(latin)
		} else {
			b.WriteRune(r)
		}
	}
	return b.String()
}
