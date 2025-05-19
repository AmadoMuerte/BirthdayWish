package wishlist

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/AmadoMuerte/BirthdayWish/API/internal/config"
	"github.com/AmadoMuerte/BirthdayWish/API/internal/storage"
)

type WishlistHandler struct {
	cfg     *config.Config
	storage *storage.Storage
	log     *slog.Logger
}

type IWishlistHandler interface {
	GetWishlist(w http.ResponseWriter, r *http.Request)
	AddToWishlist(w http.ResponseWriter, r *http.Request)
	RemoveFromWishlist(w http.ResponseWriter, r *http.Request)
}

func New(cfg *config.Config, storage *storage.Storage, log *slog.Logger) *WishlistHandler {
	return &WishlistHandler{cfg, storage, log}
}

type WishItem struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	Price     float64   `json:"price"`
	Link      string    `json:"link"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
