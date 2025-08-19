package model

import (
	"bufio"
	"encoding/csv"
	"errors"
	"os"
	"strings"
)

// Item ist ein vereinheitlichtes Modell für das PDF.
type Item struct {
	Title    string
	Category string
	Vault    string
	Username string
	Password string
	URLs     []string
	Notes    string
	TOTP     string
	// RawFields enthält alle Label->Value-Paare, die nicht in die Standardfelder fielen.
	RawFields map[string]string
}

// FromCSV parst eine 1Password-CSV (Logins). Spalten können je nach Export variieren.
func FromCSV(path, delimiter string) ([]Item, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	sep := ','
	if delimiter != "" {
		sepRunes := []rune(delimiter)
		sep = sepRunes[0]
	}
	r := csv.NewReader(bufio.NewReader(f))
	r.Comma = sep
	r.FieldsPerRecord = -1

	rows, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, errors.New("leere CSV")
	}

	// Header heuristisch
	header := rows[0]
	idx := func(name string) int {
		lname := strings.ToLower(name)
		for i, h := range header {
			if strings.ToLower(strings.TrimSpace(h)) == lname {
				return i
			}
		}
		return -1
	}

	iTitle := idx("title")
	if iTitle == -1 { iTitle = idx("name") }
	iUser := idx("username")
	iPass := idx("password")
	iURL := idx("url")
	iNotes := idx("notes")

	var items []Item
	for _, row := range rows[1:] {
		get := func(i int) string {
			if i >= 0 && i < len(row) { return row[i] }
			return ""
		}
		it := Item{
			Title:    get(iTitle),
			Category: "login",
			Vault:    "",
			Username: get(iUser),
			Password: get(iPass),
			URLs:     nil,
			Notes:    get(iNotes),
			RawFields: map[string]string{},
		}
		if iURL >= 0 {
			u := strings.TrimSpace(get(iURL))
			if u != "" {
				it.URLs = []string{u}
			}
		}
		items = append(items, it)
	}
	return items, nil
}
