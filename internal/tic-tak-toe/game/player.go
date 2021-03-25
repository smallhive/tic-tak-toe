package game

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/smallhive/tic-tak-toe/internal/tic-tak-toe/event"
)

type Player struct {
	ID          int64         `json:"id"`
	Label       string        `json:"label"`
	IsUserStep  bool          `json:"-"`
	redisClient *redis.Client `json:"-"`
}

func (p *Player) Send(e *event.Event) error {
	b, err := json.Marshal(e)
	if err != nil {
		return err
	}

	var chanName = fmt.Sprintf("user:%d", p.ID)
	res := p.redisClient.Publish(context.Background(), chanName, b)
	return res.Err()
}
