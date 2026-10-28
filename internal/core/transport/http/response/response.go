package response

import (
	"encoding/json"
	"errors"
	"go.uber.org/zap"
	core "kontakt/internal/core/errors"
	"kontakt/internal/core/transport/http/types"
	"net/http"
)

func JSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(v)
}
func Error(w http.ResponseWriter, code int, msg string) { JSON(w, code, types.Error{Error: msg}) }
func Fail(w http.ResponseWriter, e error) {
	switch {
	case errors.Is(e, core.NotFound):
		Error(w, 404, "not found")
	case errors.Is(e, core.Conflict):
		Error(w, 409, "insufficient stock")
	default:
		zap.L().Error("request failed", zap.Error(e))
		Error(w, 500, "internal server error")
	}
}
