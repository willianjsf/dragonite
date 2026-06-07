package postgres

import (
	"context"

	"github.com/caio-bernardo/dragonite/internal/domain"
	"github.com/caio-bernardo/dragonite/internal/usecase"
)

func (u *PostgresStorage) SearchProfiles(ctx context.Context, params usecase.SearchFilter) ([]domain.Profile, error) {
	return nil, nil
}
