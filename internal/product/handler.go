package product

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

func (h *Handler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var req CreateProductPayload
	var resp Response
	var err error

	userID, err := getUserID(r)
	if err != nil {
		slog.Error(err.Error())
		response.JSON(w, http.StatusInternalServerError, response.ResponseBody{})
		return
	}

	req.UserID = userID

	err = request.DecodeJSON(w, r, &req)
	if err != nil {
		response.JSON(w, http.StatusBadRequest, response.ResponseBody{
			Message: "Failed to decode JSON",
			Error:   err.Error(),
		})
		return
	}

	err = req.Validate()
	if err != nil {
		response.JSON(w, http.StatusBadRequest, response.ResponseBody{
			Error: err.Error(),
		})
		return
	}

	resp = h.service.Create(r.Context(), req)
	response.JSON(w, resp.Code, response.ResponseBody{
		Message: resp.Message,
	})
}

func (h *Handler) GetProductList(w http.ResponseWriter, r *http.Request) {
	var req ListProductPayload
	var resp Response
	var err error

	userID, err := getUserID(r)
	if err != nil {
		slog.Error(err.Error())
		response.JSON(w, http.StatusInternalServerError, response.ResponseBody{})
		return
	}

	req.UserID = userID

	// TODO: read query params

	err = req.Validate()
	if err != nil {
		response.JSON(w, http.StatusBadRequest, response.ResponseBody{
			Error: err.Error(),
		})
		return
	}

	data, pagination, resp := h.service.List(r.Context(), req)
	response.JSON(w, resp.Code, response.ResponseBody{
		Message: resp.Message,
		Data:    data,
		Meta:    pagination,
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
