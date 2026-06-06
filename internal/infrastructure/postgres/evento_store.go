package postgres

import "github.com/jackc/pgx/v5/pgxpool"

type EventoStorage struct {
	db *pgxpool.Pool
}

func NewEventoStorage(db *pgxpool.Pool) *EventoStorage {
	return &EventoStorage{db: db}
}
