package rest

import "net/http"

// CtxValue type
type CtxValue string

// CtxValue enum
const (
	CtxValueOwner     CtxValue = "project-owner"
	CtxValueRequestID CtxValue = "request-id"
	CtxValueAuth      CtxValue = "auth"
)

// CtxOwner returns context value for project owner
func CtxOwner(r *http.Request) string {
	return r.Context().Value(CtxValueOwner).(string)
}
