package network

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"

	"github.com/smallhive/tic-tak-toe/internal/tic-tak-toe/event"
)

// Handler sends event to Game
type Handler interface {
	Handle(ctx context.Context, e *event.Event) error
}

type GameHandler struct {
	rc     *redis.Client
	config *GameProxyConfig
}

func NewGameHandler(rc *redis.Client, config *GameProxyConfig) *GameHandler {
	return &GameHandler{
		rc:     rc,
		config: config,
	}
}

func (g *GameHandler) Handle(ctx context.Context, e *event.Event) error {
	b, err := json.Marshal(e)
	if err != nil {
		return err
	}

	fmt.Println(g.config.ChanName, string(b))

	res := g.rc.Publish(ctx, g.config.ChanName, b)
	return res.Err()
}
