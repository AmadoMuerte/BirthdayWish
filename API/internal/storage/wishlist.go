package storage

import (
	"context"

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
