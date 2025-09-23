package restful

import (
	"fmt"
	"net/http"

	"github.com/viniferr33/img-processor/internal/image"
	"github.com/viniferr33/img-processor/internal/utils"
	"github.com/viniferr33/img-processor/pkg/logger"
)

type imageHandler struct {
	imageService image.ImageService
}

func NewImageHandler(imageService image.ImageService) *imageHandler {
	return &imageHandler{
		imageService: imageService,
	}
}

func (h *imageHandler) handleUploadImage(w http.ResponseWriter, r *http.Request) {
	ownerId, ok := utils.GetUserIDFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("image")
	if err != nil {
		logger.Error(fmt.Sprintf("Error retrieving the file: %v", err))
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	data := make([]byte, handler.Size)
	_, err = file.Read(data)
	if err != nil {
		logger.Error(fmt.Sprintf("Error reading the file: %v", err))
		http.Error(w, "Error reading the file", http.StatusInternalServerError)
		return
	}

	_, err = h.imageService.UploadImage(r.Context(), data, handler.Filename, "", ownerId, 1)
	if err != nil {
		logger.Error(fmt.Sprintf("Error uploading the image: %v", err))
		http.Error(w, "Error uploading the image", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Image uploaded successfully"))
}
