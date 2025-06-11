package wishlist

import (
	"net/http"
	"strconv"

	"github.com/AmadoMuerte/BirthdayWish/API/pkg/response"
	"github.com/go-chi/chi/v5"
)



func (h *WishlistHandler) GetWish(w http.ResponseWriter, r *http.Request) {
	op := "wishlist/GetWish"
	ctx := r.Context()

	wishID, err := strconv.ParseInt(chi.URLParam(r, "wish_id"), 10, 64)
	if err != nil {
		h.log.Error(op + ": invalid wish_id", "error", err)
		response.ErrorResponseJSON(w,r,http.StatusBadRequest, "Invalid request body")
		return
	}

	userID, err := strconv.ParseInt(chi.URLParam(r, "user_id"), 10, 64)
	if err != nil {
		h.log.Error(op + ": invalid user id", "error", err)
		response.ErrorResponseJSON(w,r,http.StatusBadRequest, "Invalid request body")
		return
	}

	if wish, err := h.storage.GetWish(ctx, userID, wishID); err != nil {
		h.log.Error(op + ": error getting wish", "error", err)
		response.ErrorResponseJSON(w,r,http.StatusInternalServerError, "Internal Server Error")
		return
	} else {
		response.SuccessResponse(w,r, http.StatusOK, wish)
	}
}


