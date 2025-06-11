package wishlist

import (
	"encoding/json"
	"net/http"

	"github.com/AmadoMuerte/BirthdayWish/API/apps/wishlister/internal/models"
	"github.com/AmadoMuerte/BirthdayWish/API/pkg/response"
)

func (h *WishlistHandler) AddWish(w http.ResponseWriter, r *http.Request) {
	op := "wishlist/AddWish"
	var req models.Wishlist
	ctx := r.Context()

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error(op + ": failed add wish item", "error", err)
		response.ErrorResponseJSON(w,r,http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Name == "" {
		h.log.Error(op + ": invalid name")
		response.ErrorResponseJSON(w,r,http.StatusBadRequest, "Invalid name")
		return
	}

	if err := h.storage.AddToWishlist(ctx, req); err != nil {
		h.log.Error(op + ": Failed to add wish item", "error", err)
		response.ErrorResponseJSON(w,r,http.StatusInternalServerError, "Failed to add wish item")
		return
	}

	response.SuccessResponse(w,r, http.StatusCreated, "Wish added successfully")
}
