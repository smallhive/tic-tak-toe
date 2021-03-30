package player

import (
	"context"

	"github.com/smallhive/tic-tak-toe/internal/tic-tak-toe/event"
)

type Proxy interface {
	Send(ctx context.Context, e *event.Event) error
	SendControl(ctx context.Context, e *event.Event) error
}
