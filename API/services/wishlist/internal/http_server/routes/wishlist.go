package routes

import (
	"log/slog"

	"github.com/AmadoMuerte/BirthdayWish/API/internal/config"
	"github.com/AmadoMuerte/BirthdayWish/API/internal/http_server/handlers/wishlist"
	"github.com/AmadoMuerte/BirthdayWish/API/internal/storage"
	"github.com/go-chi/chi/v5"
)

func NewWishlistRouter(cfg *config.Config, storage *storage.Storage) *chi.Mux {
	router := chi.NewRouter()

	wishlisthandler := wishlist.New(cfg, storage, slog.Default())
	router.Get("/{user_id}", wishlisthandler.GetWishlist)
	router.Post("/", wishlisthandler.AddToWishlist)
	router.Delete("/{wish_id}", wishlisthandler.RemoveFromWishlist)

	return router
}
