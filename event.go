package main

const (
	EventTypeInit        = 1
	EventTypeConnect     = 2
	EventTypeGameStarted = 3
	EventTypeYouTurn     = 4
	EventTypeNotYouTurn  = 5
	EventTypeStep        = 6
	EventTypeFieldUpdate = 7
)

type Event struct {
	Type int
	Data interface{}
}

type EventInit struct {
	Label string
}

type EventGameConnect struct {
}

type EventGameStared struct {
	IsFirstPlayer bool
}

type EventYouTurn struct {
}

type EventNotYouTurn struct {
}

type EventStep struct {
	Row  int
	Coll int
}

type EventFieldUpdate struct {
	Field [3][3]string
}
