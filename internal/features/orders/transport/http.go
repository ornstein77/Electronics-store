package transport

import (
	"github.com/google/uuid"
	"kontakt/internal/core/domain"
	"kontakt/internal/core/transport/http/request"
	"kontakt/internal/core/transport/http/response"
	"kontakt/internal/core/transport/http/types"
	"kontakt/internal/features/orders/service"
	"net/http"
)

type Handler struct{ Service *service.Service }

// Create godoc
// @Summary Create order; prices are calculated by server
// @Tags orders
// @Accept json
// @Produce json
// @Param order body domain.OrderInput true "Order"
// @Success 201 {object} domain.Order
// @Failure 409 {object} types.Error
// @Router /orders [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var in domain.OrderInput
	if e := request.Decode(w, r, &in); e != nil {
		response.Error(w, 400, "invalid order: "+e.Error())
		return
	}
	seen := map[uuid.UUID]bool{}
	for _, i := range in.Items {
		if seen[i.ID] {
			response.Error(w, 400, "duplicate product id")
			return
		}
		seen[i.ID] = true
	}
	v, e := h.Service.Create(r.Context(), in)
	if e != nil {
		response.Fail(w, e)
		return
	}
	response.JSON(w, 201, v)
}

// List godoc
// @Summary Latest 500 orders (admin)
// @Security AdminToken
// @Produce json
// @Success 200 {array} domain.Order
// @Router /orders [get]
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	v, e := h.Service.List(r.Context())
	if e != nil {
		response.Fail(w, e)
		return
	}
	response.JSON(w, 200, v)
}

// Status godoc
// @Summary Update order status (admin)
// @Security AdminToken
// @Accept json
// @Produce json
// @Param id path string true "Order UUID"
// @Param status body types.Status true "Status"
// @Success 200 {object} types.Status
// @Router /orders/{id}/status [patch]
func (h *Handler) Status(w http.ResponseWriter, r *http.Request) {
	id, e := uuid.Parse(r.PathValue("id"))
	if e != nil {
		response.Error(w, 400, "invalid id")
		return
	}
	var in types.Status
	if e = request.Decode(w, r, &in); e != nil {
		response.Error(w, 400, "invalid status")
		return
	}
	if e = h.Service.Status(r.Context(), id, in.Status); e != nil {
		response.Fail(w, e)
		return
	}
	response.JSON(w, 200, in)
}
