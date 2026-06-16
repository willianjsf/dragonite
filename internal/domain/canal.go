package domain

import "time"

type Canal struct {
	ID        string
	Versao    string
	Criador   string
	CreatedAt time.Time

	// pontas do grafo, se torna os prev_eventos de cada nova mensagem
	ForwardExtremeties []string
	EstadoAtual        []StateEntry
}

type PublicRoomEntry struct {
	RoomID           string  `json:"room_id"`
	Name             *string `json:"name,omitempty"`
	Topic            *string `json:"topic,omitempty"`
	AvatarURL        *string `json:"avatar_url,omitempty"`
	CanonicalAlias   *string `json:"canonical_alias,omitempty"`
	NumJoinedMembers int     `json:"num_joined_members"`
	WorldReadable    bool    `json:"world_readable"`
	GuestCanJoin     bool    `json:"guest_can_join"`
	JoinRule         *string `json:"join_rule,omitempty"`
}

type PublicRoomsChunck struct {
	Chunk                  []PublicRoomEntry `json:"chunk"`
	NextBatch              string            `json:"next_batch,omitempty"`
	PrevBatch              string            `json:"prev_batch,omitempty"`
	TotalRoomCountEstimate int               `json:"total_room_count_estimate,omitempty"`
}

type StateEntry struct {
	Type     string
	StateKey string
	IDEvento string
}

// função auxiliar que encontra estados dentro do canal
func (c *Canal) GetStateEventID(eventType, stateKey string) (string, bool) {
	for _, state := range c.EstadoAtual {
		if state.Type == eventType && state.StateKey == stateKey {
			return state.IDEvento, true
		}
	}
	return "", false
}
