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
		switch {
		case errors.Is(err, ErrorUnauthorized.Error):
			response.JSON(w, ErrorUnauthorized.Code, response.ResponseBody{})
			return
		case errors.Is(err, ErrorForbidden.Error):
			response.JSON(w, ErrorForbidden.Code, response.ResponseBody{})
			return
		default:
			response.JSON(w, http.StatusInternalServerError, response.ResponseBody{})
			return
		}
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

	newSchema := schema.NewDecoder()
	newSchema.IgnoreUnknownKeys(true)
	if err = newSchema.Decode(&req, r.URL.Query()); err != nil {
		slog.Error(err.Error())
		response.JSON(w, http.StatusBadRequest, response.ResponseBody{})
		return
	}

	userID, err := getUserID(r)
	if err != nil {
		if req.UserOnly {
			switch {
			case errors.Is(err, ErrorUnauthorized.Error):
				response.JSON(w, ErrorUnauthorized.Code, response.ResponseBody{})
				return
			case errors.Is(err, ErrorForbidden.Error):
				response.JSON(w, ErrorForbidden.Code, response.ResponseBody{})
				return
			default:
				response.JSON(w, http.StatusInternalServerError, response.ResponseBody{})
				return
			}
		}
	}

	req.UserID = userID

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
		switch {
		case errors.Is(err, ErrorUnauthorized.Error):
			response.JSON(w, ErrorUnauthorized.Code, response.ResponseBody{})
			return
		case errors.Is(err, ErrorForbidden.Error):
			response.JSON(w, ErrorForbidden.Code, response.ResponseBody{})
			return
		default:
			response.JSON(w, http.StatusInternalServerError, response.ResponseBody{})
			return
		}
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

func (h *Handler) GetProduct(w http.ResponseWriter, r *http.Request) {
	var req GetProductPayload
	var resp Response
	var err error

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

	resp = h.service.Get(r.Context(), req)
	response.JSON(w, resp.Code, response.ResponseBody{
		Message: resp.Message,
		Data:    resp.Data,
	})
}

func (h *Handler) PurchaseProduct(w http.ResponseWriter, r *http.Request) {
	var req PurchaseProductPayload
	var resp Response
	var err error

	userID, err := getUserID(r)
	if err != nil {
		switch {
		case errors.Is(err, ErrorUnauthorized.Error):
			response.JSON(w, ErrorUnauthorized.Code, response.ResponseBody{})
			return
		case errors.Is(err, ErrorForbidden.Error):
			response.JSON(w, ErrorForbidden.Code, response.ResponseBody{})
			return
		default:
			response.JSON(w, http.StatusInternalServerError, response.ResponseBody{})
			return
		}
	}

	req.BuyerID = userID

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

	resp = h.service.Purchase(r.Context(), req)
	response.JSON(w, resp.Code, response.ResponseBody{
		Message: resp.Message,
		Data:    resp.Data,
	})
}

func (h *Handler) UpdateStockProduct(w http.ResponseWriter, r *http.Request) {
	var req UpdateStockPayload
	var resp Response
	var err error

	userID, err := getUserID(r)
	if err != nil {
		switch {
		case errors.Is(err, ErrorUnauthorized.Error):
			response.JSON(w, ErrorUnauthorized.Code, response.ResponseBody{})
			return
		case errors.Is(err, ErrorForbidden.Error):
			response.JSON(w, ErrorForbidden.Code, response.ResponseBody{})
			return
		default:
			response.JSON(w, http.StatusInternalServerError, response.ResponseBody{})
			return
		}
	}

	req.UserID = userID

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

	resp = h.service.UpdateStock(r.Context(), req)
	response.JSON(w, resp.Code, response.ResponseBody{
		Message: resp.Message,
		Data:    resp.Data,
	})
}

func getUserID(r *http.Request) (uint64, error) {
	var userID uint64
	var err error

	if authValue, ok := r.Context().Value(middleware.ContextAuthKey{}).(string); ok {
		userID, err = strconv.ParseUint(authValue, 10, 64)
		if err != nil {
			slog.Error("getUserID: %v", err)
			return 0, ErrorInternal.Error
		}
	} else {
		slog.Error("getUserID: cannot parse auth value from context")
		return 0, ErrorUnauthorized.Error
	}

	if userID == 0 {
		slog.Error("getUserID: userID is not set")
		return 0, ErrorForbidden.Error
	}

	return userID, nil
}
