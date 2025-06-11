package imagehandler

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/AmadoMuerte/BirthdayWish/API/apps/filer/internal/models"
	"github.com/AmadoMuerte/BirthdayWish/API/pkg/response"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	minio "github.com/minio/minio-go/v7"
)

func (h *ImageHandler) Upload(w http.ResponseWriter, r *http.Request) {
	op := "image/Upload"

	var request struct {
		Data string `json:"data"`
		Type string `json:"type"`
	}

	if err := render.DecodeJSON(r.Body, &request); err != nil {
		h.log.Error(op + ": Invalid request format", "error", err)
		response.ErrorResponseJSON(w,r, http.StatusBadRequest, "Invalid request format")
		return
	}

	data, err := base64.StdEncoding.DecodeString(request.Data)
	if err != nil {
		h.log.Error(op + ": Invalid base64 data", "error", err)
		response.ErrorResponseJSON(w,r, http.StatusBadRequest, "Invalid request format")
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
		h.log.Error(op + ": Failed to upload image", "error", err)
		response.ErrorResponseJSON(w,r, http.StatusInternalServerError, "Failed to upload image")
		return
	}

	imageRecord := models.ImageRecord{
		OriginalName: uniqueName,
		StoragePath:  storagePath,
		PublicURL:    fmt.Sprintf("%s/%s/%s", h.storage.Client.EndpointURL(), h.storage.BucketName, storagePath),
	}

	response.SuccessResponse(w,r,http.StatusCreated, imageRecord)
}
