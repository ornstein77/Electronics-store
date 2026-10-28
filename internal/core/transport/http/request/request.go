package request

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"io"
	"net/http"
)

var validate = validator.New()

func Decode(w http.ResponseWriter, r *http.Request, v any) error {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if e := d.Decode(v); e != nil {
		return e
	}
	var extra any
	if e := d.Decode(&extra); e != io.EOF {
		return &json.SyntaxError{}
	}
	return validate.Struct(v)
}
