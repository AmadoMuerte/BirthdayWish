package wishlist

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/AmadoMuerte/BirthdayWish/API/apps/gateway/internal/handlers/filer"
	"github.com/AmadoMuerte/BirthdayWish/API/pkg/httphelper"
	"github.com/AmadoMuerte/BirthdayWish/API/pkg/jwt"
	"github.com/AmadoMuerte/BirthdayWish/API/pkg/response"
)

func (h *WishlistHandler) AddWish(w http.ResponseWriter, r *http.Request) {
	op := "wishlist/AddWish"
	ctx := r.Context()

	var wishItemReq wishItemReq
	if err := json.NewDecoder(r.Body).Decode(&wishItemReq); err != nil {
		h.log.Error(op+": failed to decode request body", "error", err)
		response.ErrorResponseJSON(w, r, http.StatusBadRequest, "invalid request body")
		return
	}
	defer r.Body.Close()

	claims, err := jwt.GetClaims(ctx)
	if err != nil {
		h.log.Error(op+": failed to get claims from token", "error", err)
		response.ErrorResponseJSON(w, r, http.StatusUnauthorized, "invalid token")
		return
	}

	if exists, err := h.storage.UserExistsByID(ctx, claims.UserID); err != nil || !exists {
		h.log.Error(op+": user does not exist", "user_id", claims.UserID, "error", err)
		response.ErrorResponseJSON(w, r, http.StatusUnauthorized, "user not found")
		return
	}

	wishItemRes := wishItemRes{
		UserID: claims.UserID,
		Price:  wishItemReq.Price,
		Link:   wishItemReq.Link,
		Name:   wishItemReq.Name,
	}

	var path string
	var resp *http.Response

	if wishItemReq.Image != "" {
		path := fmt.Sprintf("%s/%s", h.cfg.Services.Minio, "images")
		resp, err = httphelper.DoRequest(
			ctx,
			"POST",
			path,
			map[string]string{"data": wishItemReq.Image},
			map[string]string{"1": "application/json", "2": "Content-Type"},
		)
		if resp.StatusCode != 201 || err != nil {
			h.log.Error(op+": filer service is not create image", "error", err)
			w.WriteHeader(resp.StatusCode)
			io.Copy(w, resp.Body)
			return
		}

		var imageRecord filer.ImageRecord
		if err := json.NewDecoder(resp.Body).Decode(&imageRecord); err != nil {
			h.log.Error(op+": failed to decode request body", "error", err)
			response.ErrorResponseJSON(w, r, http.StatusBadRequest, "invalid request body")
			return
		}
		wishItemRes.Image = imageRecord.PublicURL
	}

	path = fmt.Sprintf("%s", h.cfg.Services.WishListAddr)
	headers := map[string]string{
		"1": "Content-Type",
		"2": "application/json",
	}

	resp, err = httphelper.DoRequest(ctx, "POST", path, wishItemRes, headers)
	if err != nil {
		h.log.Error(op+": failed add wish item", "error", err)
		response.ErrorResponseJSON(w, r, http.StatusInternalServerError, "internal server error")
		return
	}

	defer resp.Body.Close()

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}
