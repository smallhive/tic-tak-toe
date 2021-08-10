package game

import (
	"context"
)

// referee helps to check all players are ready
type referee struct {
	size    int
	current int

	com chan struct{}
}

func newReferee(size int, com chan struct{}) *referee {
	return &referee{size: size, com: com}
}

func (r *referee) Wait(ctx context.Context) error {
	for {
		select {
		case <-r.com:
			r.current++
			if r.current >= r.size {
				return nil
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (r *referee) Chan() chan struct{} {
	return r.com
}
