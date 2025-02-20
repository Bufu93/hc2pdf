package pdf

import (
	"bytes"
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"github.com/Bufu93/hc2pdf/internal/transport/rest"
)

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

const ChromiumExecutable = "chromium"
// const ChromiumExecutable = "C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe"


type PdfService struct {
	source []byte
}

func NewPdf(s []byte) *PdfService {
 return &PdfService{
	source: s,
 }
}

func (p *PdfService) GeneratePdfFromSource(req rest.PdfRequest) ([]byte, error) {
	// Create a complete HTML document with CSS and JS included
	htmlContent := fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<title>PDF Document</title>
			<style>%s</style>
			<script>%s</script>
		</head>
		<body>
			%s
		</body>
		</html>
	`, req.CSS, req.JS, req.HTML)

	source := []byte(htmlContent) // Use the modified HTML content

	outDir, err := filepath.Abs("./tmp/dynamic")
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return nil, err
	}

	if err := clearDirectory(outDir); err != nil {
		return nil, err
	}

	name := generateFilename()
	tmpFile := path.Join(outDir, name+".html")
	outFile := path.Join(outDir, name+".pdf")
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

func(p *PdfService) GeneratePdfFromTemplateSource() ([]byte, error) {
	outDir, err := filepath.Abs("./tmp/static")
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return nil, err
	}

	if err := clearDirectory(outDir); err != nil {
		return nil, err
	}

	name := generateFilename()
	tmpFile := path.Join(outDir, name+".html")
	outFile := path.Join(outDir, name+".pdf")
	args := append(chromeFlags, fmt.Sprintf("--print-to-pdf=%s", outFile), tmpFile)

	os.WriteFile(tmpFile, p.source, 0o644)

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

func generateFilename() string {
	var bytes = make([]byte, 10)
	rand.Read(bytes[:])
	return fmt.Sprintf("%x", bytes)
}