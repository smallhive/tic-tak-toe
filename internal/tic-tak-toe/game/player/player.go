package player

import (
	"context"

	"github.com/smallhive/tic-tak-toe/internal/tic-tak-toe/event"
)

type Player struct {
	ID         string `json:"id"`
	Label      string `json:"label"`
	IsUserStep bool   `json:"-"`
	proxy      Proxy
}

func NewPlayer(id string, label string, proxy Proxy) *Player {
	return &Player{
		ID:    id,
		Label: label,
		proxy: proxy,
	}
}

func (p *Player) Send(ctx context.Context, e *event.Event) error {
	e.UserID = p.ID
	return p.proxy.Send(ctx, e)
}

func (p *Player) SendControl(ctx context.Context, e *event.Event) error {
	e.UserID = p.ID
	return p.proxy.SendControl(ctx, e)
}
