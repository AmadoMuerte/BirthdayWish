package wish_handler

import (
	"log/slog"
	"net/http"

	"github.com/AmadoMuerte/BirthdayWish/API/apps/gateway/internal/client"
	gen "github.com/AmadoMuerte/BirthdayWish/API/apps/gateway/internal/gen"
	"github.com/AmadoMuerte/BirthdayWish/API/apps/gateway/internal/jwt"
	"github.com/AmadoMuerte/BirthdayWish/API/apps/gateway/internal/response"
	wishProto "github.com/AmadoMuerte/BirthdayWish/API/proto/wish"
)

type WishHandler struct {
	wishClient *client.WishlisterClient
	log        *slog.Logger
}

func NewWishHandler(wishClient *client.WishlisterClient, log *slog.Logger) *WishHandler {
	return &WishHandler{
		wishClient: wishClient,
		log:        log,
	}
}

func (h *WishHandler) GetWish(w http.ResponseWriter, r *http.Request, wishId int64) {
	claims, err := jwt.GetClaims(r.Context())
	if err != nil {
		h.log.Error("failed to get claims", "error", err)
		response.ErrorResponseJSON(w, r, http.StatusInternalServerError, "Authentication failed")
		return
	}

	resp, err := h.wishClient.GetWish(r.Context(), &wishProto.GetWishRequest{UserId: claims.UserID, WishId: wishId})
	if err != nil {
		h.log.Error("Wishlister service GetWish failed", "error", err)
		response.ErrorResponseJSON(w, r, http.StatusInternalServerError, "Failed to get wish")
		return
	}

	responseData := gen.WishResponse{
		Id:          &resp.Id,
		UserId:      &resp.UserId,
		Title:       &resp.Title,
		Description: &resp.Description,
		ImageUrl:    &resp.ImageUrl,
		Link:        &resp.Link,
		Price:       &resp.Price,
		Priority:    &resp.Priority,
	}

	response.SuccessResponse(w, r, http.StatusOK, responseData)
}
