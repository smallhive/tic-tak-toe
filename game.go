package main

import (
	"encoding/json"
	"strconv"
	"time"
)

const MarkCross = "X"
const MarkBigO = "O"

type Player struct {
	Client   *Client `json:"-"`
	ID       string
	Label    string
	IsActive bool
}

type World struct {
	id            int64
	playerCounter int
	field         [3][3]string
	players       map[*Client]*Player
	stepCounter   int
}

func NewWorld() *World {
	w := &World{
		id: time.Now().UnixNano(),
	}

	w.resetGame()

	return w
}

func (w *World) generateID() string {
	w.playerCounter++
	return strconv.Itoa(w.playerCounter)
}

func (w *World) IsFull() bool {
	return len(w.players) == 2
}

func (w *World) resetGame() {
	w.players = make(map[*Client]*Player)
	w.field = [3][3]string{{"_", "_", "_"}, {"_", "_", "_"}, {"_", "_", "_"}}
	w.playerCounter = 1
	w.stepCounter = 0
}

func (w *World) StartGame() {
	isFirst := true

	for _, p := range w.players {
		if isFirst {
			e := NewEventGameStared(true)
			p.IsActive = true
			isFirst = false
			// p.Client.conn.WriteJSON(e)
			p.Client.send <- e.JSON()

			e = NewEventYouTurn()

			// p.Client.conn.WriteJSON(e)
			p.Client.send <- e.JSON()
		} else {
			e := NewEventGameStared(false)
			// p.Client.conn.WriteJSON(e)
			p.Client.send <- e.JSON()

			e = NewEventNotYouTurn()
			// p.Client.conn.WriteJSON(e)
			p.Client.send <- e.JSON()
		}
	}
}

func (w *World) userLabel() string {
	if len(w.players) == 0 {
		return MarkBigO
	}

	if len(w.players) == 1 {
		return MarkCross
	}

	return "IQIIWIIWIW"
}

func (w *World) Handle(c *Client, e *Event) {
	switch e.Type {
	case EventTypeStep:
		m, _ := json.Marshal(e.Data)
		var eventStep EventStep
		json.Unmarshal(m, &eventStep)

		p, _ := w.players[c]
		if !p.IsActive {
			e := NewEventNotYouTurn()

			// p.Client.conn.WriteJSON(e)
			p.Client.send <- e.JSON()
			break
		}

		w.stepCounter++

		w.field[eventStep.Row][eventStep.Coll] = p.Label
		p.IsActive = false

		for cl, p := range w.players {
			var e *Event
			if cl != c {
				p.IsActive = true
				e = NewEventYouTurn()

			} else {
				e = NewEventNotYouTurn()
			}

			// cl.conn.WriteJSON(e)
			cl.send <- e.JSON()
		}

		e := NewEventFieldUpdate(w.field)
		c.hub.broadcast <- e.JSON()

		isWin, playerSign := checkWinCondition(w.field)
		if isWin {
			winner, loser := w.resolvePlayer(playerSign)
			if winner != nil && loser != nil {
				e := NewEventGameEnded(true)

				// winner.Client.conn.WriteJSON(e)
				winner.Client.send <- e.JSON()

				e = NewEventGameEnded(false)

				// loser.Client.conn.WriteJSON(e)
				loser.Client.send <- e.JSON()
			}

			c.hub.unregister <- winner.Client
			c.hub.unregister <- loser.Client

			w.resetGame()
			return
		}

		if w.stepCounter >= 9 {
			e := &Event{
				Type: EventTypeGameFailed,
				Data: &EventNoBody{},
			}

			me, _ := json.Marshal(e)
			c.hub.broadcast <- me
			c.hub.DisconnectAll()

			w.resetGame()
		}
	}
}

func (w *World) resolvePlayer(sign string) (*Player, *Player) {
	var winner, loser *Player

	for _, p := range w.players {
		if p.Label == sign {
			winner = p
		} else {
			loser = p
		}
	}

	return winner, loser
}

func (w *World) AddPlayer(c *Client) *Player {
	id := w.generateID()
	p := &Player{
		ID:     id,
		Label:  w.userLabel(),
		Client: c,
	}

	w.players[c] = p
	return p
}

func checkWinCondition(field [3][3]string) (bool, string) {
	winConditions := [][][2]int{
		{[2]int{0, 0}, [2]int{1, 0}, [2]int{2, 0}},
		{[2]int{0, 1}, [2]int{1, 1}, [2]int{2, 1}},
		{[2]int{0, 2}, [2]int{1, 2}, [2]int{2, 2}},

		{[2]int{0, 0}, [2]int{0, 1}, [2]int{0, 2}},
		{[2]int{1, 0}, [2]int{1, 1}, [2]int{2, 1}},
		{[2]int{2, 0}, [2]int{2, 1}, [2]int{2, 2}},

		{[2]int{0, 0}, [2]int{1, 1}, [2]int{2, 2}},
		{[2]int{0, 2}, [2]int{1, 1}, [2]int{2, 0}},
	}

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
				return true, MarkBigO
			}

			if p2 == 3 {
				return true, MarkCross
			}
		}
	}

	return false, ""
}
