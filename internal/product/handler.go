package product

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
	"github.com/gorilla/schema"
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

	newSchema := schema.NewDecoder()
	newSchema.IgnoreUnknownKeys(true)
	if err = newSchema.Decode(&req, r.URL.Query()); err != nil {
		slog.Error(err.Error())
		response.JSON(w, http.StatusBadRequest, response.ResponseBody{})
		return
	}

	req.UserID = userID

	slog.Info("req", req.Search)

	err = req.Validate()
	if err != nil {
		response.JSON(w, http.StatusBadRequest, response.ResponseBody{
			Error: err.Error(),
		})
		return
	}

	resp = h.service.List(r.Context(), req)
	response.JSON(w, resp.Code, response.ResponseBody{
		Message: resp.Message,
		Data:    resp.Data,
		Meta:    resp.Meta,
	})
}

func (h *Handler) PatchProduct(w http.ResponseWriter, r *http.Request) {
	var req UpdateProductPayload
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

	params := mux.Vars(r)
	uid, err := uuid.Parse(params["productId"])
	if err != nil {
		response.JSON(w, http.StatusBadRequest, response.ResponseBody{
			Message: "Failed to parse UUID",
			Error:   err.Error(),
		})
		return
	}

	req.ProductUID = uid

	err = req.Validate()
	if err != nil {
		response.JSON(w, http.StatusBadRequest, response.ResponseBody{
			Error: err.Error(),
		})
		return
	}

	resp = h.service.Update(r.Context(), req)
	response.JSON(w, resp.Code, response.ResponseBody{
		Message: resp.Message,
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
