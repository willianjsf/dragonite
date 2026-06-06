package usecase

import (
	"context"

	"github.com/caio-bernardo/dragonite/internal/domain"
)

type SearchFilter struct {
	IDCanais  []string // canais a procurar
	Term      string   //termo de pesquisa
	Limit     int      // limite de resultados
	NextToken string   // paginação
}

type UsuarioStorage interface {
	GetProfileByID(ctx context.Context, userID string) (*domain.Profile, error)
	UpdateProfile(ctx context.Context, profile domain.Profile) error
	SearchProfiles(ctx context.Context, filter SearchFilter) ([]domain.Profile, error)
	AddDirectMessage(ctx context.Context, senderID, receiverID string, roomID string) error
}

type CanalStorage interface {
	Create(ctx context.Context, roomID, userID string) (*domain.Canal, error)
	GetByID(ctx context.Context, canalID string) (*domain.Canal, error)
	GetUserJoinedRooms(ctx context.Context, userID string) ([]string, error)
	GetAllPublic(ctx context.Context, offset, limit int) ([]domain.Canal, error)
	UpdateForwardExtremities(ctx context.Context, canalID string, extremeties []string) error
	SaveAlias(ctx context.Context, roomID, fullAlias string) error
}

type EventoStorage interface {
	CheckNew(ctx context.Context, userID string) (bool, error)
	GetSince(ctx context.Context, userID string, since domain.SyncToken) ([]domain.Evento, domain.SyncToken, error)
	GetMaxGlobalStreamOrdering(ctx context.Context) (int64, error)

	SaveEvent(ctx context.Context, event *domain.Evento) error
}

// EventBus implementa um canal de transmissão de eventos publish-subscriber
type EventBus interface {
	Publish(ctx context.Context, canal_id string, event domain.Evento)
	PublishToUser(ctx context.Context, userID string, event domain.Evento)
	Subscribe(ctx context.Context, canal_id string) (<-chan *domain.Evento, func())
}

type DirectoryStorage interface {
	SearchDirectory(ctx context.Context, term string, limit, offset int) ([]domain.PublicRoomEntry, int, error)
}

// Executes operations inside a transaction. Commit if succeeds or rollback in failure
type WorkUnit interface {
	Execute(ctx context.Context, fn func(txCtx context.Context) error) error
}

type FederationService interface {
	QueueOutgoing(ctx context.Context, event domain.Evento) error
}
