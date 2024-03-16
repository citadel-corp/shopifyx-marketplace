package bankaccount

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/citadel-corp/shopifyx-marketplace/internal/common/middleware"
	"github.com/citadel-corp/shopifyx-marketplace/internal/common/request"
	"github.com/citadel-corp/shopifyx-marketplace/internal/common/response"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
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

	var req CreateUpdateBankAccountPayload

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

func (h *Handler) ListBankAccount(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		slog.Error(err.Error())
		response.JSON(w, http.StatusInternalServerError, response.ResponseBody{})
		return
	}

	bankAccountResp, err := h.service.List(r.Context(), userID)
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

func (h *Handler) PartialUpdateBankAccount(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		slog.Error(err.Error())
		response.JSON(w, http.StatusInternalServerError, response.ResponseBody{})
		return
	}
	params := mux.Vars(r)
	uid, err := uuid.Parse(params["uuid"])
	if err != nil {
		response.JSON(w, http.StatusNotFound, response.ResponseBody{
			Message: "Not found",
			Error:   ErrNotFound.Error(),
		})
		return
	}

	var req CreateUpdateBankAccountPayload

	err = request.DecodeJSON(w, r, &req)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, response.ResponseBody{
			Message: "Failed to decode JSON",
			Error:   err.Error(),
		})
		return
	}
	bankAccountResp, err := h.service.PartialUpdate(r.Context(), req, uid, userID)
	if errors.Is(err, ErrValidationFailed) {
		response.JSON(w, http.StatusBadRequest, response.ResponseBody{
			Message: "Bad request",
			Error:   err.Error(),
		})
		return
	}
	if errors.Is(err, ErrNotFound) {
		response.JSON(w, http.StatusNotFound, response.ResponseBody{
			Message: "Not found",
			Error:   err.Error(),
		})
		return
	}
	if errors.Is(err, ErrForbidden) {
		response.JSON(w, http.StatusForbidden, response.ResponseBody{
			Message: "Forbidden",
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
		Message: "account updated successfully",
		Data:    bankAccountResp,
	})
}

func (h *Handler) DeleteBankAccount(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		slog.Error(err.Error())
		response.JSON(w, http.StatusInternalServerError, response.ResponseBody{})
		return
	}
	params := mux.Vars(r)
	uid, err := uuid.Parse(params["uuid"])
	if err != nil {
		response.JSON(w, http.StatusNotFound, response.ResponseBody{
			Message: "Not found",
			Error:   ErrNotFound.Error(),
		})
		return
	}

	err = h.service.Delete(r.Context(), uid, userID)
	if errors.Is(err, ErrNotFound) {
		response.JSON(w, http.StatusNotFound, response.ResponseBody{
			Message: "Not found",
			Error:   err.Error(),
		})
		return
	}
	if errors.Is(err, ErrForbidden) {
		response.JSON(w, http.StatusForbidden, response.ResponseBody{
			Message: "Forbidden",
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
		Message: "account deleted successfully",
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
