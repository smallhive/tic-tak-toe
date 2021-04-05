package game

import (
	"context"

	"github.com/smallhive/tic-tak-toe/internal/tic-tak-toe/event"
)

func (s *Session) unexpectedDisconnectHandler(id string) error {
	_, secondPlayer := s.detectPlayers(id)

	var ctx = context.Background()

	secondPlayer.Send(ctx, event.NewOpponentUnexpectedDisconnect())
	secondPlayer.SendControl(ctx, event.NewControlDisconnect())

	s.terminate()

	return nil
}
