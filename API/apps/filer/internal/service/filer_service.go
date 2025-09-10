package service

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/AmadoMuerte/BirthdayWish/API/apps/filer/internal/storage"
	filerProto "github.com/AmadoMuerte/BirthdayWish/API/proto/filer"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

type FilerService struct {
	storage *storage.Storage
	log     *slog.Logger
	filerProto.UnimplementedFilerServiceServer
}

func NewFilerService(storage *storage.Storage, log *slog.Logger) *FilerService {
	return &FilerService{storage: storage, log: log}
}

func (s *FilerService) LoadImage(ctx context.Context, req *filerProto.LoadImageRequest) (*filerProto.LoadImageResponse, error) {
	data, err := base64.StdEncoding.DecodeString(req.ImageBase64)
	if err != nil {
		s.log.Error("Invalid base64 data", "error", err)
		return nil, errors.New("invalid base64 data")
	}

	contentType := http.DetectContentType(data)
	var fileExt string

	switch contentType {
	case "image/jpeg":
		fileExt = ".jpg"
	case "image/png":
		fileExt = ".png"
	default:
		s.log.Error("Unsupported image format", "contentType", contentType)
		return nil, errors.New("unsupported image format")
	}

	uniqueName := uuid.New().String() + fileExt
	storagePath := fmt.Sprintf("images/%s", uniqueName)

	_, err = s.storage.Client.PutObject(
		ctx,
		s.storage.BucketName,
		storagePath,
		bytes.NewReader(data),
		int64(len(data)),
		minio.PutObjectOptions{ContentType: contentType},
	)
	if err != nil {
		s.log.Error("Failed to upload image", "error", err)
		return nil, errors.New("failed to upload image")
	}

	publicURL := fmt.Sprintf("%s/%s/%s", s.storage.Client.EndpointURL(), s.storage.BucketName, storagePath)

	return &filerProto.LoadImageResponse{ImageUrl: publicURL}, nil
}
