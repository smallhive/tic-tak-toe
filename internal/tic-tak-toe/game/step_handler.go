package game

import (
	"github.com/smallhive/tic-tak-toe/internal/tic-tak-toe/event"
	"github.com/smallhive/tic-tak-toe/internal/tic-tak-toe/network"
)

func (s *Session) stepHandler(c *network.Client, e *event.Step) {
	p, _ := s.players[c]
	if !p.IsUserStep {
		p.Send(event.NewNotYouTurn())
		return
	}

	s.stepCounter++

	activePlayer, secondPlayer := s.detectPlayers(c)

	s.field[e.Row][e.Coll] = activePlayer.Label

	activePlayer.IsUserStep = false
	activePlayer.Send(event.NewNotYouTurn())

	secondPlayer.IsUserStep = true
	secondPlayer.Send(event.NewYouTurn())

	s.broadcast(event.NewFieldUpdate(s.field))

	isWin, playerSign, winCond := s.checkWinCondition(s.field)
	if isWin {
		winner, loser := s.resolvePlayer(playerSign)
		if winner != nil && loser != nil {
			winner.Send(event.NewGameEnded(true, winCond))
			loser.Send(event.NewGameEnded(false, winCond))
		}

		s.hub.DisconnectAll()
		s.terminate()
		return
	}

	if s.stepCounter >= 9 {
		s.broadcast(event.NewGameFailed())
		s.hub.DisconnectAll()
		s.terminate()
	}
}
