package storage

import (
	"context"
	"fmt"

	"github.com/AmadoMuerte/BirthdayWish/API/internal/models"
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
