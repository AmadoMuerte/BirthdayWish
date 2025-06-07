package wishlist

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/AmadoMuerte/BirthdayWish/API/internal/pkg/response"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)


func (h *WishlistHandler) PartialUpdateWish(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()

    wishID, err := strconv.ParseInt(chi.URLParam(r, "wish_id"), 10, 64)
    if err != nil {
        h.log.Error("invalid wish_id")
        render.JSON(w, r, response.Error("invalid wish_id"))
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    userID, err := strconv.ParseInt(chi.URLParam(r, "user_id"), 10, 64)
    if err != nil {
        h.log.Error("invalid user_id")
        render.JSON(w, r, response.Error("invalid user_id"))
        w.WriteHeader(http.StatusBadRequest)
        return
    }

    exists, err := h.storage.CheckWishExists(ctx, userID, wishID)
    if err != nil {
        h.log.Error("failed to check wish existence", "error", err)
        w.WriteHeader(http.StatusInternalServerError)
        render.JSON(w, r, response.Error("failed to check wish"))
        return
    }
    if !exists {
        h.log.Error("wish not found")
        w.WriteHeader(http.StatusNotFound)
        render.JSON(w, r, response.Error("wish not found"))
        return
    }

    var updateData map[string]any
    if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
        h.log.Error("failed to decode request body", "error", err)
        w.WriteHeader(http.StatusBadRequest)
        render.JSON(w, r, response.Error("invalid request body"))
        return
    }

    updateData["updated_at"] = time.Now()

    if err := h.storage.PartialUpdateWishItem(ctx, userID, wishID, updateData); err != nil {
        h.log.Error("failed to update wish", "error", err)
        w.WriteHeader(http.StatusInternalServerError)
        render.JSON(w, r, response.Error("failed to update wish"))
        return
    }

    w.WriteHeader(http.StatusOK)
    render.JSON(w, r, response.Success("wish updated successfully"))
}