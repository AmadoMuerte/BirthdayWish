package wishlist

import (
	"encoding/json"
	"net/http"

	"github.com/AmadoMuerte/BirthdayWish/API/pkg/response"
	"github.com/AmadoMuerte/BirthdayWish/API/services/wishlist/internal/models"
	"github.com/go-chi/render"
)

func (h *WishlistHandler) AddToWishlist(w http.ResponseWriter, r *http.Request) {
	var req models.Wishlist
	ctx := r.Context()

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("failed to decode request body")
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, response.Error("invalid request body"))
		return
	}

	if req.Name == "" {
		h.log.Error("invalid name")
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, response.Error("invalid name"))
		return
	}

	if req.Link == "" {
		h.log.Error("invalid link")
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, response.Error("invalid link"))
		return
	}

	if req.ImageUrl == "" || req.ImageName == "" {
		h.log.Error("invalid image")
		w.WriteHeader(http.StatusBadRequest)
		render.JSON(w, r, response.Error("invalid image"))
		return
	}

	if err := h.storage.AddToWishlist(ctx, req); err != nil {
		h.log.Error("failed to add wishlist item", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, response.Error("failed to add wishlist item"))
		return
	}

	w.WriteHeader(http.StatusCreated)
	render.JSON(w, r, response.Success("wishlist item added successfully"))
}
