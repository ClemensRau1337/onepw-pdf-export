package op

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/example/onepw-pdf-export/pkg/model"
)

type Vault struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// opItemListEntry spiegelt die grobe Struktur von `op item list --format json`.
type opItemListEntry struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Category string `json:"category"`
	Vault    struct {
		Name string `json:"name"`
	} `json:"vault"`
}

// opItemDetail spiegelt `op item get <id> --format json` minimal.
type opItemDetail struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Category string `json:"category"`
	Vault    struct {
		Name string `json:"name"`
	} `json:"vault"`
	Fields []struct {
		ID    string      `json:"id"`
		Label string      `json:"label"`
		Type  string      `json:"type"` // e.g., CONCEALED for password
		Value interface{} `json:"value"`
	} `json:"fields"`
	URLs []struct {
		Label string `json:"label"`
		Href  string `json:"href"`
	} `json:"urls"`
	NotesPlain string `json:"notesPlain"`
}

// ListVaults ruft alle Tresore ab.
func ListVaults() ([]Vault, error) {
	raw, err := runOpJSON("vault", "list")
	if err != nil {
		return nil, err
	}
	var v []Vault
	if err := json.Unmarshal(raw, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// ListItems ruft eine lightweight-Liste aller Items ab.
func ListItems() ([]opItemListEntry, error) {
	listRaw, err := runOpJSON("item", "list")
	if err != nil {
		return nil, err
	}
	var list []opItemListEntry
	if err := json.Unmarshal(listRaw, &list); err != nil {
		return nil, err
	}
	return list, nil
}

// GetItemDetails holt Detaildaten für eine Item-ID.
func GetItemDetails(id string) (model.Item, error) {
	detailRaw, err := runOpJSON("item", "get", id)
	if err != nil {
		return model.Item{}, err
	}
	var d opItemDetail
	if err := json.Unmarshal(detailRaw, &d); err != nil {
		return model.Item{}, err
	}
	return mapDetail(d), nil
}

// FetchAllItems (legacy helper) lädt alle Items vollständig.
func FetchAllItems() ([]model.Item, error) {
	list, err := ListItems()
	if err != nil {
		return nil, err
	}
	out := make([]model.Item, 0, len(list))
	for _, e := range list {
		it, err := GetItemDetails(e.ID)
		if err != nil {
			continue
		}
		out = append(out, it)
	}
	return out, nil
}

func mapDetail(d opItemDetail) model.Item {
	it := model.Item{
		Title:    d.Title,
		Category: d.Category,
		Vault:    d.Vault.Name,
		URLs:     []string{},
		RawFields: map[string]string{},
	}
	for _, u := range d.URLs {
		if strings.TrimSpace(u.Href) != "" {
			it.URLs = append(it.URLs, u.Href)
		}
	}
	if strings.TrimSpace(d.NotesPlain) != "" {
		it.Notes = d.NotesPlain
	}

	for _, f := range d.Fields {
		val := stringify(f.Value)
		lbl := strings.TrimSpace(f.Label)
		typ := strings.ToUpper(strings.TrimSpace(f.Type))

		// heuristisch
		switch {
		case strings.EqualFold(lbl, "username") || strings.Contains(strings.ToLower(lbl), "user"):
			if it.Username == "" { it.Username = val }
		case strings.EqualFold(lbl, "password") || typ == "CONCEALED":
			if it.Password == "" { it.Password = val }
		case strings.Contains(strings.ToLower(lbl), "otp") || strings.Contains(strings.ToLower(lbl), "totp"):
			if it.TOTP == "" { it.TOTP = val }
		default:
			if lbl != "" && val != "" {
				it.RawFields[lbl] = val
			}
		}
	}
	return it
}

func stringify(v interface{}) string {
	switch t := v.(type) {
	case string:
		return t
	case float64:
		return fmt.Sprintf("%v", t)
	default:
		b, _ := json.Marshal(t)
		return string(b)
	}
}

func runOpJSON(args ...string) ([]byte, error) {
	cmd := exec.Command("op", append(args, "--format", "json")...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		msg := stderr.String()
		if msg == "" { msg = err.Error() }
		return nil, errors.New("op: " + msg)
	}
	return stdout.Bytes(), nil
}
