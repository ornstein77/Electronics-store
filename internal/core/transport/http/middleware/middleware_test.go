package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAdmin(t *testing.T) {
	for _, tc := range []struct {
		token string
		want  int
	}{{"", 401}, {"Bearer wrong", 401}, {"Bearer secret", 204}} {
		r := httptest.NewRequest("POST", "/", nil)
		r.Header.Set("Authorization", tc.token)
		w := httptest.NewRecorder()
		Admin("secret", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(204) })(w, r)
		if w.Code != tc.want {
			t.Fatalf("got %d want %d", w.Code, tc.want)
		}
	}
}
