package service

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

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

func safeSubstring(s string, start, end int) string {
	if start > len(s) {
		return ""
	}
	if end > len(s) {
		end = len(s)
	}
	if start > end {
		return ""
	}
	return s[start:end]
}

func safeSubstringBytes(data []byte, start, end int) []byte {
	if start > len(data) {
		return []byte{}
	}
	if end > len(data) {
		end = len(data)
	}
	if start > end {
		return []byte{}
	}
	return data[start:end]
}

func decodeBase64Image(base64Data string) ([]byte, error) {
	if base64Data == "" {
		return nil, errors.New("empty base64 data")
	}

	if strings.Contains(base64Data, "base64,") {
		parts := strings.SplitN(base64Data, "base64,", 2)
		if len(parts) == 2 {
			base64Data = parts[1]
		}
	}

	base64Data = strings.TrimSpace(base64Data)
	base64Data = strings.ReplaceAll(base64Data, "\n", "")
	base64Data = strings.ReplaceAll(base64Data, "\r", "")
	base64Data = strings.ReplaceAll(base64Data, " ", "")
	base64Data = strings.ReplaceAll(base64Data, "\t", "")

	if base64Data == "" {
		return nil, errors.New("empty base64 data after cleaning")
	}

	if pad := len(base64Data) % 4; pad != 0 {
		base64Data += strings.Repeat("=", 4-pad)
	}

	data, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		data, err = base64.RawStdEncoding.DecodeString(base64Data)
		if err != nil {
			return nil, fmt.Errorf("base64 decode failed: %w", err)
		}
	}

	return data, nil
}
func (s *FilerService) LoadImage(ctx context.Context, req *filerProto.LoadImageRequest) (*filerProto.LoadImageResponse, error) {
	data, err := decodeBase64Image(req.ImageBase64)
	if err != nil {
		s.log.Error("Invalid base64 data",
			"error", err,
			"inputLength", len(req.ImageBase64),
			"sample", safeSubstring(req.ImageBase64, 0, 100))
		return nil, fmt.Errorf("invalid base64 data: %w", err)
	}

	contentType := http.DetectContentType(data)
	var fileExt string

	switch contentType {
	case "image/jpeg", "image/jpg":
		fileExt = ".jpg"
	case "image/png":
		fileExt = ".png"
	default:
		s.log.Error("Unsupported image format",
			"contentType", contentType,
			"dataLength", len(data),
			"firstBytes", fmt.Sprintf("%x", safeSubstringBytes(data, 0, 8)))
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
