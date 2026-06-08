package redis_infra

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/caio-bernardo/dragonite/internal/infrastructure"
	"github.com/redis/go-redis/v9"
)

type idempotencyCache struct {
	client *redis.Client
}

func NewIdempotencyCache(client *redis.Client) infrastructure.IdempotencyCache {
	return &idempotencyCache{client: client}
}

// buildKey creates a uniquely scoped Redis key.
// Matrix scopes TxnIDs to the specific user, so we combine them.
func buildKey(accessToken, endpoint, txnID string) string {
	hash := sha256.Sum256([]byte(accessToken))
	tokenHash := hex.EncodeToString(hash[:])
	return fmt.Sprintf("idempotency:%s:%s:%s", tokenHash, endpoint, txnID)
}

func (c *idempotencyCache) Get(ctx context.Context, accessToken, endpoint, txnID string) (string, bool) {
	key := buildKey(accessToken, endpoint, txnID)

	// Fetch the value from Redis
	val, err := c.client.Get(ctx, key).Result()

	if err == redis.Nil {
		// redis.Nil is a special error meaning "Key does not exist".
		// This is a normal cache miss.
		return "", false
	} else if err != nil {
		// A real network error occurred (Redis is down).
		// In distributed systems, it's usually better to "fail open" and return false
		// so the user's message still sends, rather than locking up the chat app.
		// log.Printf("Redis error on Get: %v", err)
		return "", false
	}

	// Cache Hit! We found the Event ID.
	return val, true
}

func (c *idempotencyCache) Set(ctx context.Context, accessToken, endpoint, txnID, eventID string, ttl time.Duration) error {
	key := buildKey(accessToken, endpoint, txnID)

	err := c.client.Set(ctx, key, eventID, ttl).Err()
	if err != nil {
		return err
	}

	return nil
}
