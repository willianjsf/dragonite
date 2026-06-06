package domain

import (
	"encoding/json"
	"time"
)

type Evento struct {
	ID               string          `json:"id"`
	Tipo             string          `json:"tipo"`
	Content          json.RawMessage `json:"content"`
	CanalID          string          `json:"canal_id"`
	Sender           string          `json:"sender"`
	OrigemServidorTS time.Time       `json:"origem_servidor_ts"`

	StreamOrdering int64 `json:"stream_ordering"`

	StateKey string `json:"state_key"`
	Redacts  string `json:"redacts"`

	PrevEventos []string `json:"prev_eventos"`
	AuthEventos []string `json:"auth_eventos"`
	Depth       int64    `json:"depth"`

	Hashes     json.RawMessage `json:"hashes"`
	Signatures json.RawMessage `json:"signature"`
	Unsigned   json.RawMessage `json:"unsigned"` // dados adicionados pelo servidor
}
