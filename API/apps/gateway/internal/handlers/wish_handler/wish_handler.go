package wish_handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/AmadoMuerte/BirthdayWish/API/apps/gateway/internal/client"
	gen "github.com/AmadoMuerte/BirthdayWish/API/apps/gateway/internal/gen"
	"github.com/AmadoMuerte/BirthdayWish/API/apps/gateway/internal/jwt"
	"github.com/AmadoMuerte/BirthdayWish/API/apps/gateway/internal/response"
	wishProto "github.com/AmadoMuerte/BirthdayWish/API/shared/proto/wish"
	"github.com/google/uuid"
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

func convertTime(timeStr string) (time.Time, error) {
	time, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return time, err
	}
	return time, nil
}

func (h *WishHandler) GetWish(w http.ResponseWriter, r *http.Request, wishId int64) {
	claims, err := jwt.GetClaims(r.Context())
	if err != nil {
		h.log.Error("failed to get claims", "error", err)
		response.ErrorResponseJSON(w, r, http.StatusInternalServerError, "Internal server error")
		return
	}

	resp, err := h.wishClient.GetWish(r.Context(), &wishProto.GetWishRequest{UserId: claims.UserID, WishId: wishId})
	if err != nil {
		h.log.Error("Wishlister service GetWish failed", "error", err)
		response.ErrorResponseJSON(w, r, http.StatusBadRequest, "Invalid request")
		return
	}

	createdAt, err := convertTime(resp.CreatedAt)
	if err != nil {
		h.log.Error("failed to convert time", "error", err)
		response.ErrorResponseJSON(w, r, http.StatusInternalServerError, "Internal server error")
		return
	}
	updatedAt, err := convertTime(resp.UpdatedAt)
	if err != nil {
		h.log.Error("failed to convert time", "error", err)
		response.ErrorResponseJSON(w, r, http.StatusInternalServerError, "Internal server error")
		return
	}

	responseData := gen.WishResponse{
		Id:          resp.Id,
		UserId:      resp.UserId,
		Title:       resp.Title,
		Description: &resp.Description,
		ImageUrl:    &resp.ImageUrl,
		Link:        &resp.Link,
		Price:       &resp.Price,
		Priority:    &resp.Priority,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}

	response.SuccessResponse(w, r, http.StatusOK, responseData)
}

func (h *WishHandler) ListWishes(w http.ResponseWriter, r *http.Request) {
	claims, err := jwt.GetClaims(r.Context())
	if err != nil {
		h.log.Error("failed to get claims", "error", err)
		response.ErrorResponseJSON(w, r, http.StatusInternalServerError, "Internal server error")
		return
	}

	resp, err := h.wishClient.ListWishes(r.Context(), &wishProto.ListWishesRequest{UserId: claims.UserID})
	if err != nil {
		h.log.Error("Wishlister service ListWishes failed", "error", err)
		response.ErrorResponseJSON(w, r, http.StatusBadRequest, "Invalid request")
		return
	}

	var wishes []gen.WishResponse
	for _, wish := range resp.Wishes {
		createdAt, err := convertTime(wish.CreatedAt)
		if err != nil {
			h.log.Error("failed to convert time", "error", err)
			response.ErrorResponseJSON(w, r, http.StatusInternalServerError, "Internal server error")
			return
		}
		updatedAt, err := convertTime(wish.UpdatedAt)
		if err != nil {
			h.log.Error("failed to convert time", "error", err)
			response.ErrorResponseJSON(w, r, http.StatusInternalServerError, "Internal server error")
			return
		}
		wishes = append(wishes, gen.WishResponse{
			Id:          wish.Id,
			UserId:      wish.UserId,
			Title:       wish.Title,
			Description: &wish.Description,
			ImageUrl:    &wish.ImageUrl,
			Link:        &wish.Link,
			Price:       &wish.Price,
			Priority:    &wish.Priority,
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
		})
	}

	responseData := gen.WishListResponse{
		Wishes:     wishes,
		TotalCount: resp.TotalCount,
	}

	response.SuccessResponse(w, r, http.StatusOK, responseData)
}

func (h *WishHandler) CreateWish(w http.ResponseWriter, r *http.Request) {
	claims, err := jwt.GetClaims(r.Context())
	if err != nil {
		h.log.Error("failed to get claims", "error", err)
		response.ErrorResponseJSON(w, r, http.StatusInternalServerError, "Internal server error")
		return
	}

	var req gen.CreateWishRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("failed to decode request body", "error", err)
		response.ErrorResponseJSON(w, r, http.StatusBadRequest, "Invalid request")
		return
	}

	resp, err := h.wishClient.CreateWish(r.Context(), &wishProto.CreateWishRequest{
		UserId:      claims.UserID,
		Title:       req.Title,
		Description: *req.Description,
		ImageBase64: *req.ImageBase64,
		Link:        *req.Link,
		Price:       *req.Price,
		Priority:    *req.Priority,
	})
	if err != nil {
		h.log.Error("Wishlister service CreateWish failed", "error", err)
		response.ErrorResponseJSON(w, r, http.StatusBadRequest, "Invalid request")
		return
	}

	createdAt, err := convertTime(resp.CreatedAt)
	if err != nil {
		h.log.Error("failed to convert time", "error", err)
		response.ErrorResponseJSON(w, r, http.StatusInternalServerError, "Internal server error")
		return
	}
	updatedAt, err := convertTime(resp.UpdatedAt)
	if err != nil {
		h.log.Error("failed to convert time", "error", err)
		response.ErrorResponseJSON(w, r, http.StatusInternalServerError, "Internal server error")
		return
	}
	responseData := gen.WishResponse{
		Id:          resp.Id,
		UserId:      resp.UserId,
		Title:       resp.Title,
		Description: &resp.Description,
		ImageUrl:    &resp.ImageUrl,
		Link:        &resp.Link,
		Price:       &resp.Price,
		Priority:    &resp.Priority,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}

	response.SuccessResponse(w, r, http.StatusOK, responseData)
}

func (h *WishHandler) UpdateWish(w http.ResponseWriter, r *http.Request, wishId int64) {
	claims, err := jwt.GetClaims(r.Context())
	if err != nil {
		h.log.Error("failed to get claims", "error", err)
		response.ErrorResponseJSON(w, r, http.StatusInternalServerError, "Internal server error")
		return
	}

	var req gen.UpdateWishRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("failed to decode request body", "error", err)
		response.ErrorResponseJSON(w, r, http.StatusBadRequest, "Invalid request")
		return
	}

	updateData := map[string]string{}
	if req.Title != nil {
		updateData["title"] = *req.Title
	}
	if req.Description != nil {
		updateData["description"] = *req.Description
	}
	if req.ImageBase64 != nil {
		updateData["image_base64"] = *req.ImageBase64
	}
	if req.Link != nil {
		updateData["link"] = *req.Link
	}
	if req.Price != nil {
		updateData["price"] = strconv.FormatFloat(*req.Price, 'f', -1, 64)
	}
	if req.Priority != nil {
		updateData["priority"] = strconv.FormatInt(int64(*req.Priority), 10)
	}

	resp, err := h.wishClient.UpdateWish(r.Context(), &wishProto.UpdateWishRequest{
		UserId:     claims.UserID,
		WishId:     wishId,
		UpdateData: updateData,
	})
	if err != nil {
		h.log.Error("failed to update wish", "error", err)
		response.ErrorResponseJSON(w, r, http.StatusInternalServerError, "Internal server error")
		return
	}

	createdAt, err := convertTime(resp.CreatedAt)
	if err != nil {
		h.log.Error("failed to convert time", "error", err)
		response.ErrorResponseJSON(w, r, http.StatusInternalServerError, "Internal server error")
		return
	}
	updatedAt, err := convertTime(resp.UpdatedAt)
	if err != nil {
		h.log.Error("failed to convert time", "error", err)
		response.ErrorResponseJSON(w, r, http.StatusInternalServerError, "Internal server error")
		return
	}
	responseData := gen.WishResponse{
		Id:          resp.Id,
		UserId:      resp.UserId,
		Title:       resp.Title,
		Description: &resp.Description,
		ImageUrl:    &resp.ImageUrl,
		Link:        &resp.Link,
		Price:       &resp.Price,
		Priority:    &resp.Priority,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}

	response.SuccessResponse(w, r, http.StatusOK, responseData)
}

func (h *WishHandler) DeleteWish(w http.ResponseWriter, r *http.Request, wishId int64) {
	claims, err := jwt.GetClaims(r.Context())
	if err != nil {
		h.log.Error("failed to get claims", "error", err)
		response.ErrorResponseJSON(w, r, http.StatusInternalServerError, "Internal server error")
		return
	}

	_, err = h.wishClient.DeleteWish(r.Context(), &wishProto.DeleteWishRequest{
		UserId: claims.UserID,
		WishId: wishId,
	})
	if err != nil {
		h.log.Error("failed to delete wish", "error", err)
		response.ErrorResponseJSON(w, r, http.StatusInternalServerError, "Internal server error")
		return
	}

	responseData := gen.WishDeletedResponse{
		Message: "Wish deleted successfully",
	}

	response.SuccessResponse(w, r, http.StatusOK, responseData)
}

func (h *WishHandler) DeleteShareLink(w http.ResponseWriter, r *http.Request) {
	var req gen.DeleteShareLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("failed to decode request body", "error", err)
		response.ErrorResponseJSON(w, r, http.StatusBadRequest, "Invalid request")
		return
	}
	err := uuid.Validate(req.ShareToken)
	if err != nil {
		h.log.Error("invalid share token", "error", err)
		response.ErrorResponseJSON(w, r, http.StatusBadRequest, "Invalid request")
		return
	}

	claims, err := jwt.GetClaims(r.Context())
	if err != nil {
		h.log.Error("failed to get claims", "error", err)
		response.ErrorResponseJSON(w, r, http.StatusInternalServerError, "Internal server error")
		return
	}

	_, err = h.wishClient.DeleteShareLink(r.Context(), &wishProto.DeleteShareLinkRequest{
		UserId:     claims.UserID,
		ShareToken: req.ShareToken,
	})
	if err != nil {
		h.log.Error("failed to delete share link", "error", err)
		response.ErrorResponseJSON(w, r, http.StatusInternalServerError, "Internal server error")
		return
	}

	responseData := gen.WishDeletedResponse{
		Message: "Wish deleted successfully",
	}

	response.SuccessResponse(w, r, http.StatusOK, responseData)
}

func (h *WishHandler) CreateShareLink(w http.ResponseWriter, r *http.Request) {
	claims, err := jwt.GetClaims(r.Context())
	if err != nil {
		h.log.Error("failed to get claims", "error", err)
		response.ErrorResponseJSON(w, r, http.StatusInternalServerError, "Internal server error")
		return
	}

	var req gen.CreateShareLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("failed to decode request body", "error", err)
		response.ErrorResponseJSON(w, r, http.StatusBadRequest, "Invalid request")
		return
	}

	resp, err := h.wishClient.CreateShareLink(r.Context(), &wishProto.CreateShareLinkRequest{
		UserId:    claims.UserID,
		ExpiresIn: req.ExpiresIn,
	})
	if err != nil {
		h.log.Error("failed to create share link", "error", err)
		response.ErrorResponseJSON(w, r, http.StatusInternalServerError, "Internal server error")
		return
	}

	expiresAt, err := convertTime(resp.ExpiresAt)
	if err != nil {
		h.log.Error("failed to convert time", "error", err)
		response.ErrorResponseJSON(w, r, http.StatusInternalServerError, "Internal server error")
		return
	}

	responseData := gen.ShareLinkResponse{
		ShareToken: &resp.ShareToken,
		ExpiresAt:  expiresAt,
	}

	response.SuccessResponse(w, r, http.StatusOK, responseData)
}

func (h *WishHandler) GetShareLinks(w http.ResponseWriter, r *http.Request) {
	claims, err := jwt.GetClaims(r.Context())
	if err != nil {
		h.log.Error("failed to get claims", "error", err)
		response.ErrorResponseJSON(w, r, http.StatusInternalServerError, "Internal server error")
		return
	}

	resp, err := h.wishClient.GetShareLinks(r.Context(), &wishProto.GetShareLinksRequest{
		UserId: claims.UserID,
	})
	if err != nil {
		h.log.Error("failed to get share links", "error", err)
		response.ErrorResponseJSON(w, r, http.StatusInternalServerError, "Internal server error")
		return
	}

	var shareLinks []gen.ShareLinkResponse
	for _, shareLink := range resp.Tokens {
		expiresAt, err := convertTime(shareLink.ExpiresAt)
		if err != nil {
			h.log.Error("failed to convert time", "error", err)
			response.ErrorResponseJSON(w, r, http.StatusInternalServerError, "Internal server error")
			return
		}
		shareLinks = append(shareLinks, gen.ShareLinkResponse{
			ShareToken: &shareLink.ShareToken,
			ExpiresAt:  expiresAt,
		})
	}

	responseData := gen.ShareLinksResponse{
		Tokens:     shareLinks,
		TotalCount: resp.TotalCount,
	}

	response.SuccessResponse(w, r, http.StatusOK, responseData)
}

func (h *WishHandler) GetSharedWishes(w http.ResponseWriter, r *http.Request) {
	var req gen.GetSharedWishesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.log.Error("failed to decode request body", "error", err)
		response.ErrorResponseJSON(w, r, http.StatusBadRequest, "Invalid request")
		return
	}

	err := uuid.Validate(req.ShareToken)
	if err != nil {
		h.log.Error("invalid share token", "error", err)
		response.ErrorResponseJSON(w, r, http.StatusBadRequest, "Invalid request")
		return
	}

	resp, err := h.wishClient.GetSharedWishes(r.Context(), &wishProto.GetSharedRequest{
		ShareToken: req.ShareToken,
	})
	if err != nil {
		h.log.Error("failed to get shared wishes", "error", err)
		response.ErrorResponseJSON(w, r, http.StatusInternalServerError, "Internal server error")
		return
	}

	var wishes []gen.WishResponse
	for _, wish := range resp.Wishes {
		createdAt, err := convertTime(wish.CreatedAt)
		if err != nil {
			h.log.Error("failed to convert time", "error", err)
			response.ErrorResponseJSON(w, r, http.StatusInternalServerError, "Internal server error")
			return
		}
		updatedAt, err := convertTime(wish.UpdatedAt)
		if err != nil {
			h.log.Error("failed to convert time", "error", err)
			response.ErrorResponseJSON(w, r, http.StatusInternalServerError, "Internal server error")
			return
		}
		wishes = append(wishes, gen.WishResponse{
			Id:          wish.Id,
			UserId:      wish.UserId,
			Title:       wish.Title,
			Description: &wish.Description,
			ImageUrl:    &wish.ImageUrl,
			Link:        &wish.Link,
			Price:       &wish.Price,
			Priority:    &wish.Priority,
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
		})
	}

	responseData := gen.WishListResponse{
		Wishes:     wishes,
		TotalCount: resp.TotalCount,
	}

	response.SuccessResponse(w, r, http.StatusOK, responseData)
}
