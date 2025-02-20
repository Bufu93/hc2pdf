package main

import (
	"bytes"
	"context"
	"errors"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Bufu93/hc2pdf/internal/server"
	pdf "github.com/Bufu93/hc2pdf/internal/service"
	"github.com/Bufu93/hc2pdf/internal/transport/rest"
	log "github.com/sirupsen/logrus"
)

var styles template.CSS
var html template.HTML


const (
	baseUrl = ":3000"
)



func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
	cssBytes, err := os.ReadFile("./templates/content.styles.css")
	if err != nil {
		panic(err)
	}
	htmlBytes, err := os.ReadFile("./templates/content.html.html")
	if err != nil {
		panic(err)
	}
	styles = template.CSS(cssBytes)
	html = template.HTML(htmlBytes)
}


var pdfTemplate = template.Must(template.ParseFiles("./templates/layout.html.tmpl" ))

type Page struct {
	Styles template.CSS
	HTML   template.HTML
}


func main() {
	source := bytes.NewBuffer([]byte{})
	if err := pdfTemplate.Execute(source, Page{
		Styles: styles,
		HTML:   html,
	}); err != nil {
		log.Fatalf("Error executing template: %v", err)
	}

	pdfInstance := pdf.NewPdf(source.Bytes()) 
	handler := rest.NewHandler(pdfInstance)
	srv := server.NewServer(handler.InitRouter())

	go func() {
		if err := srv.Run(); !errors.Is(err, http.ErrServerClosed) {
			log.Errorf("error occurred while running http server: %s\n", err.Error())
		}
	}()

	log.Info("Server started on port ", baseUrl)

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	const timeout = 5 * time.Second

	ctx, shutdown := context.WithTimeout(context.Background(), timeout)
	defer shutdown()

	if err := srv.Stop(ctx); err != nil {
		log.Errorf("failed to stop server: %v", err)
	}

	
}