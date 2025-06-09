package wishlist

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/AmadoMuerte/BirthdayWish/API/pkg/httphelper"
	"github.com/AmadoMuerte/BirthdayWish/API/pkg/jwt"
	"github.com/AmadoMuerte/BirthdayWish/API/pkg/response"
	"github.com/go-chi/chi/v5"
)


func (h *WishlistHandler) DeleteWish(w http.ResponseWriter, r *http.Request) {
	op := "wishlist/RemoveFromWishlist"
	ctx := r.Context()

	wishID, err := strconv.ParseInt(chi.URLParam(r, "wish_id"), 10, 64)
	if err != nil {
		h.log.Error(op + ": wish_id is empty", "error", err)
		response.ErrorResponseJSON(w,r, http.StatusBadRequest, "invalid request body")
		return
	}

	claims, err := jwt.GetClaims(ctx)
	if err != nil {
		h.log.Error(op + ": failed to get claims from token", "error", err)
		response.ErrorResponseJSON(w,r, http.StatusUnauthorized, "Unauthorized")
		return
	}

	exists, err := h.storage.UserExistsByID(ctx, claims.UserID)
	if err != nil || !exists {
		h.log.Error(op + ": user does not exist", "error", err)
		response.ErrorResponseJSON(w,r, http.StatusBadRequest, "invalid request body")
		return
	}

	path := fmt.Sprintf("%s/%d/%d", h.cfg.Services.WishListAddr, wishID, claims.UserID)

	resp, err := httphelper.DoRequest(ctx, "DELETE", path, nil, nil)
	if err != nil {
		h.log.Error(op + ": service call failed", "error", err)
		response.ErrorResponseJSON(w,r, http.StatusBadGateway, "service unavailable")
		return
	}

	defer resp.Body.Close()

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}