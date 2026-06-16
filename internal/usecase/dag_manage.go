package usecase

import (
	"crypto/ed25519"
	"fmt"

	"github.com/caio-bernardo/dragonite/internal/domain"
	"github.com/caio-bernardo/dragonite/internal/util"
)

// genesisState simulates the room's database state during creation.
type genesisState struct {
	latestDepth   int64
	latestEventID string            // Tracks the tip of the DAG (PrevEvents)
	authCache     map[string]string // Key: "type|state_key", Value: EventID
}

// getAuth safely pulls an event ID from the cache if it exists.
func (g *genesisState) getAuth(eventType, stateKey string) (string, bool) {
	id, exists := g.authCache[fmt.Sprintf("%s|%s", eventType, stateKey)]
	return id, exists
}

// resolveGenesisAuth dynamically calculates required auth events from the in-memory cache.
func resolveGenesisAuth(ev *domain.Evento, state *genesisState) []string {
	auths := make([]string, 0)

	appendAuth := func(eType, sKey string) {
		if id, found := state.getAuth(eType, sKey); found {
			auths = append(auths, id)
		}
	}

	// 1. The Create event is the ONLY event that requires no auth.
	if ev.Tipo == "m.room.create" {
		return auths
	}

	// 2. ALL other events require the create event.
	appendAuth("m.room.create", "")

	// 3. The Creator's initial join only needs the create event.
	if ev.Tipo == "m.room.member" && *ev.StateKey == ev.Sender {
		// If power levels don't exist yet, we just need the create event.
		if _, found := state.getAuth("m.room.power_levels", ""); found {
			appendAuth("m.room.power_levels", "")
		}
		return auths
	}

	// 4. Everything else needs Power Levels and the Sender's Membership.
	appendAuth("m.room.power_levels", "")
	appendAuth("m.room.member", ev.Sender)

	// 5. If it's an invite/join for SOMEONE ELSE, we need join_rules.
	if ev.Tipo == "m.room.member" {
		appendAuth("m.room.join_rules", "")
	}

	return auths
}

// linkAndHashGenesis sequence chains a dynamic list of events together.
func linkAndHashGenesis(events []*domain.Evento, serverName, keyID string, privateKey ed25519.PrivateKey) error {
	state := &genesisState{
		authCache: make(map[string]string),
	}

	for i, ev := range events {
		// ---------------------------------------------------
		// STEP 1: Link the Timeline (PrevEvents)
		// ---------------------------------------------------
		if i == 0 {
			// The absolute beginning of time
			ev.PrevEventos = []string{}
			ev.Depth = 0
		} else {
			// Point strictly to the event generated immediately before this one,
			// creating a perfect, single-file line (linear DAG).
			ev.PrevEventos = []string{state.latestEventID}
			ev.Depth = state.latestDepth + 1
		}

		// ---------------------------------------------------
		// STEP 2: Link the Permissions (AuthEvents)
		// ---------------------------------------------------
		ev.AuthEventos = resolveGenesisAuth(ev, state)

		// ---------------------------------------------------
		// STEP 3: Cryptographic Hash
		// ---------------------------------------------------
		eventID, err := util.HashMatrixEvent(ev) // The Canonical JSON func we wrote
		if err != nil {
			return fmt.Errorf("failed to hash event %s: %w", ev.Tipo, err)
		}
		ev.ID = eventID

		sigJSON, err := util.SignMatrixEvent(ev, serverName, keyID, privateKey)
		if err != nil {
			return fmt.Errorf("failed to sign genesis event %s: %w", ev.Tipo, err)
		}
		ev.Signatures = sigJSON

		// ---------------------------------------------------
		// STEP 4: Update the Tracker for the Next Loop!
		// ---------------------------------------------------
		state.latestEventID = eventID
		state.latestDepth = ev.Depth

		// If this is a state event, cache it so future events in this loop can use it
		if ev.StateKey != nil {
			cacheKey := fmt.Sprintf("%s|%s", ev.Tipo, *ev.StateKey)
			state.authCache[cacheKey] = eventID
		}
	}

	return nil
}
