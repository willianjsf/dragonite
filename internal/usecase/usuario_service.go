package usecase

import (
	"context"
	"fmt"

	"github.com/caio-bernardo/dragonite/internal/domain"
	"github.com/caio-bernardo/dragonite/internal/domain/types"
	"github.com/jackc/pgx/v5"
)

type UsuarioService struct {
	eventoStore  EventoStorage
	usuarioStore UsuarioStorage
	canalStore   CanalStorage
}

func NewUsuarioService(eventoStore EventoStorage, usuarioStore UsuarioStorage, canalStore CanalStorage) *UsuarioService {
	return &UsuarioService{
		eventoStore:  eventoStore,
		usuarioStore: usuarioStore,
		canalStore:   canalStore,
	}
}

func (u *UsuarioService) GetProfileByID(ctx context.Context, userID string) (*domain.Profile, error) {
	usuario, err := u.usuarioStore.GetProfileByID(ctx, userID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, types.ErrNotFound
		}
		return nil, err
	}
	return usuario, nil
}

func (u *UsuarioService) SearchProfiles(ctx context.Context, term string, limit int) ([]domain.Profile, error) {
	userID := ctx.Value(types.UserIDKey).(string)
	if term == "" {
		return nil, types.ErrInvalidSearchTerm
	}

	allowedRooms, err := u.canalStore.GetUserJoinedRooms(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("Failed to verify user membership: %w", err)
	}
	filter := SearchFilter{
		IDCanais:  allowedRooms,
		Term:      term,
		Limit:     limit,
		NextToken: "",
	}
	return u.usuarioStore.SearchProfiles(ctx, filter)
}
