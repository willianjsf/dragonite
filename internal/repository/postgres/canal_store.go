package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CanalStorage struct {
	db *pgxpool.Pool
}

func NewCanalStorage(db *pgxpool.Pool) *CanalStorage {
	return &CanalStorage{db: db}
}

func (c *CanalStorage) GetUserJoinedRooms(ctx context.Context, userID string) ([]string, error) {
	return nil, nil
}
