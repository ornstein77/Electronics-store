package request

import (
	"kontakt/internal/core/domain"
	"kontakt/internal/core/transport/http/types"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestStatusValidation(t *testing.T) {
	for _, s := range []string{"Новый", "Подтверждён", "Собирается", "Готов к выдаче", "Завершён"} {
		r := httptest.NewRequest("PATCH", "/", strings.NewReader(`{"status":"`+s+`"}`))
		var v types.Status
		if e := Decode(httptest.NewRecorder(), r, &v); e != nil {
			t.Fatalf("%s: %v", s, e)
		}
	}
}
func TestRejectInvalidOrders(t *testing.T) {
	for _, body := range []string{`{}`, `{"name":"Test","phone":"1234567890","delivery":"delivery","items":[{"id":"00000000-0000-4000-8000-000000000001","qty":1}]}`, `{"name":"Test","phone":"1234567890","delivery":"pickup","items":[{"id":"00000000-0000-4000-8000-000000000001","qty":-1}]}`, `{} {}`} {
		r := httptest.NewRequest("POST", "/", strings.NewReader(body))
		var v domain.OrderInput
		if Decode(httptest.NewRecorder(), r, &v) == nil {
			t.Fatalf("accepted %s", body)
		}
	}
}
