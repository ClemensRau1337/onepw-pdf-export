package fonts

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/jung-kurt/gofpdf"
)

const FontName = "DejaVuSans"
const FontFile = "DejaVuSans.ttf"
const FontURL = "https://github.com/dejavu-fonts/dejavu-fonts/blob/version_2_37/ttf/DejaVuSans.ttf?raw=1"

// EnsureUTF8Font ensures a UTF‑8 capable font is registered in gofpdf.
// It tries to find/download DejaVuSans.ttf into a cache dir, then registers it.
func EnsureUTF8Font(pdf *gofpdf.Fpdf) error {
	path, err := ensureFontFile()
	if err != nil {
		return err
	}
	pdf.AddUTF8Font(FontName, "", path)
	return nil
}

func ensureFontFile() (string, error) {
	// Preferred cache dir
	cache := defaultCacheDir()
	if cache == "" {
		cache = "."
	}
	path := filepath.Join(cache, FontFile)
	if _, err := os.Stat(path); err == nil {
		return path, nil
	}
	if err := os.MkdirAll(cache, 0o755); err != nil {
		return "", err
	}
	// Download
	if err := download(FontURL, path); err != nil {
		return "", fmt.Errorf("failed to download %s: %w", FontFile, err)
	}
	return path, nil
}

func defaultCacheDir() string {
	if x := os.Getenv("ONEPW_PDF_FONT_DIR"); x != "" {
		return x
	}
	if xdg := os.Getenv("XDG_CACHE_HOME"); xdg != "" {
		return filepath.Join(xdg, "onepw-pdf-export")
	}
	home, _ := os.UserHomeDir()
	if home == "" {
		return ""
	}
	switch runtime.GOOS {
	case "windows":
		if local := os.Getenv("LOCALAPPDATA"); local != "" {
			return filepath.Join(local, "onepw-pdf-export", "fonts")
		}
		return filepath.Join(home, "AppData", "Local", "onepw-pdf-export", "fonts")
	default:
		return filepath.Join(home, ".cache", "onepw-pdf-export", "fonts")
	}
}

func download(url, dest string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New(resp.Status)
	}
	f, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, resp.Body)
	return err
}
