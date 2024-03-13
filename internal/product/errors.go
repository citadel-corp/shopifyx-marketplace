package product

import "net/http"

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

var (
	ErrorUnauthorized  = Error{Code: http.StatusForbidden, Message: "Unauthorized"}
	ErrorRequiredField = Error{Code: http.StatusBadRequest, Message: "Required field"}
	ErrorInternal      = Error{Code: http.StatusInternalServerError, Message: "Internal Server Error"}
)
