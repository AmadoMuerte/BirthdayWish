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

func (h *WishlistHandler) GetWish(w http.ResponseWriter, r *http.Request) {
	op := "wishlist/GetWish"
	ctx := r.Context()

	wishID, err := strconv.ParseInt(chi.URLParam(r, "wish_id"), 10, 64)
	if err != nil {
		h.log.Error(op+": invalid wish id", "error", err)
		response.ErrorResponseJSON(w, r, http.StatusBadRequest, "invalid wish id")
		return
	}

	claims, err := jwt.GetClaims(ctx)

	path := fmt.Sprintf("%s/%d/%d", h.cfg.Services.WishListAddr, wishID, claims.UserID)

	resp, err := httphelper.DoRequest(ctx, "GET", path, nil, nil)
	if err != nil {
		h.log.Error(op+": service call failed", "error", err)
		response.ErrorResponseJSON(w, r, http.StatusBadGateway, "service unavailable")
		return
	}

	defer resp.Body.Close()

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}
