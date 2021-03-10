package network

import (
	"github.com/smallhive/tic-tak-toe/internal/tic-tak-toe/event"
)

type Handler interface {
	Handle(c *Client, e *event.Event)
}
