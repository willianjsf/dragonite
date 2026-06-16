package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/caio-bernardo/dragonite/internal/domain"
	"github.com/caio-bernardo/dragonite/internal/util"
)

type SyncService struct {
	eventStore EventoStorage
	notifier   Notifier
}

func NewSyncService(eventStore EventoStorage, notifier Notifier) *SyncService {
	return &SyncService{
		eventStore: eventStore,
		notifier:   notifier,
	}
}

func (s *SyncService) SyncClient(ctx context.Context, userID string, since domain.SyncToken, timeout time.Duration) ([]domain.Evento, domain.SyncToken, error) {

	eventos, err := s.eventStore.GetSince(ctx, userID, since)
	if err != nil {
		return nil, since, err
	}

	// sem long-polling, envia eventos
	if len(eventos) > 0 || timeout <= 0 {
		return eventos, util.GenerateNextSinceToken(since, eventos), nil
	}

	// Long-polling.
	pollCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// espera por uma notificação do banco
	err = s.notifier.WaitForEvents(pollCtx, userID)
	if err != nil {
		// deu timeout, retornamos apenas uma lista vazia.
		if errors.Is(err, context.DeadlineExceeded) {
			return []domain.Evento{}, since, nil
		}
		return nil, since, err
	}

	eventos, err = s.eventStore.GetSince(pollCtx, userID, since)
	if err != nil {
		return nil, since, err
	}

	return eventos, util.GenerateNextSinceToken(since, eventos), nil
}
