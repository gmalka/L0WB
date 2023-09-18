package resttransport

import (
	"fmt"
	"html/template"
	"l0wb/models"
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

type Orderer interface {
	Add(models.Order) error
	Get(OrderUID string) (models.Order, error)
}

type Handler struct {
	s    Orderer
	path string
}

func NewHandler(s Orderer, path string) Handler {
	return Handler{
		s:    s,
		path: path,
	}
}

func (h Handler) Init() http.Handler {
	r := chi.NewRouter()

	r.Get("/", h.OrderGet)

	return r
}

func (h Handler) OrderGet(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("uid")
	tmpl, err := template.ParseFiles(h.path + "index.html")
	if err != nil {
		log.Printf("cant parse file: %v\n", err)
		http.Error(w, "some server error", http.StatusBadRequest)
		return
	}

	if name == "" {
		w.WriteHeader(http.StatusOK)

		if err := tmpl.Execute(w, nil); err != nil {
			log.Printf("cant execute template: %v\n", err)
			http.Error(w, "some server error", http.StatusBadRequest)
			return
		}
		return
	}

	order, err := h.s.Get(name)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err := tmpl.Execute(w, fmt.Sprintf("cant find order with id %v", name)); err != nil {
			log.Printf("cant execute template: %v\n", err)
			http.Error(w, "some server error", http.StatusBadRequest)
			return
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := tmpl.Execute(w, string(order.Order)); err != nil {
		log.Printf("cant execute template: %v\n", err)
		http.Error(w, "some server error", http.StatusBadRequest)
		return
	}
}
