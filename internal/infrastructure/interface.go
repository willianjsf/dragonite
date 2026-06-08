package infrastructure

import (
	"context"
	"time"
)

type IdempotencyCache interface {
	Get(ctx context.Context, accessToken, endpoint, txnId string) (eventID string, exists bool)
	Set(ctx context.Context, accessToken, endpoint, txnId, eventID string, ttl time.Duration) error
}
