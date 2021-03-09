package main

import (
	"encoding/json"
)

const (
	EventTypeInit        = 1
	EventTypeConnect     = 2
	EventTypeGameStarted = 3
	EventTypeYouTurn     = 4
	EventTypeNotYouTurn  = 5
	EventTypeStep        = 6
	EventTypeFieldUpdate = 7
	EventTypeGameEnded   = 8
	EventTypeGameFailed  = 9
)

type EventNoBody struct {
}

type Event struct {
	Type int
	Data interface{}
}

func (e *Event) JSON() []byte {
	b, _ := json.Marshal(e)
	return b
}

type EventInit struct {
	Label string
}

func NewEventInit(label string) *Event {
	return &Event{
		Type: EventTypeInit,
		Data: &EventInit{Label: label},
	}
}

type EventGameConnect struct {
}

type EventGameStared struct {
	IsFirstPlayer bool
}

func NewEventGameStared(IsFirstPlayer bool) *Event {
	return &Event{
		Type: EventTypeGameStarted,
		Data: &EventGameStared{IsFirstPlayer: IsFirstPlayer},
	}
}

func NewEventYouTurn() *Event {
	return &Event{
		Type: EventTypeYouTurn,
		Data: &EventNoBody{},
	}
}

func NewEventNotYouTurn() *Event {
	return &Event{
		Type: EventTypeNotYouTurn,
		Data: &EventNoBody{},
	}
}

type EventStep struct {
	Row  int
	Coll int
}

type EventFieldUpdate struct {
	Field [3][3]string
}

func NewEventFieldUpdate(field [3][3]string) *Event {
	return &Event{
		Type: EventTypeFieldUpdate,
		Data: &EventFieldUpdate{Field: field},
	}
}

type EventGameEnded struct {
	IsWin bool
}

func NewEventGameEnded(IsWin bool) *Event {
	return &Event{
		Type: EventTypeGameEnded,
		Data: &EventGameEnded{IsWin: IsWin},
	}
}
