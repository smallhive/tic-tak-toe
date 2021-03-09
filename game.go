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
	id      int64
	counter int
	field   [3][3]string
	players map[*Client]*Player
}

func NewWorld() *World {
	w := &World{
		id:      time.Now().UnixNano(),
		counter: 1,
		field:   [3][3]string{{"_", "_", "_"}, {"_", "_", "_"}, {"_", "_", "_"}},
		players: make(map[*Client]*Player),
	}

	return w
}

func (w *World) generateID() string {
	w.counter++
	return strconv.Itoa(w.counter)
}

func (w *World) IsFull() bool {
	return len(w.players) == 2
}

func (w *World) StartGame() {
	isFirst := true

	for _, p := range w.players {
		if isFirst {
			e := &Event{
				Type: EventTypeGameStarted,
				Data: &EventGameStared{IsFirstPlayer: true},
			}

			p.IsActive = true
			isFirst = false
			p.Client.conn.WriteJSON(e)

			e = &Event{
				Type: EventTypeYouTurn,
				Data: &EventYouTurn{},
			}

			p.Client.conn.WriteJSON(e)
		} else {
			e := &Event{
				Type: EventTypeGameStarted,
				Data: &EventGameStared{IsFirstPlayer: false},
			}
			p.Client.conn.WriteJSON(e)
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

	return "5"
}

func (w *World) Handle(c *Client, e *Event) {
	switch e.Type {
	case EventTypeStep:
		m, _ := json.Marshal(e.Data)
		var eventStep EventStep
		json.Unmarshal(m, &eventStep)

		p, _ := w.players[c]
		if !p.IsActive {
			e := &Event{
				Type: EventTypeNotYouTurn,
				Data: &EventNotYouTurn{},
			}

			p.Client.conn.WriteJSON(e)
			break
		}

		w.field[eventStep.Row][eventStep.Coll] = p.Label
		p.IsActive = false

		for cl, p := range w.players {
			if cl != c {
				p.IsActive = true
				e := &Event{
					Type: EventTypeYouTurn,
					Data: &EventYouTurn{},
				}
				p.Client.conn.WriteJSON(e)

				break
			}
		}

		e := &Event{
			Type: EventTypeFieldUpdate,
			Data: &EventFieldUpdate{Field: w.field},
		}

		me, _ := json.Marshal(e)
		c.hub.broadcast <- me
	}
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
