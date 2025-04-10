package wishlist

import (
	"net/http"
	"strconv"

	"github.com/AmadoMuerte/BirthdayWish/API/internal/lib/response"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func (h *WishlistHandler) GetWishlist(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user_id_str := chi.URLParam(r, "user_id")
	user_id, err := strconv.ParseInt(user_id_str, 10, 64)
	if err != nil {
		h.log.Error("user id is empty")
		render.JSON(w, r, response.Error("user id is empty"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if wishlist, err := h.storage.GetWishlist(ctx, user_id); err != nil {
		h.log.Error("error getting wishlist", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		render.JSON(w, r, response.Error("Internal Server Error"))
		return
	} else {
		var response = []WishItem{}

		for _, wishItem := range wishlist {
			response = append(response, WishItem{
				ID:        wishItem.ID,
				UserID:    wishItem.UserID,
				Price:     wishItem.Price,
				Link:      wishItem.Link,
				Name:      wishItem.Name,
				CreatedAt: wishItem.CreatedAt,
				UpdatedAt: wishItem.UpdatedAt,
			})
		}

		render.JSON(w, r, response)
	}
}
