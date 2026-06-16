package domain

import "time"

// ReadReceipt rastreia a última mensagem que um usuário leu explicitamente em uma sala.
type ReadReceipt struct {
	RoomID    string    `json:"room_id"`
	UserID    string    `json:"user_id"`
	EventID   string    `json:"event_id"` // O evento que o usuário leu
	Timestamp time.Time `json:"ts"`
}
