package product

import "net/http"

var (
	ErrorUnauthorized  = Response{Code: http.StatusForbidden, Message: "Unauthorized"}
	ErrorRequiredField = Response{Code: http.StatusBadRequest, Message: "Required field"}
	ErrorInternal      = Response{Code: http.StatusInternalServerError, Message: "Internal Server Error"}
)