package postgres

import (
	"context"

	"github.com/caio-bernardo/dragonite/internal/domain"
	"github.com/caio-bernardo/dragonite/internal/usecase"
)

func (u *PostgresStorage) CreateUsuarioAndProfile(ctx context.Context, userProps domain.Usuario) (*domain.Usuario, error) {

}

func (u *PostgresStorage) GetUsuarioByID(ctx context.Context, userID string) (*domain.Usuario, error) {

}

func (u *PostgresStorage) GetProfileByID(ctx context.Context, userID string) (*domain.Profile, error) {

}

func (u *PostgresStorage) UpdateProfile(ctx context.Context, profile domain.Profile) error {

}

func (u *PostgresStorage) SearchProfiles(ctx context.Context, params usecase.SearchFilter) ([]domain.Profile, error) {
	return nil, nil
}

func (u *PostgresStorage) AddDirectMessage(ctx context.Context, senderID, receiverID string, roomID string) error {

}
