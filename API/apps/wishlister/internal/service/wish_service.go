package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/AmadoMuerte/BirthdayWish/API/apps/wishlister/internal/storage"
	wishProto "github.com/AmadoMuerte/BirthdayWish/API/proto/wish"
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

// CreateWish, DeleteWish, CreateShareLink, GetSharedWishes
