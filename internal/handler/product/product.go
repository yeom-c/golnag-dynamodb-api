package product

import (
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/yeom-c/golnag-dynamodb-api/internal/controller/product"
	EntityProduct "github.com/yeom-c/golnag-dynamodb-api/internal/entity/product"
	"github.com/yeom-c/golnag-dynamodb-api/internal/handler"
	"github.com/yeom-c/golnag-dynamodb-api/internal/repository/adapter"
	"github.com/yeom-c/golnag-dynamodb-api/internal/rule"
	RuleProduct "github.com/yeom-c/golnag-dynamodb-api/internal/rule/product"
	Http "github.com/yeom-c/golnag-dynamodb-api/util/http"
)

type Handler struct {
	handler.Interface
	Controller product.Interface
	Rules      rule.Interface
}

func NewHandler(repository adapter.Interface) handler.Interface {
	return &Handler{
		Controller: product.NewController(repository),
		Rules:      RuleProduct.NewRule(),
	}
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	if chi.URLParam(r, "ID") != "" {
		h.getOne(w, r)
	} else {
		h.getAll(w, r)
	}
}

func (h *Handler) getOne(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "ID"))
	if err != nil {
		Http.StatusBadRequest(w, r, errors.New("ID is not uuid valid"))
		return
	}

	response, err := h.Controller.ListOne(id)
	if err != nil {
		Http.StatusInternalServerError(w, r, err)
		return
	}

	Http.StatusOK(w, r, response)
}

func (h *Handler) getAll(w http.ResponseWriter, r *http.Request) {
	response, err := h.Controller.ListAll()
	if err != nil {
		Http.StatusInternalServerError(w, r, err)
		return
	}

	Http.StatusOK(w, r, response)
}

func (h *Handler) Post(w http.ResponseWriter, r *http.Request) {
	productBody, err := h.getBodyAndValidate(r, uuid.Nil)
	if err != nil {
		Http.StatusBadRequest(w, r, err)
		return
	}

	id, err := h.Controller.Create(productBody)
	if err != nil {
		Http.StatusInternalServerError(w, r, err)
		return
	}
	Http.StatusOK(w, r, map[string]interface{}{"id": id})
}

func (h *Handler) Put(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "ID"))
	if err != nil {
		Http.StatusBadRequest(w, r, errors.New("ID is not uuid valid"))
		return
	}

	productBody, err := h.getBodyAndValidate(r, id)
	if err != nil {
		Http.StatusBadRequest(w, r, err)
		return
	}

	if err := h.Controller.Update(id, productBody); err != nil {
		Http.StatusInternalServerError(w, r, err)
		return
	}

	Http.StatusNoContent(w, r)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "ID"))
	if err != nil {
		Http.StatusBadRequest(w, r, errors.New("ID is not uuid valid"))
		return
	}

	if err := h.Controller.Delete(id); err != nil {
		Http.StatusInternalServerError(w, r, err)
		return
	}

	Http.StatusNoContent(w, r)
}

func (h *Handler) Options(w http.ResponseWriter, r *http.Request) {
	Http.StatusNoContent(w, r)
}

func (h *Handler) getBodyAndValidate(r *http.Request, id uuid.UUID) (*EntityProduct.Product, error) {
	productBody := &EntityProduct.Product{}
	body, err := h.Rules.ConvertIoReaderToStruct(r.Body, productBody)
	if err != nil {
		return &EntityProduct.Product{}, errors.New("body is required")
	}

	productParsed, err := EntityProduct.InterfaceToModel(body)
	if err != nil {
		return &EntityProduct.Product{}, errors.New("error on converting body to model")
	}

	setDefaultValues(productParsed, id)
	return productParsed, h.Rules.Validate(productParsed)
}

func setDefaultValues(product *EntityProduct.Product, id uuid.UUID) {
	if id == uuid.Nil {
		product.ID = uuid.New().String()
		product.CreatedAt = time.Now()
	} else {
		product.ID = id.String()
	}
}
