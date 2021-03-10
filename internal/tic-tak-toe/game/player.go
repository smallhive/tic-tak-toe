package game

import (
	"encoding/json"

	"github.com/smallhive/tic-tak-toe/internal/tic-tak-toe/event"
	"github.com/smallhive/tic-tak-toe/internal/tic-tak-toe/network"
)

type Player struct {
	Client     *network.Client `json:"-"`
	ID         int64           `json:"id"`
	Label      string          `json:"label"`
	IsUserStep bool            `json:"-"`
}

func (p *Player) Send(e *event.Event) error {
	b, err := json.Marshal(e)
	if err != nil {
		return err
	}

	p.Client.Send(b)
	return nil
}
