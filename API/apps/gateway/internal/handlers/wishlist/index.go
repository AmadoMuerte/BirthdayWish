package wishlist

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/AmadoMuerte/BirthdayWish/API/apps/gateway/internal/config"
	"github.com/AmadoMuerte/BirthdayWish/API/apps/gateway/internal/storage"
)

type WishlistHandler struct {
	cfg     *config.Config
	storage *storage.Storage
	log     *slog.Logger
}

type wishItemReq struct {
	Price     float64 `json:"price"`
	Link      string  `json:"link"`
	Image     string  `json:"image_data"`
	ImageType string  `json:"image_type"`
	Name      string  `json:"name"`
}

type wishItemRes struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Price     float64   `json:"price"`
	Link      string    `json:"link"`
	ImageUrl  string    `json:"image_url"`
	ImageName string    `json:"image_name"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type IWishlistHandler interface {
	AddWish(w http.ResponseWriter, r *http.Request)
	GetWish(w http.ResponseWriter, r *http.Request)
	UpdateWish(w http.ResponseWriter, r *http.Request)
	DeleteWish(w http.ResponseWriter, r *http.Request)
	GetWishlist(w http.ResponseWriter, r *http.Request)
	GetShareList(w http.ResponseWriter, r *http.Request)
}

func New(cfg *config.Config, storage *storage.Storage, log *slog.Logger) IWishlistHandler {
	return &WishlistHandler{cfg, storage, log}
}




