package event

const (
	TypeControlDisconnect      = 1000
	TypeControlLinkGameHandler = 1001
)

func NewControlDisconnect() *Event {
	return &Event{
		Type: TypeControlDisconnect,
		Data: &NoBody{},
	}
}

type ControlLinkGameHandler struct {
	ID string
}

func NewControlLinkGameHandler(id string) *Event {
	return &Event{
		Type: TypeControlLinkGameHandler,
		Data: &ControlLinkGameHandler{ID: id},
	}
}
