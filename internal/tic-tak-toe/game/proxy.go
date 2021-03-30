package game

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"

	"github.com/smallhive/tic-tak-toe/internal/tic-tak-toe/event"
	"github.com/smallhive/tic-tak-toe/internal/tic-tak-toe/network"
)

type PlayerProxy struct {
	redisClient *redis.Client
	config      *network.PlayerProxyConfig
}

func NewPlayerProxy(rc *redis.Client, config *network.PlayerProxyConfig) *PlayerProxy {
	return &PlayerProxy{
		redisClient: rc,
		config:      config,
	}
}

func (p *PlayerProxy) Send(ctx context.Context, e *event.Event) error {
	b, err := json.Marshal(e)
	if err != nil {
		return err
	}

	fmt.Println(p.config.UserChanName, string(b))

	res := p.redisClient.Publish(ctx, p.config.UserChanName, b)
	return res.Err()
}

func (p *PlayerProxy) SendControl(ctx context.Context, e *event.Event) error {
	b, err := json.Marshal(e)
	if err != nil {
		return err
	}

	fmt.Println(p.config.ControlChanName, string(b))

	res := p.redisClient.Publish(ctx, p.config.ControlChanName, b)
	return res.Err()
}
