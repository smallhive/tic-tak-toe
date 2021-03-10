package game

import (
	"encoding/json"
	"time"

	"github.com/smallhive/tic-tak-toe/internal/tic-tak-toe/event"
	"github.com/smallhive/tic-tak-toe/internal/tic-tak-toe/network"
)

var (
	winConditions = [][][2]int{
		{[2]int{0, 0}, [2]int{1, 0}, [2]int{2, 0}},
		{[2]int{0, 1}, [2]int{1, 1}, [2]int{2, 1}},
		{[2]int{0, 2}, [2]int{1, 2}, [2]int{2, 2}},

		{[2]int{0, 0}, [2]int{0, 1}, [2]int{0, 2}},
		{[2]int{1, 0}, [2]int{1, 1}, [2]int{1, 2}},
		{[2]int{2, 0}, [2]int{2, 1}, [2]int{2, 2}},

		{[2]int{0, 0}, [2]int{1, 1}, [2]int{2, 2}},
		{[2]int{0, 2}, [2]int{1, 1}, [2]int{2, 0}},
	}
)

const (
	MarkCross = "X"
	MarkBigO  = "O"
)

type Session struct {
	hub          *network.Hub
	completeChan SessionCompleteChan

	id          int64
	field       [3][3]string
	players     map[*network.Client]*Player
	stepCounter int
}

func NewSession(hub *network.Hub, completeChan SessionCompleteChan) *Session {
	return &Session{
		hub:          hub,
		completeChan: completeChan,
		id:           time.Now().UnixNano(),
		field:        [3][3]string{{"_", "_", "_"}, {"_", "_", "_"}, {"_", "_", "_"}},
		players:      make(map[*network.Client]*Player),
		stepCounter:  0,
	}
}

func (s *Session) Hub() *network.Hub {
	return s.hub
}

func (s *Session) ID() int64 {
	return s.id
}

func (s *Session) IsFull() bool {
	return len(s.players) == 2
}

func (s *Session) userMark() string {
	switch len(s.players) {
	case 0:
		return MarkBigO
	case 1:
		return MarkCross
	default:
		return "NoneMark"
	}
}

func (s *Session) AddPlayer(c *network.Client) *Player {
	p := &Player{
		Client: c,
		ID:     time.Now().UnixNano(),
		Label:  s.userMark(),
	}

	s.players[c] = p
	return p
}

func (s *Session) Start() {
	isFirst := true
	for _, p := range s.players {
		p.Send(event.NewGameStared(isFirst))

		if isFirst {
			p.IsUserStep = true
			p.Send(event.NewYouTurn())
		} else {
			p.Send(event.NewNotYouTurn())
		}

		isFirst = false
	}
}

func (s *Session) Handle(c *network.Client, e *event.Event) {
	switch e.Type {
	case event.TypeStep:
		m, _ := json.Marshal(e.Data)
		var eventStep event.Step
		json.Unmarshal(m, &eventStep)

		s.stepHandler(c, &eventStep)
	}
}

func (s *Session) detectPlayers(c *network.Client) (*Player, *Player) {
	var p1, p2 *Player

	p1 = s.players[c]

	for _, p := range s.players {
		if p != p1 {
			p2 = p
		}
	}

	return p1, p2
}

func (s *Session) broadcast(e *event.Event) error {
	b, err := json.Marshal(e)
	if err != nil {
		return err
	}

	s.hub.Broadcast(b)
	return nil
}

func (s *Session) checkWinCondition(field [3][3]string) (bool, string, [][2]int) {
	var v string

	for _, cond := range winConditions {
		p1 := 0
		p2 := 0

		for i := 0; i < len(cond); i++ {
			row := cond[i][0]
			cel := cond[i][1]

			v = field[row][cel]
			if v == MarkBigO {
				p1++
			}

			if v == MarkCross {
				p2++
			}

			if p1 == 3 {
				return true, MarkBigO, cond
			}

			if p2 == 3 {
				return true, MarkCross, cond
			}
		}
	}

	return false, "", nil
}

func (s *Session) resolvePlayer(sign string) (*Player, *Player) {
	var winner, loser *Player

	for _, p := range s.players {
		if p.Label == sign {
			winner = p
		} else {
			loser = p
		}
	}

	return winner, loser
}

func (s *Session) terminate() {
	s.completeChan <- s
}
