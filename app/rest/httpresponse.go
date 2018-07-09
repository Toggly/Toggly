package rest

import (
	"net/http"

	"github.com/go-chi/render"
)

// ErrorResponse creates {error: message} json body and responds with error code
func ErrorResponse(w http.ResponseWriter, r *http.Request, err error, code int) {
	render.Status(r, code)
	render.JSON(w, r, map[string]interface{}{"error": err.Error()})
}

// JSONResponse creates json body
func JSONResponse(w http.ResponseWriter, r *http.Request, data interface{}) {
	render.JSON(w, r, data)
}

// NotFoundResponse creates empty json body and responds with 404 code
func NotFoundResponse(w http.ResponseWriter, r *http.Request) {
	render.Status(r, 404)
	render.JSON(w, r, map[string]interface{}{})
}
