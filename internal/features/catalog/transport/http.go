package transport

import (
	"github.com/google/uuid"
	"kontakt/internal/core/domain"
	"kontakt/internal/core/transport/http/request"
	"kontakt/internal/core/transport/http/response"
	"kontakt/internal/features/catalog/service"
	"net/http"
)

type Handler struct{ Service *service.Service }

// List godoc
// @Summary List catalog products
// @Tags catalog
// @Produce json
// @Success 200 {array} domain.Product
// @Router /products [get]
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	v, e := h.Service.List(r.Context())
	if e != nil {
		response.Fail(w, e)
		return
	}
	response.JSON(w, 200, v)
}

// Get godoc
// @Summary Get product
// @Produce json
// @Param id path string true "Product UUID"
// @Success 200 {object} domain.Product
// @Router /products/{id} [get]
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id, e := uuid.Parse(r.PathValue("id"))
	if e != nil {
		response.Error(w, 400, "invalid id")
		return
	}
	v, e := h.Service.Get(r.Context(), id)
	if e != nil {
		response.Fail(w, e)
		return
	}
	response.JSON(w, 200, v)
}

// Save godoc
// @Summary Create or replace product (admin)
// @Security AdminToken
// @Accept json
// @Produce json
// @Param product body domain.Product true "Product"
// @Success 200 {object} domain.Product
// @Success 201 {object} domain.Product
// @Router /products [post]
// @Router /products/{id} [put]
func (h *Handler) Save(w http.ResponseWriter, r *http.Request) {
	var p domain.Product
	if e := request.Decode(w, r, &p); e != nil {
		response.Error(w, 400, "invalid product: "+e.Error())
		return
	}
	create := r.Method == "POST"
	if !create {
		id, e := uuid.Parse(r.PathValue("id"))
		if e != nil {
			response.Error(w, 400, "invalid id")
			return
		}
		p.ID = id
	}
	p, e := h.Service.Save(r.Context(), p, create)
	if e != nil {
		response.Fail(w, e)
		return
	}
	code := 200
	if create {
		code = 201
	}
	response.JSON(w, code, p)
}
