package wishlist

import (
	"net/http"
	"strconv"

	"github.com/AmadoMuerte/BirthdayWish/API/pkg/apimodels"
	"github.com/AmadoMuerte/BirthdayWish/API/pkg/jwt"
	"github.com/AmadoMuerte/BirthdayWish/API/pkg/redis"
	"github.com/AmadoMuerte/BirthdayWish/API/pkg/response"
	"github.com/go-chi/chi/v5"
)

func (h *WishlistHandler) DeleteWish(w http.ResponseWriter, r *http.Request) {
	op := "wishlist/DeleteWish"
	ctx := r.Context()

	wishID, err := strconv.ParseInt(chi.URLParam(r, "wish_id"), 10, 64)
	if err != nil {
		h.log.Error(op+": wish_id is empty", "error", err)
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
		h.log.Error(op+": user does not exist", "error", err)
		response.ErrorResponseJSON(w, r, http.StatusBadRequest, "invalid request body")
		return
	}

	task := apimodels.DeleteWishTask{
		WishID: wishID,
		UserID: claims.UserID,
	}

	if err := h.RedisClient.PushToQueue(ctx, redis.TaskDeleteWish, task); err != nil {
		h.log.Error(op+": failed to queue task", "error", err)
		response.ErrorResponseJSON(w, r, http.StatusServiceUnavailable, "service unavailable")
		return
	}

	response.SuccessResponse(w, r, http.StatusAccepted, map[string]string{
		"status":  "queued",
		"message": "delete request is being processed",
	})
}
