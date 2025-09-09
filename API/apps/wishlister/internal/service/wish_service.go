package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/AmadoMuerte/BirthdayWish/API/apps/wishlister/internal/models"
	"github.com/AmadoMuerte/BirthdayWish/API/apps/wishlister/internal/storage"
	wishProto "github.com/AmadoMuerte/BirthdayWish/API/proto/wish"
	"github.com/google/uuid"
)

type WishService struct {
	storage *storage.Storage
	log     *slog.Logger
	wishProto.UnimplementedWishServiceServer
}

func NewWishService(storage *storage.Storage, log *slog.Logger) *WishService {
	return &WishService{
		storage: storage,
		log:     log,
	}
}

func (s *WishService) GetWish(ctx context.Context, req *wishProto.GetWishRequest) (*wishProto.WishResponse, error) {
	wish, err := s.storage.GetWish(ctx, req.UserId, req.WishId)
	if err != nil {
		s.log.Error("failed to get wish", "error", err, "user_id", req.UserId, "wish_id", req.WishId)
		return nil, err
	}

	return &wishProto.WishResponse{
		Id:          wish.ID,
		UserId:      wish.UserID,
		Title:       wish.Title,
		Description: wish.Description,
		ImageUrl:    wish.ImageURL,
		Link:        wish.Link,
		Price:       wish.Price,
		Priority:    wish.Priority,
		CreatedAt:   wish.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   wish.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (s *WishService) ListWishes(ctx context.Context, req *wishProto.ListWishesRequest) (*wishProto.WishListResponse, error) {
	wishes, err := s.storage.GetWishlist(ctx, req.UserId)
	if err != nil {
		s.log.Error("failed to get wishlist", "error", err, "user_id", req.UserId)
		return nil, err
	}

	var protoWishes []*wishProto.WishResponse
	for _, wish := range wishes {
		protoWishes = append(protoWishes, &wishProto.WishResponse{
			Id:          wish.ID,
			UserId:      wish.UserID,
			Title:       wish.Title,
			Description: wish.Description,
			ImageUrl:    wish.ImageURL,
			Link:        wish.Link,
			Price:       wish.Price,
			Priority:    wish.Priority,
			CreatedAt:   wish.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   wish.UpdatedAt.Format(time.RFC3339),
		})
	}

	return &wishProto.WishListResponse{
		Wishes:     protoWishes,
		TotalCount: int32(len(protoWishes)),
	}, nil
}

func (s *WishService) UpdateWish(ctx context.Context, req *wishProto.UpdateWishRequest) (*wishProto.WishResponse, error) {
	updateData := make(map[string]any)
	for k, v := range req.UpdateData {
		updateData[k] = v
	}

	// TODO: Upload image to S3
	if updateData["image_base64"] != nil {
		updateData["image_url"] = "https://www.pinterest.com/pin/63261569757704014/"
	}

	wish, err := s.storage.PartialUpdateWishItem(ctx, req.UserId, req.WishId, updateData)
	if err != nil {
		s.log.Error("failed to update wish", "error", err, "user_id", req.UserId, "wish_id", req.WishId)
		return nil, err
	}

	return &wishProto.WishResponse{
		Id:          wish.ID,
		UserId:      wish.UserID,
		Title:       wish.Title,
		Description: wish.Description,
		ImageUrl:    wish.ImageURL,
		Link:        wish.Link,
		Price:       wish.Price,
		Priority:    wish.Priority,
		CreatedAt:   wish.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   wish.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (s *WishService) CreateWish(ctx context.Context, req *wishProto.CreateWishRequest) (*wishProto.WishResponse, error) {
	// TODO: Upload image to S3

	wish := &models.Wish{
		UserID:      req.UserId,
		Title:       req.Title,
		Description: req.Description,
		ImageURL:    "https://www.pinterest.com/pin/63261569757704014/",
		Link:        req.Link,
		Price:       req.Price,
		Priority:    req.Priority,
	}

	wish, err := s.storage.AddToWishlist(ctx, *wish)
	if err != nil {
		s.log.Error("failed to create wish", "error", err, "user_id", req.UserId)
		return nil, err
	}

	return &wishProto.WishResponse{
		Id:          wish.ID,
		UserId:      wish.UserID,
		Title:       wish.Title,
		Description: wish.Description,
		ImageUrl:    wish.ImageURL,
		Link:        wish.Link,
		Price:       wish.Price,
		Priority:    wish.Priority,
		CreatedAt:   wish.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   wish.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (s *WishService) DeleteWish(ctx context.Context, req *wishProto.DeleteWishRequest) (*wishProto.EmptyResponse, error) {
	err := s.storage.RemoveFromWishlist(ctx, req.UserId, req.WishId)
	if err != nil {
		s.log.Error("failed to delete wish", "error", err, "user_id", req.UserId, "wish_id", req.WishId)
		return nil, err
	}

	return &wishProto.EmptyResponse{}, nil
}

func (s *WishService) CreateShareLink(ctx context.Context, req *wishProto.CreateShareLinkRequest) (*wishProto.ShareLinkResponse, error) {
	shareLink := &models.ShareLink{
		UserID:    req.UserId,
		ExpiresAt: time.Now().Add(time.Duration(req.ExpiresIn) * time.Second),
		Token:     uuid.New().String(),
	}

	shareLink, err := s.storage.AddShareLink(ctx, *shareLink)
	if err != nil {
		s.log.Error("failed to create share link", "error", err, "user_id", req.UserId)
		return nil, err
	}

	return &wishProto.ShareLinkResponse{
		ShareToken: shareLink.Token,
		ExpiresAt:  shareLink.ExpiresAt.Format(time.RFC3339),
	}, nil
}

func (s *WishService) GetShareLinks(ctx context.Context, req *wishProto.GetShareLinksRequest) (*wishProto.ShareLinksResponse, error) {
	shareLinks, err := s.storage.GetShareLinks(ctx, req.UserId)
	if err != nil {
		s.log.Error("failed to get share links", "error", err, "user_id", req.UserId)
		return nil, err
	}

	var protoShareLinks []*wishProto.ShareLinkResponse
	for _, shareLink := range shareLinks {
		protoShareLinks = append(protoShareLinks, &wishProto.ShareLinkResponse{
			ShareToken: shareLink.Token,
			ExpiresAt:  shareLink.ExpiresAt.Format(time.RFC3339),
		})
	}

	return &wishProto.ShareLinksResponse{
		Tokens:     protoShareLinks,
		TotalCount: int32(len(protoShareLinks)),
	}, nil
}

func (s *WishService) DeleteShareLink(ctx context.Context, req *wishProto.DeleteShareLinkRequest) (*wishProto.EmptyResponse, error) {
	err := s.storage.DeleteShareLink(ctx, req.UserId, req.ShareToken)
	if err != nil {
		s.log.Error("failed to delete share link", "error", err, "user_id", req.UserId, "share_token", req.ShareToken)
		return nil, err
	}
	return &wishProto.EmptyResponse{}, nil
}

func (s *WishService) GetSharedWishes(ctx context.Context, req *wishProto.GetSharedRequest) (*wishProto.WishListResponse, error) {
	shareLink, err := s.storage.GetShareLink(ctx, req.ShareToken)
	if err != nil {
		s.log.Error("failed to get share link", "error", err, "share_token", req.ShareToken)
		return nil, err
	}

	if shareLink.ExpiresAt.Before(time.Now()) {
		s.log.Error("share link expired", "error", err, "share_token", req.ShareToken)
		return nil, err
	}

	return s.ListWishes(ctx, &wishProto.ListWishesRequest{
		UserId: shareLink.UserID,
	})
}
