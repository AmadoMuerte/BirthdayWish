package wishlist

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/AmadoMuerte/BirthdayWish/API/apps/gateway/internal/handlers/filer"
	"github.com/AmadoMuerte/BirthdayWish/API/pkg/httphelper"
	"github.com/AmadoMuerte/BirthdayWish/API/pkg/jwt"
	"github.com/AmadoMuerte/BirthdayWish/API/pkg/response"
	"github.com/go-chi/chi/v5"
)

func (h *WishlistHandler) UpdateWish(w http.ResponseWriter, r *http.Request) {
	op := "wishlist/PartialUpdateWish"
	ctx := r.Context()

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		h.log.Error(op+": failed to read request body", "error", err)
		response.ErrorResponseJSON(w, r, http.StatusBadRequest, "invalid request body")
		return
	}
	defer r.Body.Close()

	var wishItemReq wishItemReq
	if err := json.Unmarshal(bodyBytes, &wishItemReq); err != nil {
		h.log.Error(op+": failed to decode request body", "error", err)
		response.ErrorResponseJSON(w, r, http.StatusBadRequest, "invalid request body")
		return
	}

	var updateData map[string]any
	if err := json.Unmarshal(bodyBytes, &updateData); err != nil {
		h.log.Error(op+": failed to decode request body", "error", err)
		response.ErrorResponseJSON(w, r, http.StatusBadRequest, "invalid request body")
		return
	}

	wishID, err := strconv.ParseInt(chi.URLParam(r, "wish_id"), 10, 64)
	if err != nil {
		h.log.Error(op+": invalid wish_id", "error", err)
		response.ErrorResponseJSON(w, r, http.StatusBadRequest, "invalid request body")
		return
	}

	claims, err := jwt.GetClaims(ctx)
	if err != nil {
		h.log.Error(op+": failed to get claims from token", "error", err)
		response.ErrorResponseJSON(w, r, http.StatusUnauthorized, "Unauthorized")
		return
	}

	exists, err := h.storage.UserExistsByID(ctx, claims.UserID)
	if err != nil || !exists {
		h.log.Error(op+"user does not exist", "user_id", claims.UserID, "error", err)
		response.ErrorResponseJSON(w, r, http.StatusNotFound, "user not found")
		return
	}

	if wishItemReq.Image != "" {
		imagePath := fmt.Sprintf("%s/images", h.cfg.Services.Minio)
		imageResp, err := httphelper.DoRequest(
			ctx,
			"POST",
			imagePath,
			map[string]string{"data": wishItemReq.Image},
			map[string]string{"Content-Type": "application/json"},
		)
		if err != nil || imageResp.StatusCode != http.StatusCreated {
			h.log.Error(op+"filer service failed", "error", err)
			response.ErrorResponseJSON(w, r, http.StatusBadGateway, "service unavailable")
			return
		}
		defer imageResp.Body.Close()

		var imageRecord filer.ImageRecord
		if err := json.NewDecoder(imageResp.Body).Decode(&imageRecord); err != nil {
			h.log.Error(op+"failed to decode image response", "error", err)
			response.ErrorResponseJSON(w, r, http.StatusInternalServerError, "service unavailable")
			return
		}

		updateData["image"] = imageRecord.PublicURL
	}

	path := fmt.Sprintf("%s/%d/%d", h.cfg.Services.WishListAddr, wishID, claims.UserID)

	resp, err := httphelper.DoRequest(
		ctx,
		"PATCH",
		path,
		updateData,
		map[string]string{"Content-Type": "application/json"},
	)
	if err != nil {
		h.log.Error(op+"service call failed", "error", err)
		response.ErrorResponseJSON(w, r, http.StatusBadGateway, "service unavailable")
		return
	}
	defer resp.Body.Close()

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}
