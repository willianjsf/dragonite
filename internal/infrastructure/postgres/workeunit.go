package postgres

import "context"

// NOTE: used to avoid type conflicts
type contextKey string

const txKey contextKey = "tx"

// Executes a function under a transaction
func (s *PostgresStorage) Execute(
	ctx context.Context,
	fn func(txCtx context.Context) error,
) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	txCtx := context.WithValue(ctx, txKey, tx)

	err = fn(txCtx)
	if err != nil {
		return err
	}
	return tx.Commit(txCtx)
}
