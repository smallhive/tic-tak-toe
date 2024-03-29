package game

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/smallhive/tic-tak-toe/internal/logger"

	"github.com/smallhive/tic-tak-toe/internal/tic-tak-toe/game/player"
	"github.com/smallhive/tic-tak-toe/internal/tic-tak-toe/network"
)

type Queue struct {
	redis   *redis.Client
	gm      *Manager
	setName string
}

func NewQueue(r *redis.Client, gm *Manager) *Queue {
	return &Queue{
		redis:   r,
		setName: "userQueue",
		gm:      gm,
	}
}

func (q *Queue) Reset(ctx context.Context) error {
	r := q.redis.Del(ctx, q.setName)
	return r.Err()
}

func (q *Queue) Add(ctx context.Context, id string) error {
	r := q.redis.SAdd(ctx, q.setName, id)
	return r.Err()
}

func (q *Queue) two(ctx context.Context) ([]string, error) {
	r := q.redis.SRandMemberN(ctx, q.setName, 2)
	return r.Result()
}

func (q *Queue) MemberAmount(ctx context.Context) (int64, error) {
	r := q.redis.SCard(ctx, q.setName)
	return r.Result()
}

func (q *Queue) StartGame(ctx context.Context) error {
	ids, err := q.two(ctx)
	if err != nil {
		return err
	}

	var session = q.gm.Session()
	for _, id := range ids {
		q.redis.SRem(ctx, q.setName, id)

		var p = player.NewPlayer(
			id,
			session.UserMark(),
			network.NewProxy(q.redis, network.PlayerProxyChanName(id)),
			network.NewProxy(q.redis, network.PlayerProxyCommandChanName(id)),
		)

		session.AddPlayer(p)
	}

	var gameProxyChanName = network.GameProxyChanName(session.id)
	var pubSub = q.redis.Subscribe(ctx, gameProxyChanName)

	logger.Warn(ctx, "Sub", gameProxyChanName)
	session.Start(pubSub.Channel())

	return nil
}
