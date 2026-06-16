package usecase

import (
	"context"
	"fmt"

	"github.com/caio-bernardo/dragonite/internal/domain"
	"github.com/caio-bernardo/dragonite/internal/domain/types"
)

type DirectoryService struct {
	dirStore   DirectoryStorage
	userStore  UsuarioStorage
	canalStore CanalStorage
}

func NewDirectoryService(userStore UsuarioStorage, canalStore CanalStorage) *DirectoryService {
	return &DirectoryService{
		userStore:  userStore,
		canalStore: canalStore,
	}
}

func (s *DirectoryService) ListPublic(ctx context.Context, term string, limit int, offset int) (*domain.PublicRoomsChunck, error) {
	if limit < 0 || limit > 100 {
		limit = 50
	}

	entries, totalCount, err := s.dirStore.SearchDirectory(ctx, term, limit, offset)
	if err != nil {
		return nil, err
	}

	response := domain.PublicRoomsChunck{
		Chunk:                  entries,
		TotalRoomCountEstimate: totalCount,
		NextBatch:              fmt.Sprintf("%d", len(entries)+offset),
		PrevBatch:              fmt.Sprintf("%d", offset),
	}

	return &response, nil
}

func (s *DirectoryService) SearchProfiles(ctx context.Context, query string, limit int) ([]domain.Profile, error) {
	userID := ctx.Value(types.UserIDKey).(string)
	if query == "" {
		return nil, types.ErrInvalidSearchTerm
	}

	allowedRooms, err := s.canalStore.GetUserJoinedRooms(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("Failed to verify user membership: %w", err)
	}
	// +1 para saber se há mais resultados
	filter := SearchFilter{
		IDCanais:  allowedRooms,
		Term:      query,
		Limit:     limit + 1,
		NextToken: "",
	}
	return s.userStore.SearchProfiles(ctx, filter)
}
