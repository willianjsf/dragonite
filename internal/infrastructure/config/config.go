package config

import (
	"crypto/ed25519"
	"errors"
	"os"
	"strconv"

	"github.com/caio-bernardo/dragonite/internal/util"
	_ "github.com/joho/godotenv/autoload"
)

type AppConfig struct {
	ServerName  string `env:"SERVER_NAME"`
	ServerPort  int    `env:"BACKEND_PORT"`
	Version     string `env:"VERSION"`
	JWTToken    string `env:"JWT_TOKEN"`
	KeyID       string `env:"KEY_ID"`
	DatabaseURL string `env:"DATABASE_URL"`
	PublicKey   ed25519.PublicKey
	PrivateKey  ed25519.PrivateKey
}

func LoadConfig() (*AppConfig, error) {
	if os.Getenv("BACKEND_PORT") == "" {
		return nil, errors.New("BACKEND_PORT variable is not set")
	}
	port, err := strconv.Atoi(os.Getenv("BACKEND_PORT"))
	if err != nil {
		return nil, errors.New("BACKEND_PORT is not a valid integer")
	}

	if os.Getenv("SERVER_NAME") == "" {
		return nil, errors.New("SERVER_NAME variable is not set")
	}

	if os.Getenv("DATABASE_URL") == "" {
		return nil, errors.New("DATABASE_URL variable is not set")
	}

	if os.Getenv("VERSION") == "" {
		return nil, errors.New("VERSION variable is not set")
	}

	if os.Getenv("JWT_TOKEN") == "" {
		return nil, errors.New("JWT_TOKEN variable is not set")
	}

	key, pubKey, privKey, err := util.GenerateServerKey(os.Getenv("SERVER_NAME"), os.Getenv("VERSION"))
	if err != nil {
		return nil, err
	}

	config := AppConfig{
		ServerName:  os.Getenv("SERVER_NAME"),
		ServerPort:  port,
		Version:     os.Getenv("VERSION"),
		JWTToken:    os.Getenv("JWT_TOKEN"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
		KeyID:       key,
		PublicKey:   pubKey,
		PrivateKey:  privKey,
	}

	return &config, nil
}
