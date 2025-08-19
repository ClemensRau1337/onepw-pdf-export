package pdfwriter

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/jung-kurt/gofpdf"
	"github.com/example/onepw-pdf-export/pkg/fonts"
	"github.com/example/onepw-pdf-export/pkg/model"
)

type Options struct {
	Template     string // compact | detailed
	MaskPassword bool
	Source       string // csv | live/op | 1pux
	UserPassword string // PDF user password (required)
}

func randomOwnerPassword() string {
	var b [16]byte
	_, _ = rand.Read(b[:])
	return hex.EncodeToString(b[:])
}

func mask(s string, on bool) string {
	if !on { return s }
	if s == "" { return "" }
	return strings.Repeat("•", 8)
}

func WritePDF(path string, items []model.Item, opt Options) error {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetTitle("1Password Export", false)
	pdf.SetAuthor("onepw-pdf-export", false)

	// UTF-8 capable font
	if err := fonts.EnsureUTF8Font(pdf); err == nil {
		pdf.SetFont(fonts.FontName, "", 12)
	} else {
		// Fallback (no full UTF‑8)
		pdf.SetFont("Helvetica", "", 12)
	}

	// Protection
	if opt.UserPassword == "" {
		return fmt.Errorf("pdfwriter: UserPassword ist leer")
	}
	pdf.SetProtection(gofpdf.CnProtectPrint, opt.UserPassword, randomOwnerPassword())

	pdf.AddPage()

	// Header
	pdf.Cell(0, 6, fmt.Sprintf("Exportzeitpunkt: %s", time.Now().Format(time.RFC3339)))
	pdf.Ln(6)
	pdf.Cell(0, 6, fmt.Sprintf("Quelle: %s | Items: %d", opt.Source, len(items)))
	pdf.Ln(10)

	for _, it := range items {
		writeItem(pdf, it, opt)
	}

	return pdf.OutputFileAndClose(path)
}

func writeItem(pdf *gofpdf.Fpdf, it model.Item, opt Options) {
	w, h := 190.0, 6.0

	// Titelzeile
	pdf.SetFontStyle("B")
	title := it.Title
	if title == "" { title = "(ohne Titel)" }
	meta := it.Vault
	if it.Category != "" {
		if meta != "" { meta += " · " }
		meta += it.Category
	}
	pdf.CellFormat(0, 7, title, "", 1, "", false, 0, "")
	if meta != "" {
		pdf.SetFontStyle("")
		pdf.SetFontSize(10)
		pdf.CellFormat(0, 5, meta, "", 1, "", false, 0, "")
		pdf.SetFontSize(11)
	}

	// Inhalt
	pdf.SetFontStyle("")

	kv := func(k, v string) {
		if v == "" { return }
		pdf.CellFormat(30, h, k, "", 0, "", false, 0, "")
		pdf.MultiCell(w, h, v, "", "", false)
	}

	if opt.Template == "compact" {
		kv("Username", it.Username)
		kv("Passwort", mask(it.Password, opt.MaskPassword))
		if len(it.URLs) > 0 { kv("URL", strings.Join(it.URLs, " ")) }
		kv("TOTP", it.TOTP)
		if it.Notes != "" {
			pdf.SetFontSize(10)
			kv("Notizen", it.Notes)
			pdf.SetFontSize(11)
		}
	} else {
		kv("Username", it.Username)
		kv("Passwort", mask(it.Password, opt.MaskPassword))
		if len(it.URLs) > 0 { kv("URL", strings.Join(it.URLs, " ")) }
		kv("TOTP", it.TOTP)
		if it.Notes != "" { kv("Notizen", it.Notes) }
		for k, v := range it.RawFields {
			if strings.EqualFold(k, "username") || strings.EqualFold(k, "password") || strings.EqualFold(k, "notes") {
				continue
			}
			kv(k, v)
		}
	}

	pdf.Ln(2)
	if pdf.GetY() > 270 {
		pdf.AddPage()
	}
}
