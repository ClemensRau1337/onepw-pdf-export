# onepw-pdf-export

## 🇩🇪 Deutsch

**[Deutsch 🇩🇪](#deutsch-) | [English 🇬🇧](#english-)**

---

## Deutsch 🇩🇪

⚠️ **HOHES RISIKO**: Dieses Tool exportiert **alle Passwörter** aus 1Password in ein **verschlüsseltes PDF**.  
Verwende es **nur für dich selbst** und behandle die PDF-Datei mit höchster Vorsicht. In Unternehmensumgebungen können Exporte gegen Richtlinien verstoßen oder Audits auslösen.

---

### ✨ Funktionen
- Export **live** via 1Password CLI (`op`)
- Export aus **CSV** (offizieller Export) oder **1PUX** (experimentell)
- Pflicht: PDF ist **immer** passwortgeschützt (AES/RC4 via gofpdf)
- Interaktive Passwortabfrage oder Übergabe per Flag
- Layouts: kompakt oder detailliert
- Filter: Vaults, Suchbegriffe
- Maskieren von Passwörtern optional möglich

---

### 🧰 Voraussetzungen
1. **1Password-CLI installieren**
   - macOS:  
     ```bash
     brew install --cask 1password-cli
     ```
   - Windows:  
     ```powershell
     winget install AgileBits.1Password.CLI
     ```
   - Linux: Siehe offizielle 1Password-Doku  
   Danach in der 1Password-App: *Einstellungen → Developer → „Integrate with 1Password CLI“ aktivieren*  
   Test:  
   ```bash
   op --version
   ```

2. **Go installieren**
   - macOS: `brew install go`
   - Windows: [Download von go.dev](https://go.dev/dl)
   - Linux: `sudo apt-get install golang` oder Download von go.dev  
   Test:  
   ```bash
   go version
   ```

3. **Repo bauen**
   ```bash
   git clone git@github.com:ClemensRau1337/onepw-pdf-export.git
   cd onepw-pdf-export
   go mod tidy
   go build -o onepw-pdf-export ./
   ```

---

### 🚀 Nutzung

**Interaktiv:** Wenn du Flags weglässt, fragt das Tool **ALLE Parameter** ab (Ausgabedatei, Risikoakzeptanz, Tresor-Auswahl, Passwort, Layout, Maskierung, Suche). Während des Exports siehst du einen **Spinner/Fortschritt**.


#### Standard via 1Password CLI

### Login to 1Password CLI (Siehe Requirements)
```bash
op signin
```

### Starten (interaktiv)
```bash
./onepw-pdf-export live
# prompts for all details
```

### CLI paramter
```bash
./onepw-pdf-export ./onepw-pdf-export --out my-vault.pdf --i-understand-the-risk
# prompts for password
```

#### CSV
```bash
./onepw-pdf-export csv --input ./examples/sample.csv --out logins.pdf --i-understand-the-risk
```

#### 1PUX (experimentell)
```bash
./onepw-pdf-export onepux --input export.1pux --out everything.pdf --i-understand-the-risk
```

---

### ⚙️ Flags

- `--out <file.pdf>` (**Pflicht**) – Zieldatei
- `--vault <name>` – Tresorfilter (nur Live-Modus, mehrfach)
- `--search <query>` – Textsuche über Titel, Benutzername, URLs
- `--template compact|detailed` (Standard: `compact`)
- `--mask-passwords` – ersetzt Passwörter durch •••••
- `--password <PW>` – setzt PDF-Passwort ohne Rückfrage  
- Ohne `--password`: verdeckte Eingabe mit Bestätigung
- `--i-understand-the-risk` (**Pflicht**) – Sicherheitsbestätigung

---

### 🔐 Sicherheit
- PDF immer verschlüsselt
- Keine temporären Dateien mit Klartext-Passwörtern
- Logs enthalten keine Geheimnisse
- Zusätzliche Absicherung: PDF in ein verschlüsseltes Archiv (7z, gpg, age) legen
- **Firmenumgebungen:** Policies/Audits beachten

---

### ❓ Troubleshooting

- **`op: command not found`** → 1Password-CLI installieren, PATH prüfen, Integration aktivieren
- **`missing go.sum entry`** → `go mod tidy`
- **Build-Fehler zu Imports** → ungenutzte Imports löschen (`sed -i.bak '/"path\/filepath"/d' pkg/onepux/onepux.go`)
- **macOS zsh PATH** (Apple Silicon):  
  ```bash
  echo 'export PATH="/opt/homebrew/bin:$PATH"' >> ~/.zshrc
  source ~/.zshrc
  ```

---

### 📄 Lizenz
MIT – siehe [LICENSE](LICENSE). Nutzung auf **eigene Gefahr**.  
Dieses Projekt ist **nicht mit 1Password affiliiert**.

---

## English 🇬🇧

⚠️ **HIGH RISK**: This tool exports **all your passwords** from 1Password into an **encrypted PDF**.  
Use it **only for yourself** and treat the resulting PDF with extreme caution. In corporate environments, exports may violate policies or trigger audits.

---

### ✨ Features
- Export **live** via 1Password CLI (`op`)
- Export from **CSV** (official export) or **1PUX** (experimental)
- Mandatory: PDF is **always** password-protected (AES/RC4 via gofpdf)
- Interactive password prompt or via flag
- Layouts: compact or detailed
- Filters: vaults, search queries
- Optional password masking

---

**UTF‑8:** Das PDF verwendet eine Unicode-Schrift (DejaVuSans). Fehlt sie lokal, lädt das Tool die Schrift automatisch herunter. Setze `ONEPW_PDF_FONT_DIR`, um den Speicherort zu steuern.

### 🧰 Requirements
1. **Install 1Password CLI**
   - macOS:  
     ```bash
     brew install --cask 1password-cli
     ```
   - Windows:  
     ```powershell
     winget install AgileBits.1Password.CLI
     ```
   - Linux: see official 1Password docs  
   Then in the 1Password app: *Settings → Developer → “Integrate with 1Password CLI”*  
   Test:  
   ```bash
   op --version
   ```

2. **Install Go**
   - macOS: `brew install go`
   - Windows: [Download from go.dev](https://go.dev/dl)
   - Linux: `sudo apt-get install golang` or download from go.dev  
   Test:  
   ```bash
   go version
   ```

3. **Build this repo**
   ```bash
   git clone git@github.com:ClemensRau1337/onepw-pdf-export.git
   cd onepw-pdf-export
   go mod tidy
   go build -o onepw-pdf-export ./
   ```

---

### 🚀 Usage

**Interactive:** If you omit flags, the tool will **prompt for everything** (output file, risk acceptance, vault selection, password, layout, masking, search). It also shows a **spinner/progress** during export.


#### Live via 1Password CLI

### Login to 1Password CLI (See Requirements)
```bash
op signin
```

### Starten (interaktiv)
```bash
./onepw-pdf-export live
# prompts for all details
```

### CLI paramter
```bash
./onepw-pdf-export ./onepw-pdf-export --out my-vault.pdf --i-understand-the-risk
# prompts for password
```

#### CSV
```bash
./onepw-pdf-export csv --input ./examples/sample.csv --out logins.pdf --i-understand-the-risk
```

#### 1PUX (experimental)
```bash
./onepw-pdf-export onepux --input export.1pux --out everything.pdf --i-understand-the-risk
```

---

### ⚙️ Flags

- `--out <file.pdf>` (**required**) – output file
- `--vault <name>` – filter by vault (live mode only, repeatable)
- `--search <query>` – text search over title, username, URLs
- `--template compact|detailed` (default: `compact`)
- `--mask-passwords` – replace passwords with •••••
- `--password <PW>` – set PDF password without prompt  
- Without `--password`: hidden interactive input with confirmation
- `--i-understand-the-risk` (**required**) – safety confirmation

---

### 🔐 Security
- PDF always encrypted
- No temporary plaintext files
- Logs never contain secrets
- Extra safety: place PDF inside an encrypted archive (7z, gpg, age)
- **Corporate environments:** respect policies/audits

---

### ❓ Troubleshooting

- **`op: command not found`** → install 1Password CLI, check PATH, enable integration
- **`missing go.sum entry`** → run `go mod tidy`
- **Build error about imports** → remove unused imports (`sed -i.bak '/"path\/filepath"/d' pkg/onepux/onepux.go`)
- **macOS zsh PATH** (Apple Silicon):  
  ```bash
  echo 'export PATH="/opt/homebrew/bin:$PATH"' >> ~/.zshrc
  source ~/.zshrc
  ```

---

### 📄 License
MIT – see [LICENSE](LICENSE). Use at your **own risk**.  
This project is **not affiliated with 1Password**.
