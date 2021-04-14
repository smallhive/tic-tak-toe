package network

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"

	"github.com/smallhive/tic-tak-toe/internal/tic-tak-toe/event"
)

func PlayerProxyChanName(id string) string {
	return fmt.Sprintf("user:%s", id)
}

func PlayerProxyCommandChanName(id string) string {
	return fmt.Sprintf("user:control:%s", id)
}

func GameProxyChanName(id string) string {
	return fmt.Sprintf("game:%s", id)
}

type Proxy interface {
	Send(ctx context.Context, e *event.Event) error
}

type proxy struct {
	redisClient *redis.Client
	chanName    string
}

func NewProxy(rc *redis.Client, chanName string) Proxy {
	return &proxy{
		redisClient: rc,
		chanName:    chanName,
	}
}

func (p *proxy) Send(ctx context.Context, e *event.Event) error {
	b, err := json.Marshal(e)
	if err != nil {
		return err
	}

	res := p.redisClient.Publish(ctx, p.chanName, b)
	return res.Err()
}
