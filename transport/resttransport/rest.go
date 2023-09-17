package resttransport

import (
	"fmt"
	"l0wb/models"
	"net/http"

	"github.com/go-chi/chi"
)

type Orderer interface {
	Add(models.Order) error
	Get(OrderUID string) (models.Order, error)
}

type Handler struct {
	s Orderer
}

func NewHandler(s Orderer) Handler {
	return Handler{
		s: s,
	}
}

func (h Handler) Init() http.Handler {
	r := chi.NewRouter()

	r.Get("/{name}", h.OrderGet)

	return r
}

func (h Handler) OrderGet(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")

	order, err := h.s.Get(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprintln(w, string(order.Order))
	w.WriteHeader(http.StatusOK)
}
