package rest

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type Pdf interface {
	GeneratePdfFromSource(source PdfRequest) ([]byte, error)
	GeneratePdfFromTemplateSource() ([]byte, error)
}

type Handler struct {
	pdfService Pdf

}

type PdfRequest struct {
	CSS  string `json:"css"` 
	HTML string `json:"html"`
	JS string `json:"js"` 
}

func NewHandler(pdf Pdf) *Handler {
	return &Handler{
		pdfService: pdf,
	}
	
	
}

func (h *Handler) InitRouter() *mux.Router {
	r := mux.NewRouter()
	r.Use(loggingMiddleware)
    r.HandleFunc("/live-pdf", h.livePdf).Methods(http.MethodGet)
    r.HandleFunc("/create-pdf", h.createPdf).Methods(http.MethodPost)
    return r
}


func(h *Handler) createPdf(w http.ResponseWriter, r *http.Request) {
	var req PdfRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Generate PDF from the HTML content
	pdf, err := h.pdfService.GeneratePdfFromSource(req) // Convert HTML string to byte slice
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Set the headers to indicate a file download
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=generated.pdf")
	w.WriteHeader(http.StatusOK)
	w.Write(pdf) 
}

func(h *Handler) livePdf(w http.ResponseWriter,r *http.Request) {
	pdf, err := h.pdfService.GeneratePdfFromTemplateSource()
		if err != nil {
			w.Header().Add("Content-Type", "text/plain; charset=utf-8")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}
	w.Write(pdf)
}






