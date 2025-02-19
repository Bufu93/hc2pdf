package main

import (
	"bytes"
	"html/template"
	"net/http"
	"os"

	pdf "github.com/Bufu93/hc2pdf/internal/service"
	"github.com/Bufu93/hc2pdf/internal/transport/rest"
	log "github.com/sirupsen/logrus"
)

var styles template.CSS


const (
	baseUrl = ":3000"
)



func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
	contents, err := os.ReadFile("./templates/content.styles.css")
	if err != nil {
		panic(err)
	}
	styles = template.CSS(contents)
}


var pdfTemplate = template.Must(template.ParseFiles("./templates/layout.html.tmpl", "./templates/content.html.tmpl", ))

type Page struct {
	Styles template.CSS
}


func main() {
	source := bytes.NewBuffer([]byte{})
	if err := pdfTemplate.Execute(source, Page{
		Styles: styles,
	}); err != nil {
		log.Fatalf("Error executing template: %v", err)
	}

	pdfInstance := pdf.NewPdf(source.Bytes()) 
	handler := rest.NewHandler(pdfInstance)

	srv := &http.Server{
		Addr:    baseUrl,
		Handler: handler.InitRouter(),
	}

	// Start the server
	log.Info("Starting server on port ", baseUrl)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

	
	
	
	
}