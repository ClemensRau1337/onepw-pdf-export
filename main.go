package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"golang.org/x/term"

	"github.com/example/onepw-pdf-export/pkg/model"
	"github.com/example/onepw-pdf-export/pkg/onepux"
	"github.com/example/onepw-pdf-export/pkg/op"
	"github.com/example/onepw-pdf-export/pkg/pdfwriter"
)

var version = "0.4.0"

func main() {
	// Flags (can be omitted; we will prompt interactively)
	var (
		out          string
		template     string
		maskPw       bool
		confirmRisk  bool
		search       string
		password     string
		noInteractive bool
		csvPath      string
		onepuxPath   string
		vaults       multiFlag
	)

	flag.StringVar(&out, "out", "", "Zieldatei (PDF)")
	flag.StringVar(&template, "template", "", "Layout-Vorlage: compact|detailed (optional)")
	flag.BoolVar(&maskPw, "mask-passwords", false, "Passwörter maskieren (optional)")
	flag.BoolVar(&confirmRisk, "i-understand-the-risk", false, "Sicherheitsbestätigung (required unless interactive confirmed)")
	flag.StringVar(&search, "search", "", "Einfache Volltextsuche (optional)")
	flag.StringVar(&password, "password", "", "PDF-Passwort (ansonsten verdeckte Abfrage)")
	flag.BoolVar(&noInteractive, "no-interactive", false, "Interaktive Eingaben ausschalten (z. B. CI)")
	flag.StringVar(&csvPath, "csv", "", "CSV-Datei als Quelle statt op (optional)")
	flag.StringVar(&onepuxPath, "onepux", "", ".1pux-Datei als Quelle statt op (optional)")
	flag.Var(&vaults, "vault", "Name eines Tresors (mehrfach möglich; nur mit op)")
	flag.Parse()

	mode := detectMode(csvPath, onepuxPath)
	if !noInteractive {
		// 1) Risk acceptance if not given
		if !confirmRisk {
			if !promptYesNo("⚠️  Du exportierst ALLE Passwörter in ein PDF. Bist du dir des Risikos bewusst und willst fortfahren? (ja/nein): ") {
				fail(errors.New("abgebrochen"))
			}
			confirmRisk = true
		}

		// 2) Choose mode if not specified
		if mode == "" {
			fmt.Fprintln(os.Stderr, "Quelle nicht angegeben. Wähle:")
			fmt.Fprintln(os.Stderr, "  [1] Direkt aus 1Password (op)")
			fmt.Fprintln(os.Stderr, "  [2] CSV-Datei")
			fmt.Fprintln(os.Stderr, "  [3] 1PUX-Datei (experimentell)")
			choice := promptString("Auswahl (1/2/3): ")
			switch strings.TrimSpace(choice) {
			case "2":
				csvPath = promptString("Pfad zur CSV-Datei: ")
				mode = "csv"
			case "3":
				onepuxPath = promptString("Pfad zur 1PUX-Datei: ")
				mode = "1pux"
			default:
				mode = "op"
			}
		}

		// 3) Output file
		if strings.TrimSpace(out) == "" {
			out = promptStringDefault("Pfad zur Ausgabe-PDF", "vault-export.pdf")
		}
		// ensure .pdf extension
		if strings.ToLower(filepath.Ext(out)) != ".pdf" {
			out += ".pdf"
		}

		// 4) Template
		if template == "" {
			template = strings.ToLower(promptStringDefault("Layout (compact/detailed)", "compact"))
			if template != "compact" && template != "detailed" {
				template = "compact"
			}
		}

		// 5) Password
		if strings.TrimSpace(password) == "" {
			var err error
			password, err = promptPassword()
			if err != nil { fail(err) }
		}

		// 6) Optional: mask and search
		if !maskPw {
			if promptYesNo("Passwörter maskieren (•)? (ja/nein): ") {
				maskPw = true
			}
		}
		if strings.TrimSpace(search) == "" {
			search = promptStringDefault("Optional: Suchbegriff (leer lassen für alle)", "")
		}

		// 7) Vault selection (op only)
		if mode == "op" && len(vaults) == 0 {
			// List vaults
			fmt.Fprintln(os.Stderr, "Lade Tresore...")
			vlist, err := op.ListVaults()
			if err != nil { fail(fmt.Errorf("op vault list: %w", err)) }
			if len(vlist) == 0 { fail(errors.New("keine Tresore gefunden")) }
			fmt.Fprintln(os.Stderr, "Bitte wähle Tresor(e) (Mehrfachauswahl mit Komma):")
			for i, v := range vlist {
				fmt.Fprintf(os.Stderr, "  [%d] %s\n", i+1, v.Name)
			}
			sel := promptString("Auswahl (z. B. 1,3 oder leer für alle): ")
			if strings.TrimSpace(sel) != "" {
				parts := strings.Split(sel, ",")
				for _, p := range parts {
					p = strings.TrimSpace(p)
					var idx int
					_, err := fmt.Sscanf(p, "%d", &idx)
					if err == nil && idx >= 1 && idx <= len(vlist) {
						vaults = append(vaults, vlist[idx-1].Name)
					}
				}
			}
		}
	} else {
		// Non-interactive: require required params
		if !confirmRisk {
			fail(errors.New("--i-understand-the-risk erforderlich im --no-interactive Modus"))
		}
		if strings.TrimSpace(out) == "" {
			fail(errors.New("--out ist erforderlich im --no-interactive Modus"))
		}
		if strings.TrimSpace(password) == "" {
			fail(errors.New("--password ist erforderlich im --no-interactive Modus"))
		}
		if template == "" {
			template = "compact"
		}
	}

	// Run export
	switch mode {
	case "csv":
		runCSV(csvPath, out, template, maskPw, search, password)
	case "1pux":
		runOnePUX(onepuxPath, out, template, maskPw, search, password)
	default:
		runOP(vaults, out, template, maskPw, search, password)
	}
}

func detectMode(csvPath, onepuxPath string) string {
	if strings.TrimSpace(csvPath) != "" {
		return "csv"
	}
	if strings.TrimSpace(onepuxPath) != "" {
		return "1pux"
	}
	return "op"
}

func promptString(label string) string {
	fmt.Fprint(os.Stderr, label)
	if !strings.HasSuffix(label, " ") { fmt.Fprint(os.Stderr, " ") }
	reader := bufio.NewReader(os.Stdin)
	s, _ := reader.ReadString('\n')
	return strings.TrimSpace(s)
}

func promptStringDefault(label, def string) string {
	if def != "" {
		returned := promptString(fmt.Sprintf("%s [%s]: ", label, def))
		if strings.TrimSpace(returned) == "" {
			return def
		}
		return returned
	}
	return promptString(label + ": ")
}

func promptYesNo(label string) bool {
	for {
		ans := strings.ToLower(promptString(label))
		if ans == "j" || ans == "ja" || ans == "y" || ans == "yes" {
			return true
		}
		if ans == "n" || ans == "nein" || ans == "no" {
			return false
		}
		fmt.Fprintln(os.Stderr, "Bitte 'ja' oder 'nein' eingeben.")
	}
}

func promptPassword() (string, error) {
	fmt.Fprint(os.Stderr, "PDF-Passwort: ")
	pw1, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Fprintln(os.Stderr)
	if err != nil {
		return "", err
	}
	if len(strings.TrimSpace(string(pw1))) == 0 {
		return "", errors.New("leeres Passwort ist nicht erlaubt")
	}
	fmt.Fprint(os.Stderr, "Passwort wiederholen: ")
	pw2, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Fprintln(os.Stderr)
	if err != nil {
		return "", err
	}
	if string(pw1) != string(pw2) {
		return "", errors.New("Passwörter stimmen nicht überein")
	}
	return string(pw1), nil
}

func runOP(vaults []string, out, template string, maskPw bool, search, password string) {
	// 1) Items via op (liste)
	fmt.Fprintln(os.Stderr, "Lade Item-Liste...")
	list, err := op.ListItems()
	if err != nil {
		fail(fmt.Errorf("op: %w", err))
	}

	// 2) Vorfilterung anhand Vault/Query auf Listeneinträgen
	wantVault := map[string]bool{}
	for _, v := range vaults {
		wantVault[strings.ToLower(strings.TrimSpace(v))] = true
	}

	type pair struct{ id, vault, title string }
	ids := make([]pair, 0, len(list))
	q := strings.TrimSpace(strings.ToLower(search))
	for _, e := range list {
		if len(wantVault) > 0 && !wantVault[strings.ToLower(e.Vault.Name)] {
			continue
		}
		if q != "" && !strings.Contains(strings.ToLower(e.Title), q) {
			continue
		}
		ids = append(ids, pair{e.ID, e.Vault.Name, e.Title})
	}

	if len(ids) == 0 {
		fail(errors.New("keine Items nach Filter gefunden"))
	}

	// 3) Fortschritt anzeigen
	fmt.Fprintf(os.Stderr, "Lade Details & rendere PDF (%d Items)...\n", len(ids))
	stop := make(chan struct{})
	go spinner("Bitte warten", stop)

	items := make([]model.Item, 0, len(ids))
	for i, p := range ids {
		it, err := op.GetItemDetails(p.id)
		if err == nil {
			items = append(items, it)
		}
		if (i+1)%50 == 0 {
			fmt.Fprintf(os.Stderr, "… %d/%d verarbeitet\n", i+1, len(ids))
		}
	}
	close(stop)
	fmt.Fprintln(os.Stderr, "Details geladen. Erzeuge PDF...")

	if err := pdfwriter.WritePDF(out, items, pdfwriter.Options{
		Template:     template,
		MaskPassword: maskPw,
		Source:       "op",
		UserPassword: password,
	}); err != nil {
		fail(err)
	}
	fmt.Println("OK:", out)
}

func runCSV(csvPath, out, template string, maskPw bool, search, password string) {
	if strings.TrimSpace(csvPath) == "" {
		fail(errors.New("--csv Pfad fehlt"))
	}
	items, err := model.FromCSV(csvPath, ",")
	if err != nil {
		fail(err)
	}
	filtered := filterItems(items, nil, search)
	if err := pdfwriter.WritePDF(out, filtered, pdfwriter.Options{
		Template:     template,
		MaskPassword: maskPw,
		Source:       "csv",
		UserPassword: password,
	}); err != nil {
		fail(err)
	}
	fmt.Println("OK:", out)
}

func runOnePUX(onepuxPath, out, template string, maskPw bool, search, password string) {
	if strings.TrimSpace(onepuxPath) == "" {
		fail(errors.New("--onepux Pfad fehlt"))
	}
	items, err := onepux.FromFile(onepuxPath)
	if err != nil {
		fail(err)
	}
	filtered := filterItems(items, nil, search)
	if err := pdfwriter.WritePDF(out, filtered, pdfwriter.Options{
		Template:     template,
		MaskPassword: maskPw,
		Source:       "1pux",
		UserPassword: password,
	}); err != nil {
		fail(err)
	}
	fmt.Println("OK:", out)
}

type multiFlag []string
func (m *multiFlag) String() string { return strings.Join(*m, ",") }
func (m *multiFlag) Set(v string) error { *m = append(*m, v); return nil }

func filterItems(in []model.Item, vaults []string, query string) []model.Item {
	vset := map[string]bool{}
	for _, v := range vaults {
		vset[strings.ToLower(strings.TrimSpace(v))] = true
	}
	out := make([]model.Item, 0, len(in))
	for _, it := range in {
		if len(vset) > 0 && !vset[strings.ToLower(it.Vault)] {
			continue
		}
		if q := strings.TrimSpace(strings.ToLower(query)); q != "" {
			if !strings.Contains(strings.ToLower(it.Title), q) &&
				!strings.Contains(strings.ToLower(it.Username), q) &&
				!strings.Contains(strings.ToLower(strings.Join(it.URLs, " ")), q) {
				continue
			}
		}
		out = append(out, it)
	}
	return out
}

func spinner(msg string, stop <-chan struct{}) {
	frames := []rune{'⠋','⠙','⠹','⠸','⠼','⠴','⠦','⠧','⠇','⠏'}
	i := 0
	for {
		select {
		case <-stop:
			fmt.Fprint(os.Stderr, "\r")
			return
		default:
			fmt.Fprintf(os.Stderr, "\r%s %c", msg, frames[i%len(frames)])
			time.Sleep(90 * time.Millisecond)
			i++
		}
	}
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, "Fehler:", err)
	os.Exit(1)
}
