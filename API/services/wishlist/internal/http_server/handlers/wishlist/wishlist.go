package wishlist

import (
	"log/slog"
	"net/http"

	"github.com/AmadoMuerte/BirthdayWish/API/services/wishlist/internal/config"
	"github.com/AmadoMuerte/BirthdayWish/API/services/wishlist/internal/storage"
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
	PartialUpdateWish(w http.ResponseWriter, r *http.Request)
}

func New(cfg *config.Config, storage *storage.Storage, log *slog.Logger) *WishlistHandler {
	return &WishlistHandler{cfg, storage, log}
}
