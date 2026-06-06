package util

import (
	"crypto/ed25519"
	"crypto/rand"
	"fmt"
)

func GenerateServerKey(serverName string, version string) (string, ed25519.PublicKey, ed25519.PrivateKey, error) {
	pubKey, privKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return "", nil, nil, fmt.Errorf("failed to generate key: %w", err)
	}
	return fmt.Sprintf("ed25519:%s", version), pubKey, privKey, nil
}
