package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/AmadoMuerte/BirthdayWish/API/apps/gateway/internal/models"
)

func (s *Storage) CreateUser(ctx context.Context, user *models.User) error {
	_, err := s.DB.NewInsert().Model(user).Exec(ctx)
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

func (s *Storage) UserExists(ctx context.Context, username, email string) (bool, error) {
	exists, err := s.DB.NewSelect().
		Model((*models.User)(nil)).
		Where("email = ? OR name = ?", email, username).
		Exists(ctx)
	if err != nil {
		return false, fmt.Errorf("check user exists: %w", err)
	}
	return exists, nil
}

func (s *Storage) UserExistsByID(ctx context.Context, userID int64) (bool, error) {
	exists, err := s.DB.NewSelect().
		Model((*models.User)(nil)).
		Where("id = ?", userID).
		Exists(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to check user existence: %w", err)
	}
	return exists, nil
}

func (s *Storage) GetUserByUsername(ctx context.Context, name string) (models.User, error) {
	var user models.User

	err := s.DB.NewSelect().
		Model(&user).
		Where("name = ?", name).
		Scan(ctx)

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (s *Storage) GetUserIDByToken(ctx context.Context, token string) (int64, error) {
	var userID int64

	err := s.DB.NewSelect().
		Column("user_id").
		Table("share_wishlist_access").
		Where("access_token = ?", token).
		Limit(1).
		Scan(ctx, &userID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("token not found")
		}
		return 0, fmt.Errorf("failed to get user ID: %w", err)
	}

	return userID, nil
}
