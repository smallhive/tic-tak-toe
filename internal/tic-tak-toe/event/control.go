package event

const (
	TypeControlDisconnect = 1000
	TypeControlGameStared = 1001
)

func NewControlDisconnect() *Event {
	return &Event{
		Type: TypeControlDisconnect,
		Data: &NoBody{},
	}
}

type ControlGameStarted struct {
	ID int64
}

func NewControlGameStarted(id int64) *Event {
	return &Event{
		Type: TypeControlGameStared,
		Data: &ControlGameStarted{ID: id},
	}
}
