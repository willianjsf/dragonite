package postgres

import (
	"context"
	"testing"
	"time"

	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
)

func TestPostgresPingDB(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	container, err := tcpostgres.Run(ctx,
		"postgres:16-alpine",
		tcpostgres.WithDatabase("dragonite"),
		tcpostgres.WithUsername("dragonite"),
		tcpostgres.WithPassword("dragonite"),
	)
	if err != nil {
		t.Fatalf("start postgres container: %v", err)
	}
	defer container.Terminate(ctx)

	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("connection string: %v", err)
	}

	var poolClose func()
	var store *PostgresStorage
	deadline := time.Now().Add(15 * time.Second)
	for {
		pool, err := ConnectBD(ctx, connStr)
		if err == nil {
			poolClose = pool.Close
			store = NewPostgresStorage(pool)
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("connect db: %v", err)
		}
		time.Sleep(500 * time.Millisecond)
	}
	defer poolClose()

	status := store.PingDB()

	if status["status"] != "up" {
		t.Fatalf("expected status up, got %+v", status)
	}
}
