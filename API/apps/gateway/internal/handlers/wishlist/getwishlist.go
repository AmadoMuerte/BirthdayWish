package wishlist

import (
	"fmt"
	"io"
	"net/http"

	"github.com/AmadoMuerte/BirthdayWish/API/pkg/httphelper"
	"github.com/AmadoMuerte/BirthdayWish/API/pkg/jwt"
	"github.com/AmadoMuerte/BirthdayWish/API/pkg/response"
)

func (h *WishlistHandler) GetWishlist(w http.ResponseWriter, r *http.Request) {
	op := "wishlist/GetWishlist"
	ctx := r.Context()

	claims, err := jwt.GetClaims(ctx)
	if err != nil {
		h.log.Error(op + ": failed to get claims", "error", err)
		response.ErrorResponseJSON(w,r, http.StatusInternalServerError, "internal server error")
		return
	}

	path := fmt.Sprintf("%s/%d", h.cfg.Services.WishListAddr, claims.UserID)

	resp, err := httphelper.DoRequest(ctx, "GET", path, nil, nil)
	if err != nil {
		h.log.Error(op + ": service call failed", "error", err)
		response.ErrorResponseJSON(w,r, http.StatusBadGateway, "service unavailable")
		return
	}

	defer resp.Body.Close()

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}