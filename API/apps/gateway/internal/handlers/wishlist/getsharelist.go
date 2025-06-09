package wishlist

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/AmadoMuerte/BirthdayWish/API/pkg/httphelper"
	"github.com/AmadoMuerte/BirthdayWish/API/pkg/response"
)



func (h *WishlistHandler) GetShareList(w http.ResponseWriter, r *http.Request) {
	op := "wishlist/getShareList"
	ctx := r.Context()

	authHeader := r.Header.Get("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == authHeader {
		h.log.Error(op+ ": token error")
		response.ErrorResponseJSON(w,r, http.StatusBadRequest, "internal error")
		return
	}

	userID, err := h.storage.GetUserIDByToken(ctx, token)
	if err != nil {
		h.log.Error(op + ": userID", "error", err)
		response.ErrorResponseJSON(w,r, http.StatusBadRequest, "internal error")
		return
	}

	path := fmt.Sprintf("%s/%d", h.cfg.Services.WishListAddr, userID)
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