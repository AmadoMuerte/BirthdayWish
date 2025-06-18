package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/AmadoMuerte/BirthdayWish/API/apps/wishlister/internal/models"
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

func (s *Storage) GetWish(ctx context.Context, userID, wishID int64) (models.Wishlist, error) {
	var wishlist models.Wishlist

	err := s.DB.NewSelect().
		Model(&wishlist).
		Where("id = ?", wishID).
		Where("user_id = ?", userID).
		Scan(ctx)
	if err != nil {
		return wishlist, err
	}

	return wishlist, nil
}

func (s *Storage) CheckWishExists(ctx context.Context, userID int64, itemID int64) (bool, error) {
	exists, err := s.DB.NewSelect().
		Model((*models.Wishlist)(nil)).
		Where("user_id = ?", userID).
		Where("id = ?", itemID).
		Exists(ctx)

	if err != nil {
		return false, fmt.Errorf("failed to check wishlist item existence: %w", err)
	}

	return exists, nil
}

func (s *Storage) AddToWishlist(ctx context.Context, item models.Wishlist) (*models.Wishlist, error) {
	now := time.Now()
	item.CreatedAt = now
	item.UpdatedAt = now

	_, err := s.DB.NewInsert().
		Model(&item).
		Returning("*").
		Exec(ctx, &item)
	if err != nil {
		return nil, fmt.Errorf("failed to add item to wishlist: %w", err)
	}

	return &item, nil
}

func (s *Storage) PartialUpdateWishItem(ctx context.Context, userID, itemID int64, updates map[string]any) (*models.Wishlist, error) {
	var updatedItem models.Wishlist
	query := s.DB.NewUpdate().
		Model(&updatedItem).
		Where("id = ?", itemID).
		Where("user_id = ?", userID).
		Set("updated_at = ?", time.Now()).
		Returning("*")

	query = query.Set("updated_at = ?", time.Now())

	if name, ok := updates["name"].(string); ok {
		query = query.Set("name = ?", name)
	}
	if link, ok := updates["link"].(string); ok {
		query = query.Set("link = ?", link)
	}
	if imgURL, ok := updates["image_url"].(string); ok {
		query = query.Set("image_url = ?", imgURL)
	}
	if price, ok := updates["price"].(float64); ok {
		query = query.Set("price = ?", price)
	}

	_, err := query.Exec(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to partially update wishlist item: %w", err)
	}

	return &updatedItem, nil
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
