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
