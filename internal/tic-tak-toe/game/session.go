package game

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/smallhive/tic-tak-toe/internal/tic-tak-toe/event"
	"github.com/smallhive/tic-tak-toe/internal/tic-tak-toe/game/player"
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
	MarkEmpty = "_"
)

type Session struct {
	completeChan SessionCompleteChan

	id          string
	field       [3][3]string
	players     map[string]*player.Player
	stepCounter int
	cmdChan     <-chan *redis.Message
}

func NewSession(completeChan SessionCompleteChan) *Session {
	return &Session{
		completeChan: completeChan,
		id:           strconv.FormatInt(time.Now().UnixNano(), 16),
		field:        [3][3]string{{MarkEmpty, MarkEmpty, MarkEmpty}, {MarkEmpty, MarkEmpty, MarkEmpty}, {MarkEmpty, MarkEmpty, MarkEmpty}},
		players:      make(map[string]*player.Player),
		stepCounter:  0,
	}
}

func (s *Session) ID() string {
	return s.id
}

func (s *Session) IsFull() bool {
	return len(s.players) == 2
}

func (s *Session) UserMark() string {
	switch len(s.players) {
	case 0:
		return MarkBigO
	case 1:
		return MarkCross
	default:
		return "NoneMark"
	}
}

func (s *Session) AddPlayer(p *player.Player) *player.Player {
	s.players[p.ID] = p

	return p
}

func (s *Session) Start(cmdChan <-chan *redis.Message) {
	isFirst := true
	fmt.Println("GameStarting", s.id)
	for _, p := range s.players {
		p.Send(context.Background(), event.NewGameStared(isFirst, s.id))
		p.SendControl(context.Background(), event.NewControlGameStarted(s.id))

		if isFirst {
			p.IsUserStep = true
			p.Send(context.Background(), event.NewYouTurn())
		} else {
			p.Send(context.Background(), event.NewNotYouTurn())
		}

		isFirst = false
	}

	s.cmdChan = cmdChan

	go func() {
		for {
			message, ok := <-s.cmdChan
			if !ok {
				break
			}

			var e event.Event
			if err := json.Unmarshal([]byte(message.Payload), &e); err != nil {
				fmt.Println(err)
			} else {
				if err = s.Handle(&e); err != nil {
					fmt.Println(err)
				}
			}
		}
	}()
}

func (s *Session) Handle(e *event.Event) error {
	switch e.Type {
	case event.TypeStep:
		m, _ := json.Marshal(e.Data)
		var eventStep event.Step
		if err := json.Unmarshal(m, &eventStep); err != nil {
			return err
		}

		return s.stepHandler(e.UserID, &eventStep)
	}

	return nil
}

func (s *Session) detectPlayers(id string) (*player.Player, *player.Player) {
	var p1, p2 *player.Player

	p1 = s.players[id]

	for _, p := range s.players {
		if p != p1 {
			p2 = p
		}
	}

	return p1, p2
}

func (s *Session) broadcast(e *event.Event) error {
	for _, p := range s.players {
		p.Send(context.Background(), e)
	}

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

func (s *Session) resolvePlayer(sign string) (*player.Player, *player.Player) {
	var winner, loser *player.Player

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
