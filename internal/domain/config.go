package domain

import "crypto/ed25519"

// Config for the server
type Config struct {
	ServerName string
	Version    string
	Port       int
	KeyID      string
	PublicKey  ed25519.PublicKey
	PrivateKey ed25519.PrivateKey
}
