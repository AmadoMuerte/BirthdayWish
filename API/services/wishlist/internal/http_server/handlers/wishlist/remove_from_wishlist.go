package wishlist

import (
	"net/http"
	"strconv"

	"github.com/AmadoMuerte/BirthdayWish/API/pkg/jwt"
	"github.com/AmadoMuerte/BirthdayWish/API/pkg/response"
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

	claims, err := jwt.GetClaims(ctx)
	if err != nil {
		h.log.Error("failed to get claims from token", "error", err)
		w.WriteHeader(http.StatusUnauthorized)
		render.JSON(w, r, response.Error("invalid token"))
		return
	}

	// exists, err := h.storage.UserExistsByID(ctx, claims.UserID)
	// if err != nil || !exists {
	// 	fmt.Printf("user id is %d", claims.UserID)
	// 	h.log.Error("user does not exist", "error", err)
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	render.JSON(w, r, response.Error("user does not exist"))
	// 	return
	// }

	err = h.storage.RemoveFromWishlist(ctx, wishID, claims.UserID)
	if err != nil {
		h.log.Error("remove from wishlist: ", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, response.Error("internal server error"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
