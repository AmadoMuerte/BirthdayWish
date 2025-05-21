package imagehandler

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/AmadoMuerte/BirthdayWish/API/services/minio/internal/models"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

func (h *ImageHandler) Upload(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Data string `json:"data"`
		Type string `json:"type"`
	}

	if err := render.DecodeJSON(r.Body, &request); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid request format"})
		return
	}

	data, err := base64.StdEncoding.DecodeString(request.Data)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"error": "Invalid base64 data"})
		return
	}

	var fileExt string
	switch request.Type {
	case "image/jpeg":
		fileExt = ".jpg"
	case "image/png":
		fileExt = ".png"
	}

	uniqueName := uuid.New().String() + uuid.New().String() + uuid.New().String() + fileExt
	storagePath := fmt.Sprintf("images/%s", uniqueName)

	_, err = h.storage.Client.PutObject(
		r.Context(),
		h.storage.BucketName,
		storagePath,
		bytes.NewReader(data),
		int64(len(data)),
		minio.PutObjectOptions{ContentType: request.Type},
	)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "Failed to upload image"})
		return
	}

	imageRecord := models.ImageRecord{
		OriginalName: uniqueName,
		StoragePath:  storagePath,
		PublicURL:    fmt.Sprintf("%s/%s/%s", h.storage.Client.EndpointURL(), h.storage.BucketName, storagePath),
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, imageRecord)
}
