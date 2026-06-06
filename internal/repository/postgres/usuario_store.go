package postgres

import (
	"context"

	"github.com/caio-bernardo/dragonite/internal/domain"
	"github.com/caio-bernardo/dragonite/internal/usecase"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UsuarioStorage struct {
	db *pgxpool.Pool
}

func NewUsuarioStorage(db *pgxpool.Pool) *UsuarioStorage {
	return &UsuarioStorage{db: db}
}

func (u *UsuarioStorage) SearchProfiles(ctx context.Context, params usecase.SearchFilter) ([]domain.Profile, error) {
	return nil, nil
}
