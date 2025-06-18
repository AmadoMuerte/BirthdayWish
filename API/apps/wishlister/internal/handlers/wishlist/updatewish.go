package wishlist

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/AmadoMuerte/BirthdayWish/API/pkg/response"
	"github.com/go-chi/chi/v5"
)

func (h *WishlistHandler) UpdateWish(w http.ResponseWriter, r *http.Request) {
	op := "wishlist/UpdateWish"
	ctx := r.Context()

	wishID, err := strconv.ParseInt(chi.URLParam(r, "wish_id"), 10, 64)
	if err != nil {
		h.log.Error(op+": invalid wish_id", "error", err)
		response.ErrorResponseJSON(w, r, http.StatusBadRequest, "Invalid request body")
		return
	}

	userID, err := strconv.ParseInt(chi.URLParam(r, "user_id"), 10, 64)
	if err != nil {
		h.log.Error(op+": invalid user_id", "error", err)
		response.ErrorResponseJSON(w, r, http.StatusBadRequest, "Invalid request body")
		return
	}

	exists, err := h.storage.CheckWishExists(ctx, userID, wishID)
	if err != nil {
		h.log.Error(op+": failed to check wish existence", "error", err)
		response.ErrorResponseJSON(w, r, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	if !exists {
		h.log.Error(op+": wish not found", "error", err)
		response.ErrorResponseJSON(w, r, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	var updateData map[string]any
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		h.log.Error(op+": failed to decode request body", "error", err)
		response.ErrorResponseJSON(w, r, http.StatusBadRequest, "Invalid request body")
		return
	}

	updateData["updated_at"] = time.Now()

	wish, err := h.storage.PartialUpdateWishItem(ctx, userID, wishID, updateData)
	if err != nil {
		h.log.Error(op+": failed to update wish", "error", err)
		response.ErrorResponseJSON(w, r, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	response.SuccessResponse(w, r, http.StatusOK, wish)
}
