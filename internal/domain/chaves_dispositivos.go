package domain

import "encoding/json"

// ChavesDispositivo representa as chaves de identidade E2EE (Olm/Megolm) de um dispositivo
type ChavesDispositivo struct {
	DispositivoID   string          `json:"-"`
	UsuarioID       string          `json:"-"`
	NomeDispositivo string          `json:"-"` // usado para preencher unsigned.device_display_name na resposta
	Algorithms      []string        `json:"algorithms"`
	Keys            json.RawMessage `json:"keys"`
	Signatures      json.RawMessage `json:"signatures"`
}

// ChaveUsoUnico representa uma one-time key de um dispositivo, consumida ao ser reivindicada
type ChaveUsoUnico struct {
	DispositivoID string
	KeyID         string // formato "<algorithm>:<id>", ex: "signed_curve25519:AAAAHg"
	Algorithm     string
	KeyData       json.RawMessage // KeyObject: {key, signatures}
}

// ChaveFallback representa a fallback key de um dispositivo para um algoritmo (reutilizável até ser substituída)
type ChaveFallback struct {
	DispositivoID string
	Algorithm     string
	KeyID         string
	KeyData       json.RawMessage
	Usada         bool
}