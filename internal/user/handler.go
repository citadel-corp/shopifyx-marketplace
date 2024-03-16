package user

import (
	"errors"
	"net/http"

	"github.com/citadel-corp/shopifyx-marketplace/internal/common/request"
	"github.com/citadel-corp/shopifyx-marketplace/internal/common/response"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserPayload

	err := request.DecodeJSON(w, r, &req)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, response.ResponseBody{
			Message: "Failed to decode JSON",
			Error:   err.Error(),
		})
		return
	}
	userResp, err := h.service.Create(r.Context(), req)
	if errors.Is(err, ErrUsernameAlreadyExists) {
		response.JSON(w, http.StatusConflict, response.ResponseBody{
			Message: "User already exists",
			Error:   err.Error(),
		})
		return
	}
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
	response.JSON(w, http.StatusCreated, response.ResponseBody{
		Message: "User registered successfully",
		Data:    userResp,
	})
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginPayload

	err := request.DecodeJSON(w, r, &req)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, response.ResponseBody{
			Message: "Failed to decode JSON",
			Error:   err.Error(),
		})
		return
	}
	userResp, err := h.service.Login(r.Context(), req)
	if errors.Is(err, ErrUserNotFound) {
		response.JSON(w, http.StatusNotFound, response.ResponseBody{
			Message: "Not found",
			Error:   err.Error(),
		})
		return
	}
	if errors.Is(err, ErrWrongPassword) {
		response.JSON(w, http.StatusBadRequest, response.ResponseBody{
			Message: "Bad request",
			Error:   err.Error(),
		})
		return
	}
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
		Message: "User logged successfully",
		Data:    userResp,
	})
}
