package rest

import (
	"net/http"

	"github.com/go-chi/render"
)

// DataCachedStatus type
type DataCachedStatus string

// Cached status enum
const (
	DataCached    DataCachedStatus = "Yes"
	DataNotCached DataCachedStatus = "No"
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

// JSONResponseFromBytes creates json body
func JSONResponseFromBytes(w http.ResponseWriter, r *http.Request, data []byte) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(data)
}

// NotFoundResponse creates empty json body and responds with 404 code
func NotFoundResponse(w http.ResponseWriter, r *http.Request) {
	render.Status(r, 404)
	render.JSON(w, r, map[string]interface{}{})
}
