package health

import (
	"errors"
	"net/http"

	"github.com/yeom-c/golnag-dynamodb-api/internal/handler"
	"github.com/yeom-c/golnag-dynamodb-api/internal/repository/adapter"
	Http "github.com/yeom-c/golnag-dynamodb-api/util/http"
)

type Handler struct {
	handler.Interface
	Repository adapter.Interface
}

func NewHandler(repository adapter.Interface) handler.Interface {
	return &Handler{
		Repository: repository,
	}
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	if !h.Repository.Health() {
		Http.StatusInternalServerError(w, r, errors.New("Relational database not alive"))
	}
	Http.StatusOK(w, r, "Service OK")
}

func (h *Handler) Post(w http.ResponseWriter, r *http.Request) {
	Http.StatusMethodNotAllowed(w, r)
}

func (h *Handler) Put(w http.ResponseWriter, r *http.Request) {
	Http.StatusMethodNotAllowed(w, r)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	Http.StatusMethodNotAllowed(w, r)
}

func (h *Handler) Options(w http.ResponseWriter, r *http.Request) {
	Http.StatusNoContent(w, r)
}
