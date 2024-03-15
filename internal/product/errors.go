package product

import (
	"errors"
	"net/http"
)

var (
	ErrorForbidden     = Response{Code: http.StatusForbidden, Message: "Forbidden", Error: errors.New("Forbidden")}
	ErrorUnauthorized  = Response{Code: http.StatusUnauthorized, Message: "Unauthorized", Error: errors.New("Unauthorized")}
	ErrorRequiredField = Response{Code: http.StatusBadRequest, Message: "Required field"}
	ErrorInternal      = Response{Code: http.StatusInternalServerError, Message: "Internal Server Error", Error: errors.New("Internal Server Error")}
	ErrorBadRequest    = Response{Code: http.StatusBadRequest, Message: "Bad Request"}
	ErrorNoRecords     = Response{Code: http.StatusOK, Message: "No records found"}

	ErrorNotPurchasable    = Response{Code: http.StatusBadRequest, Message: "product is not purchasable"}
	ErrorInsufficientStock = Response{Code: http.StatusBadRequest, Message: "insufficient product stock"}
)
