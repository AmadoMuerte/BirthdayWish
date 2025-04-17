package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/AmadoMuerte/BirthdayWish/API/internal/models"
)

func (s *Storage) GetWishlist(ctx context.Context, userID int64) ([]models.Wishlist, error) {
	var wishlist []models.Wishlist

	err := s.DB.NewSelect().
		Model(&wishlist).
		Where("user_id = ?", userID).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	return wishlist, nil
}

func (s *Storage) AddToWishlist(ctx context.Context, item models.Wishlist) error {
	now := time.Now()
	item.CreatedAt = now
	item.UpdatedAt = now

	_, err := s.DB.NewInsert().
		Model(&item).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to add item to wishlist: %w", err)
	}

	return nil
}
