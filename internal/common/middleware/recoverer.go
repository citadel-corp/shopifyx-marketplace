package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/citadel-corp/shopifyx-marketplace/internal/common/response"
)

func PanicRecoverer(next func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			r := recover()
			if r != nil {
				if r != http.ErrAbortHandler {
					slog.Error(fmt.Sprintf("Recovered from panic: %s", string(debug.Stack())))
				}
				response.JSON(w, http.StatusInternalServerError, response.ResponseBody{
					Message: "Internal server error",
				})
			}
		}()
		next(w, r)
	}
}
