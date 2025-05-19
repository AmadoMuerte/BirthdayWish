package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/AmadoMuerte/BirthdayWish/API/services/wishlist/internal/models"
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

func (s *Storage) RemoveFromWishlist(ctx context.Context, wishID, userID int64) error {
	var wishlist models.Wishlist

	res, err := s.DB.NewDelete().
		Model(&wishlist).
		Where("user_id = ?", userID).
		Where("id = ?", wishID).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to remove item from wishlist: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("item not found in wishlist or already deleted")
	}

	return nil
}
