package usecase

import (
	"context"
	"fmt"
	"time"

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

func (u *UsuarioService) Sync(ctx context.Context, since domain.SyncToken, timeout time.Duration) ([]domain.Evento, *domain.SyncToken, error) {

	userID := ctx.Value("userID").(string)
	// Lógica de Long-Polling
	if since.TimelinePosition != 0 {
		hasEvents, err := u.eventoStore.CheckNew(ctx, userID, since)
		if err != nil {
			return nil, nil, err
		}

		if !hasEvents && timeout > 0 {
			// sem eventos, long-polling
			ch := u.notifier.Subscribe(userID)
			defer u.notifier.Unsubscribe(userID, ch)

			select {
			case <-ch:
				// Novo evento, pode acessar o banco
			case <-time.After(timeout):
				// Deu timeout antes de um novo evento, cria novo token e retorna
				maxGlobal, _ := u.eventoStore.GetMaxGlobalStreamOrdering(ctx)
				if maxGlobal > since.TimelinePosition {
					since.TimelinePosition = maxGlobal
				}
				return nil, &since, types.ErrTimeout
			case <-ctx.Done():
				// o client se desconectou
				return nil, nil, types.ErrLooseConnection
			}
		}
	}

	// accesso ao banco
	events, newToken, err := u.eventoStore.GetSince(ctx, userID, since)
	if err != nil {
		return nil, nil, err
	}
	return events, newToken, nil
}
