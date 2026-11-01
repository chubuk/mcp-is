package pdn

import (
	"encoding/csv"
	"strings"
)

// MaskCSV маскирует PII-ячейки в CSV по именам колонок (152-ФЗ).
// Структура/квотирование сохраняются; нераспознанный CSV возвращается как есть.
func MaskCSV(text string) string {
	r := csv.NewReader(strings.NewReader(text))
	rows, err := r.ReadAll()
	if err != nil || len(rows) < 2 {
		return text // без строк данных маскировать нечего; однорядный «CSV» — не CSV
	}
	w := &strings.Builder{}
	wr := csv.NewWriter(w)
	colMask := map[int]func(string) string{}
	for i, h := range rows[0] {
		switch strings.ToLower(strings.TrimSpace(h)) {
		case "email":
			colMask[i] = MaskEmail
		case "phone", "formatted_phone":
			colMask[i] = MaskPhone
		case "first_name", "last_name", "name", "surname", "middlename", "full_name", "client_name":
			colMask[i] = MaskName
		}
	}
	for _, row := range rows[1:] {
		for i, fn := range colMask {
			if i < len(row) {
				row[i] = fn(row[i])
			}
		}
	}
	if err := wr.WriteAll(rows); err != nil {
		return text
	}
	return w.String()
}
