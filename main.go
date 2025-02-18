package main

import (
	"bytes"
	"crypto/rand"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
)

var styles template.CSS

func init() {
	contents, err := os.ReadFile("./templates/styles.css")
	if err != nil {
		panic(err)
	}
	styles = template.CSS(contents)
}


var pdfTemplate = template.Must(template.ParseFiles("./templates/layout.html.tmpl"))

type Page struct {
	Styles template.CSS
}

func generateFilename() string {
	var bytes = make([]byte, 10)
	rand.Read(bytes[:])
	return fmt.Sprintf("%x", bytes)
}

var chromeFlags = []string{
	"--headless",
    "--disable-accelerated-2d-canvas",
    "--disable-gpu",
    "--allow-pre-commit-input",
    "--disable-background-networking",
    "--disable-background-timer-throttling",
    "--disable-backgrounding-occluded-windows",
    "--disable-breakpad",
    "--disable-client-side-phishing-detection",
    "--disable-component-extensions-with-background-pages",
    "--disable-component-update",
    "--disable-default-apps",
    "--disable-extensions",
    "--disable-features=Translate,BackForwardCache,AcceptCHFrame,MediaRouter,OptimizationHints",
    "--disable-hang-monitor",
    "--disable-ipc-flooding-protection",
    "--disable-popup-blocking",
    "--disable-prompt-on-repost",
    "--disable-renderer-backgrounding",
    "--disable-sync",
    "--enable-automation",
    "--enable-features=NetworkServiceInProcess2",
    "--export-tagged-pdf",
    "--force-color-profile=srgb",
    "--hide-scrollbars",
    "--metrics-recording-only",
    "--no-default-browser-check",
    "--no-first-run",
    "--no-service-autorun",
    "--password-store=basic",
    "--use-mock-keychain",
	"--no-sandbox",
}

const ChromiumExecutable = "C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe"

func generatePdfFromSource(source []byte) ([]byte, error) {
	outDir, err := filepath.Abs("./out")
	if err != nil {
		return nil, err
	}

	// Clear the output directory
	if err := clearDirectory(outDir); err != nil {
		return nil, err
	}
	

	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return nil, err
	}

	name := generateFilename()
	tmpFile := path.Join(outDir, name + ".html")
	outFile := path.Join(outDir, name + ".pdf")
	args := append(chromeFlags, fmt.Sprintf("--print-to-pdf=%s", outFile), tmpFile)
	
	os.WriteFile(tmpFile, source, 0o644)

	out := bytes.NewBuffer([]byte{})

	cmd := exec.Command(ChromiumExecutable, args...)
	cmd.Stderr = out
	cmd.Stdout = out

	if err := cmd.Run(); err != nil {
		log.Print(err)
		return nil, errors.New(out.String())
	}
	
	return os.ReadFile(outFile)
}

// Clear the directory function (same as before)
func clearDirectory(dir string) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		err := os.RemoveAll(path.Join(dir, file.Name()))
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	source := bytes.NewBuffer([]byte{})
	pdfTemplate.Execute(source, Page{
		Styles: styles,
	})

	http.Handle("/pdf", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pdf, err := generatePdfFromSource(source.Bytes())
		if err != nil {
			w.Header().Add("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
		w.Write(pdf)
	}))
	log.Print("Listening on port :3000")
	log.Fatal(http.ListenAndServe(":3000", nil))

	fmt.Println(source.String())
}