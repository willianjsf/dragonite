package domain

import "encoding/json"

// Filter define as regras de escopo que o cliente impõe sobre o que quer receber no /sync.
type Filter struct {
	ID         string          // Gerado pelo servidor (ex: "1", "2")
	UserID     string          // Dono do filtro
	Definition json.RawMessage // O JSON estruturado com regras de inclusão/exclusão de eventos
}
