package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/AmadoMuerte/BirthdayWish/API/apps/gateway/internal/config"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

type Storage struct {
	DB *bun.DB
}

func NewStorage(cfg *config.Config) (*Storage, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&timeout=5s",
		cfg.DB.User,
		cfg.DB.Pass,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.Name,
	)

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	err := ApplyMigrations(sqldb, cfg.DB.Name)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to apply migrations: %w", err)
	}

	return &Storage{DB: db}, nil
}

func (s *Storage) Close() error {
	if s == nil || s.DB == nil {
		return nil
	}
	return s.DB.Close()
}
