package types

import (
	"github.com/golang-jwt/jwt/v5"
)

type MatrixClaims struct {
	UserID   string `json:"user_id"`
	DeviceID string `json:"device_id"`
	jwt.RegisteredClaims
}

type contextKey string

const UserIDKey contextKey = "user_id"
const DeviceIDKey contextKey = "device_id"
