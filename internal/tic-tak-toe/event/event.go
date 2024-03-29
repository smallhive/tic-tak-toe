package event

import (
	"encoding/json"
)

const (
	TypeInit        = 1
	TypeConnect     = 2
	TypeGameStarted = 3
	TypeYouTurn     = 4
	TypeNotYouTurn  = 5
	TypeStep        = 6
	TypeFieldUpdate = 7
	TypeGameEnded   = 8
	TypeGameFailed  = 9

	// TypeUnexpectedDisconnect will be sent to Game, when one player has disconnected
	TypeUnexpectedDisconnect = 10

	// TypeOpponentUnexpectedDisconnect will be send to Second player if first quit
	TypeOpponentUnexpectedDisconnect = 11

	TypeSetNick         = 12
	TypeSetOpponentNick = 13

	TypeAreYouReady = 14
	TypeIamReady    = 15
)

type NoBody struct {
}

type Event struct {
	UserID string      `json:"id,omitempty"`
	Type   int         `json:"type"`
	Data   interface{} `json:"data"`
}

func (e *Event) JSON() []byte {
	b, _ := json.Marshal(e)
	return b
}

type Init struct {
	Label  string
	GameID int64
}

func NewInit(label string, gameID int64) *Event {
	return &Event{
		Type: TypeInit,
		Data: &Init{
			Label:  label,
			GameID: gameID,
		},
	}
}

type GameConnect struct {
}

type GameStarted struct {
	IsFirstPlayer bool
	ID            string
}

func NewGameStared(IsFirstPlayer bool, id string) *Event {
	return &Event{
		Type: TypeGameStarted,
		Data: &GameStarted{IsFirstPlayer: IsFirstPlayer, ID: id},
	}
}

func NewYouTurn() *Event {
	return &Event{
		Type: TypeYouTurn,
		Data: &NoBody{},
	}
}

func NewNotYouTurn() *Event {
	return &Event{
		Type: TypeNotYouTurn,
		Data: &NoBody{},
	}
}

type Step struct {
	Row  int
	Coll int
}

type FieldUpdate struct {
	Field [3][3]string
}

func NewFieldUpdate(field [3][3]string) *Event {
	return &Event{
		Type: TypeFieldUpdate,
		Data: &FieldUpdate{Field: field},
	}
}

type GameEnded struct {
	IsWin     bool
	Condition [][2]int
}

func NewGameEnded(IsWin bool, Condition [][2]int) *Event {
	return &Event{
		Type: TypeGameEnded,
		Data: &GameEnded{IsWin: IsWin, Condition: Condition},
	}
}

func NewGameFailed() *Event {
	return &Event{
		Type: TypeGameFailed,
		Data: &NoBody{},
	}
}

func NewUnexpectedDisconnect(id string) *Event {
	return &Event{
		UserID: id,
		Type:   TypeUnexpectedDisconnect,
		Data:   &NoBody{},
	}
}

func NewOpponentUnexpectedDisconnect() *Event {
	return &Event{
		Type: TypeOpponentUnexpectedDisconnect,
		Data: &NoBody{},
	}
}

type Nick struct {
	Nick string
}

func NewSetNick(nick string) *Event {
	return &Event{
		Type: TypeSetNick,
		Data: &Nick{Nick: nick},
	}
}

func NewTypeSetOpponentNick(nick string) *Event {
	return &Event{
		Type: TypeSetOpponentNick,
		Data: &Nick{Nick: nick},
	}
}

type AreYouReady struct {
	ID string
}

func NewAreYouReady(id string) *Event {
	return &Event{
		Type: TypeAreYouReady,
		Data: &AreYouReady{ID: id},
	}
}

func NewIamReady() *Event {
	return &Event{
		Type: TypeIamReady,
		Data: &NoBody{},
	}
}
