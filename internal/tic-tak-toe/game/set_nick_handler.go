package game

import (
	"context"

	"github.com/smallhive/tic-tak-toe/internal/tic-tak-toe/event"
)

func (s *Session) setNickHandler(id string, setNick event.Nick) error {
	activePlayer, secondPlayer := s.detectPlayers(id)
	activePlayer.Nick = setNick.Nick

	var ctx = context.Background()
	secondPlayer.Send(ctx, event.NewTypeSetOpponentNick(setNick.Nick))

	return nil
}
