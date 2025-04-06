package storage

import (
	"context"

	"github.com/AmadoMuerte/BirthdayWish/API/internal/models"
)

func (s *Storage) CreateUser(user *models.User) error {
	ctx := context.Background()
	_, err := s.DB.NewInsert().Model(user).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}
