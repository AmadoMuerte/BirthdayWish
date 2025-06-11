package wishlist

import (
	"net/http"
	"strconv"

	"github.com/AmadoMuerte/BirthdayWish/API/apps/wishlister/internal/models"
	"github.com/AmadoMuerte/BirthdayWish/API/pkg/response"
	"github.com/go-chi/chi/v5"
)

func (h *WishlistHandler) GetWishlist(w http.ResponseWriter, r *http.Request) {
	op := "wishlist/GetWishlist"
	ctx := r.Context()

	user_id_str := chi.URLParam(r, "user_id")
	user_id, err := strconv.ParseInt(user_id_str, 10, 64)
	if err != nil {
		h.log.Error(op + ": user id is empty", "error", err)
		response.ErrorResponseJSON(w,r,http.StatusBadRequest, "Invalid request body")
		return
	}

	if wishlist, err := h.storage.GetWishlist(ctx, user_id); err != nil {
		h.log.Error(op + ": error getting wishlist", "error", err)
		response.ErrorResponseJSON(w,r,http.StatusInternalServerError, "Internal Server Error")
		return
	} else {
		var res []models.Wishlist
		for _, wish := range wishlist {
			res = append(res, wish)
		}
		response.SuccessResponse(w,r, http.StatusOK, res)
	}
}
