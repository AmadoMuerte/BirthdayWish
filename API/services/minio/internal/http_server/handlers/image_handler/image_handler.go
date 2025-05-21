package imagehandler

import (
	"log/slog"
	"net/http"

	"github.com/AmadoMuerte/BirthdayWish/API/services/minio/internal/config"
	"github.com/AmadoMuerte/BirthdayWish/API/services/minio/internal/storage"
)

type ImageHandler struct {
	cfg     *config.Config
	storage *storage.Storage
	log     *slog.Logger
}

type IImageHandler interface {
	Upload(w http.ResponseWriter, r *http.Request)
}

func New(cfg *config.Config, storage *storage.Storage, log *slog.Logger) IImageHandler {
	return &ImageHandler{cfg, storage, log}
}
