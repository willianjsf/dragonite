package usecase

type SyncService struct{}

// userID := ctx.Value(types.UserIDKey).(string)
// 	// Lógica de Long-Polling
// 	if since.TimelinePosition != 0 {
// 		hasEvents, err := u.eventoStore.CheckNew(ctx, userID, since)
// 		if err != nil {
// 			return nil, nil, err
// 		}

// 		if !hasEvents && timeout > 0 {
// 			// sem eventos, long-polling
// 			ch := u.notifier.Subscribe(userID)
// 			defer u.notifier.Unsubscribe(userID, ch)

// 			select {
// 			case <-ch:
// 				// Novo evento, pode acessar o banco
// 			case <-time.After(timeout):
// 				// Deu timeout antes de um novo evento, cria novo token e retorna
// 				maxGlobal, _ := u.eventoStore.GetMaxGlobalStreamOrdering(ctx)
// 				if maxGlobal > since.TimelinePosition {
// 					since.TimelinePosition = maxGlobal
// 				}
// 				return nil, &since, types.ErrTimeout
// 			case <-ctx.Done():
// 				// o client se desconectou
// 				return nil, nil, types.ErrLooseConnection
// 			}
// 		}
// 	}

// 	// accesso ao banco
// 	events, newToken, err := u.eventoStore.GetSince(ctx, userID, since)
// 	if err != nil {
// 		return nil, nil, err
// 	}
// 	return events, newToken, nil
