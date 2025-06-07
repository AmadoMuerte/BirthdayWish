package wishlist

import (
	"net/http"
	"strconv"

	"github.com/AmadoMuerte/BirthdayWish/API/internal/pkg/response"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func (h *WishlistHandler) RemoveFromWishlist(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	wishID, err := strconv.ParseInt(chi.URLParam(r, "wish_id"), 10, 64)
	if err != nil {
		h.log.Error("wish_id is empty")
		render.JSON(w, r, response.Error("wish_id is empty"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseInt(chi.URLParam(r, "user_id"), 10, 64)
	if err != nil {
		h.log.Error("user_id is empty")
		render.JSON(w, r, response.Error("error on the service"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.storage.RemoveFromWishlist(ctx, wishID, userID)
	if err != nil {
		h.log.Error("remove from wishlist: ", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, response.Error("internal server error"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
