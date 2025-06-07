package wishlist

import (
	"net/http"
	"strconv"

	"github.com/AmadoMuerte/BirthdayWish/API/internal/pkg/response"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)



func (h *WishlistHandler) GetWish(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	wishID, err := strconv.ParseInt(chi.URLParam(r, "wish_id"), 10, 64)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error("invalid wish id"))
		return
	}

	userID, err := strconv.ParseInt(chi.URLParam(r, "user_id"), 10, 64)
	if err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, response.Error("invalid user id"))
		return
	}

	if wish, err := h.storage.GetWish(ctx, userID, wishID); err != nil {
		h.log.Error("error getting wish", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, response.Error("Internal Server Error"))
		return
	} else {
		render.JSON(w, r, wish)
	}
}


