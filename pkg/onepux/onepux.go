package onepux

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/example/onepw-pdf-export/pkg/model"
)

// FromFile parst eine .1pux ZIP best-effort.
// Der Export kann zwischen Versionen variieren; wir suchen nach plausiblen JSON-Dateien
// mit Item-ähnlicher Struktur. Nicht alle 1PUX-Varianten werden unterstützt.
func FromFile(path string) ([]model.Item, error) {
	f, err := os.Open(path)
	if err != nil { return nil, err }
	defer f.Close()

	st, err := f.Stat()
	if err != nil { return nil, err }

	zr, err := zip.NewReader(f, st.Size())
	if err != nil {
		return nil, errors.New("konnte .1pux nicht als ZIP lesen: " + err.Error())
	}

	var items []model.Item

	for _, zf := range zr.File {
		low := strings.ToLower(zf.Name)
		// heuristisch: JSON-Dateien unterhalb von vaults/ oder data/-Ordnern
		if !(strings.HasSuffix(low, ".json") && (strings.Contains(low, "vault") || strings.Contains(low, "item") || strings.Contains(low, "data"))) {
			continue
		}
		rc, err := zf.Open()
		if err != nil { continue }
		b, _ := io.ReadAll(rc)
		_ = rc.Close()

		// Versuche mehrere Strukturen
		var generic []map[string]interface{}
		if json.Unmarshal(b, &generic) == nil {
			for _, g := range generic {
				if it, ok := mapGeneric(g); ok {
					items = append(items, it)
				}
			}
			continue
		}
		var single map[string]interface{}
		if json.Unmarshal(b, &single) == nil {
			if it, ok := mapGeneric(single); ok {
				items = append(items, it)
			}
			continue
		}
	}

	if len(items) == 0 {
		return nil, errors.New("keine unterstützten Items in 1PUX gefunden (experimentell)")
	}
	return items, nil
}

func mapGeneric(g map[string]interface{}) (model.Item, bool) {
	// sehr grobe Heuristik
	getS := func(keys ...string) string {
		for _, k := range keys {
			if v, ok := g[k]; ok {
				if s, ok2 := v.(string); ok2 && strings.TrimSpace(s) != "" {
					return s
				}
			}
		}
		return ""
	}
	title := getS("title", "name")
	if title == "" { return model.Item{}, false }

	it := model.Item{
		Title:    title,
		Vault:    getS("vault", "vaultName"),
		Category: getS("category", "type"),
		URLs:     []string{},
		RawFields: map[string]string{},
	}
	// häufige Felder
	user := getS("username", "user", "login", "loginUsername")
	pass := getS("password", "loginPassword")
	notes := getS("notes", "notesPlain", "note")
	totp := getS("totp", "otp", "oneTimePassword")

	it.Username = user
	it.Password = pass
	it.Notes = notes
	it.TOTP = totp

	// URLs (einzelne Felder)
	u := getS("url", "website")
	if u != "" { it.URLs = append(it.URLs, u) }

	// zusätzliche Felder pauschal übernehmen (nur strings)
	for k, v := range g {
		if _, ok := v.(string); ok {
			ks := strings.ToLower(k)
			if ks == "title" || ks == "name" || ks == "username" || ks == "password" || ks == "notes" || ks == "notesplain" || ks == "totp" || ks == "otp" || ks == "onetimepassword" || ks == "url" || ks == "website" {
				continue
			}
			it.RawFields[k] = v.(string)
		}
	}

	return it, true
}
