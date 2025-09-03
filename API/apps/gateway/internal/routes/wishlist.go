package routes

import (
	"log/slog"

	"github.com/AmadoMuerte/BirthdayWish/API/apps/gateway/internal/handlers/wishlist"
	"github.com/AmadoMuerte/BirthdayWish/API/apps/gateway/internal/storage"
	"github.com/AmadoMuerte/BirthdayWish/API/pkg/config"
	"github.com/AmadoMuerte/BirthdayWish/API/pkg/redis"
	"github.com/go-chi/chi/v5"
)

func NewWishlistRouter(cfg *config.Config, storage *storage.Storage, rdb *redis.RDB) *chi.Mux {
	router := chi.NewRouter()

	wishlisthandler := wishlist.New(cfg, storage, rdb, slog.Default())
	router.Get("/{wish_id}", wishlisthandler.GetWish)
	router.Post("/", wishlisthandler.AddWish)
	router.Patch("/{wish_id}", wishlisthandler.UpdateWish)
	router.Delete("/{wish_id}", wishlisthandler.DeleteWish)
	router.Get("/", wishlisthandler.GetWishlist)

	return router
}
