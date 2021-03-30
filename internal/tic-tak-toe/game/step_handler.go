package game

import (
	"context"

	"github.com/smallhive/tic-tak-toe/internal/tic-tak-toe/event"
)

func (s *Session) stepHandler(id int64, e *event.Step) error {
	p, _ := s.players[id]
	if !p.IsUserStep {
		p.Send(context.Background(), event.NewNotYouTurn())
		return nil
	}

	activePlayer, secondPlayer := s.detectPlayers(id)

	if s.field[e.Row][e.Coll] != MarkEmpty {
		return nil
	}

	s.stepCounter++

	s.field[e.Row][e.Coll] = activePlayer.Label

	activePlayer.IsUserStep = false
	activePlayer.Send(context.Background(), event.NewNotYouTurn())

	secondPlayer.IsUserStep = true
	secondPlayer.Send(context.Background(), event.NewYouTurn())

	s.broadcast(event.NewFieldUpdate(s.field))

	var ctx = context.Background()

	isWin, playerSign, winCond := s.checkWinCondition(s.field)
	if isWin {
		winner, loser := s.resolvePlayer(playerSign)
		if winner != nil && loser != nil {
			winner.Send(ctx, event.NewGameEnded(true, winCond))
			loser.Send(ctx, event.NewGameEnded(false, winCond))
		}

		winner.SendControl(ctx, event.NewControlDisconnect())
		loser.SendControl(ctx, event.NewControlDisconnect())

		// s.hub.DisconnectAll()
		s.terminate()
		return nil
	}

	if s.stepCounter >= 9 {
		s.broadcast(event.NewGameFailed())
		// s.hub.DisconnectAll()
		activePlayer.SendControl(ctx, event.NewControlDisconnect())
		secondPlayer.SendControl(ctx, event.NewControlDisconnect())

		s.terminate()
	}

	return nil
}
