package bankaccount

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/citadel-corp/shopifyx-marketplace/internal/common/middleware"
	"github.com/citadel-corp/shopifyx-marketplace/internal/common/request"
	"github.com/citadel-corp/shopifyx-marketplace/internal/common/response"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateBankAccount(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		slog.Error(err.Error())
		response.JSON(w, http.StatusInternalServerError, response.ResponseBody{})
		return
	}

	var req CreateBankAccountPayload

	err = request.DecodeJSON(w, r, &req)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, response.ResponseBody{
			Message: "Failed to decode JSON",
			Error:   err.Error(),
		})
		return
	}
	bankAccountResp, err := h.service.Create(r.Context(), req, userID)
	if errors.Is(err, ErrValidationFailed) {
		response.JSON(w, http.StatusBadRequest, response.ResponseBody{
			Message: "Bad request",
			Error:   err.Error(),
		})
		return
	}
	if err != nil {
		response.JSON(w, http.StatusInternalServerError, response.ResponseBody{
			Message: "Internal server error",
			Error:   err.Error(),
		})
		return
	}
	response.JSON(w, http.StatusOK, response.ResponseBody{
		Message: "success",
		Data:    bankAccountResp,
	})
}

func getUserID(r *http.Request) (uint64, error) {
	var userID uint64
	var err error

	if authValue, ok := r.Context().Value(middleware.ContextAuthKey{}).(string); ok {
		userID, err = strconv.ParseUint(authValue, 10, 64)
		if err != nil {
			return 0, err
		}
	} else {
		slog.Error("cannot parse auth value from context")
		return 0, errors.New("cannot parse auth value from context")
	}

	return userID, nil
}
