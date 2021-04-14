package player

import (
	"context"

	"github.com/smallhive/tic-tak-toe/internal/tic-tak-toe/event"
	"github.com/smallhive/tic-tak-toe/internal/tic-tak-toe/network"
)

// Player is a main entity for game (session) logic
type Player struct {
	ID           string `json:"id"`
	Label        string `json:"label"`
	IsUserStep   bool   `json:"-"`
	proxy        network.Proxy
	proxyControl network.Proxy
}

func NewPlayer(id string, label string, proxy network.Proxy, proxyControl network.Proxy) *Player {
	return &Player{
		ID:           id,
		Label:        label,
		proxy:        proxy,
		proxyControl: proxyControl,
	}
}

// Send sends game logic message to player
func (p *Player) Send(ctx context.Context, e *event.Event) error {
	e.UserID = p.ID
	return p.proxy.Send(ctx, e)
}

// SendControl sends control/system message to player
func (p *Player) SendControl(ctx context.Context, e *event.Event) error {
	e.UserID = p.ID
	return p.proxyControl.Send(ctx, e)
}
