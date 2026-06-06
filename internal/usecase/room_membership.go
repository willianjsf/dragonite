package usecase

import (
	"context"
	"fmt"

	"github.com/caio-bernardo/dragonite/internal/util"
)

type RoomMembershipService struct {
	eventBus  EventBus
	uow       WorkUnit
	canalRepo CanalStorage
	eventRepo EventoStorage
}

func NewRoomMembershipService(eventBus EventBus, canalRepo CanalStorage, eventRepo EventoStorage) *RoomMembershipService {
	return &RoomMembershipService{
		eventBus:  eventBus,
		canalRepo: canalRepo,
		eventRepo: eventRepo,
	}
}

func (s *RoomMembershipService) LeaveRoom(ctx context.Context, userID, roomID string) error {
	// 1. Verify they are actually in the room (or invited)
	currentStatus, err := s.canalRepo.GetUserMembership(ctx, roomID, userID)
	if err != nil || (currentStatus != "join" && currentStatus != "invite") {
		return fmt.Errorf("user is not in a state to leave this room")
	}

	// 2. Build the Leave Event
	leaveEvent := buildLeaveEvent(roomID, userID)

	// 3. Resolve DAG dependencies (Use the method we built earlier)
	prevs, auths, err := s.resolveEventDependencies(ctx, roomID, userID, "m.room.member")
	if err != nil {
		return err
	}
	leaveEvent.PrevEventos = prevs
	leaveEvent.AuthEventos = auths

	// 4. Hash it
	eventID, _ := util.HashMatrixEvent(leaveEvent)
	leaveEvent.ID = eventID

	// 5. Transactional Save
	err = s.uow.Execute(ctx, func(txCtx context.Context) error {
		// A. Save the event to the DAG
		if err := s.eventRepo.SaveEvent(txCtx, leaveEvent); err != nil {
			return err
		}

		// B. Update the denormalized fast-lookup table
		if err := s.canalRepo.UpsertMembership(txCtx, roomID, userID, "leave"); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	// 6. Notify local clients (so the room disappears from their UI instantly)
	s.eventBus.Publish(ctx, roomID, *leaveEvent)
	return nil
}

func (s *RoomMembershipService) JoinLocalRoom(ctx context.Context, userID, roomID string) error {

	// 1. Fetch the Join Rules
	joinRule, err := s.canalRepo.GetJoinRule(ctx, roomID)
	if err != nil {
		return fmt.Errorf("room not found or missing join rules")
	}

	// 2. Validate Access
	if joinRule == "invite" {
		// They MUST have a pending invite in our database
		status, _ := s.canalRepo.GetUserMembership(ctx, roomID, userID)
		if status != "invite" {
			return fmt.Errorf("M_FORBIDDEN: You must be invited to join this room")
		}
	} else if joinRule != "public" {
		return fmt.Errorf("M_FORBIDDEN: Room is not public")
	}

	// 3. Build the Join Event

	joinEvent := buildJoinEvent(roomID, userID)

	// 4. Resolve DAG Dependencies
	prevs, auths, err := s.resolveEventDependencies(ctx, roomID, userID, "m.room.member")
	if err != nil {
		return err
	}
	joinEvent.PrevEventos = prevs
	joinEvent.AuthEventos = auths

	// 5. Hash it
	eventID, _ := util.HashMatrixEvent(joinEvent)
	joinEvent.ID = eventID

	// 6. Transactional Save
	err = s.uow.Execute(ctx, func(txCtx context.Context) error {
		if err := s.eventRepo.SaveEvent(txCtx, joinEvent); err != nil {
			return err
		}

		// Update their status from "invite" (or null) to "join"
		if err := s.canalRepo.UpsertMembership(txCtx, roomID, userID, "join"); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}

	// 7. Notify clients
	s.eventBus.Publish(ctx, roomID, *joinEvent)
	return nil
}

func (s *RoomMembershipService) resolveEventDependencies(
	ctx context.Context,
	roomID, sender string,
	eventType string,
) ([]string, []string, error) {

	// 1. Resolve PrevEvents (The tips of the DAG)
	// This is simply querying the room_forward_extremities table.
	prevEvents, err := s.canalRepo.GetForwardExtremities(ctx, roomID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get extremities: %w", err)
	}

	// 2. Resolve AuthEvents (The VIP Pass)
	// The exact state events needed change based on what kind of event is being sent.
	// But almost ALL events require at least the create, power_levels, and sender's membership.

	authEvents := make([]string, 0)

	// Helper to fetch and append state safely
	appendState := func(stateType, stateKey string) {
		if eventID, found := s.canalRepo.GetStateEventID(ctx, roomID, stateType, stateKey); found {
			authEvents = append(authEvents, eventID)
		}
	}

	// ALWAYS REQUIRED: The room creation event
	appendState("m.room.create", "")

	// ALWAYS REQUIRED: The room power levels
	appendState("m.room.power_levels", "")

	// ALWAYS REQUIRED: The sender's membership (proving they are in the room)
	appendState("m.room.member", sender)

	// CONDITIONAL: If they are inviting someone, we also need the invitee's current membership state
	if eventType == "m.room.member" {
		// (Assuming you passed the target stateKey into this function)
		// appendState("m.room.member", targetUserID)

		// And we need the join_rules to see if invites are allowed
		appendState("m.room.join_rules", "")
	}

	return prevEvents, authEvents, nil
}
