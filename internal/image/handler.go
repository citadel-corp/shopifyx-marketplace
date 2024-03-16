package image

import (
	"net/http"
	"strings"

	"github.com/citadel-corp/shopifyx-marketplace/internal/common/response"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) UploadToS3(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 2*1024*1024) // 2 MB

	if err := r.ParseMultipartForm(2 * 1024 * 1024); err != nil {
		response.JSON(w, http.StatusBadRequest, response.ResponseBody{
			Message: "File must be smaller than 2 MB",
		})
		return
	}
	file, header, err := r.FormFile("file")
	if file == nil {
		response.JSON(w, http.StatusBadRequest, response.ResponseBody{
			Message: "File should not be empty",
		})
		return
	}

	if header.Size < 1024*10 {
		response.JSON(w, http.StatusBadRequest, response.ResponseBody{
			Message: "File must be larger than 10 KB",
		})
		return
	}

	if err != nil {
		response.JSON(w, http.StatusBadRequest, response.ResponseBody{
			Message: "Failed to parse file",
			Error:   err.Error(),
		})
		return
	}
	defer file.Close()
	mimeType := header.Header.Get("Content-Type")
	if mimeType != "image/jpeg" {
		fileNameSplit := strings.Split(header.Filename, ".")
		ext := fileNameSplit[len(fileNameSplit)-1]
		isJPEG := ext == "jpeg" || ext == "jpg"
		if mimeType == "application/octet-stream" && !isJPEG {
			response.JSON(w, http.StatusBadRequest, response.ResponseBody{
				Message: "File is not a jpg/jpeg type",
			})
		}

		return
	}

	url, err := h.service.UploadToS3(r.Context(), file, header.Filename)
	if err != nil {
		response.JSON(w, http.StatusInternalServerError, response.ResponseBody{
			Message: "Unable to upload file",
			Error:   err.Error(),
		})
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{
		"imageUrl": url,
	})
}
