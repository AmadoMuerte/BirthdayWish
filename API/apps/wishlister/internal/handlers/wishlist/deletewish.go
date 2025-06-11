package wishlist

import (
	"net/http"
	"strconv"

	"github.com/AmadoMuerte/BirthdayWish/API/pkg/response"
	"github.com/go-chi/chi/v5"
)

func (h *WishlistHandler) DeleteWish(w http.ResponseWriter, r *http.Request) {
	op := "wishlist/DeleteWish"
	ctx := r.Context()

	wishID, err := strconv.ParseInt(chi.URLParam(r, "wish_id"), 10, 64)
	if err != nil {
		h.log.Error(op + ": wish_id is empty", "error", err)
		response.ErrorResponseJSON(w,r,http.StatusBadRequest, "Invalid request body")
		return
	}

	userID, err := strconv.ParseInt(chi.URLParam(r, "user_id"), 10, 64)
	if err != nil {
		h.log.Error(op + ": user_id is empty", "error", err)
		response.ErrorResponseJSON(w,r,http.StatusBadRequest, "Invalid request body")
		return
	}

	err = h.storage.RemoveFromWishlist(ctx, wishID, userID)
	if err != nil {
		h.log.Error(op + ": internal server error", "error", err)
		response.ErrorResponseJSON(w,r,http.StatusInternalServerError, "internal server error")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
