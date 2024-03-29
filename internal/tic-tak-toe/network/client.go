package network

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"

	"github.com/smallhive/tic-tak-toe/internal/logger"
	"github.com/smallhive/tic-tak-toe/internal/tic-tak-toe/closer"
	"github.com/smallhive/tic-tak-toe/internal/tic-tak-toe/event"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 1000 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 600 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	id string

	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	// Channel for receiving messages from Game
	proxy *redis.PubSub
	// Channel for receiving Control/System messages from Game
	control *redis.PubSub

	redis   *redis.Client
	handler Proxy

	closer *closer.Closer
}

func NewClient(id string, hub *Hub, conn *websocket.Conn, redis *redis.Client, proxy *redis.PubSub, control *redis.PubSub, waiter *closer.Closer) *Client {
	c := &Client{
		id:      id,
		hub:     hub,
		conn:    conn,
		send:    make(chan []byte, maxMessageSize),
		proxy:   proxy,
		control: control,
		redis:   redis,
		closer:  waiter,
	}

	hub.register <- c

	return c
}

func (c *Client) Send(data []byte) {
	c.send <- data
}

// ReadPump pumps messages from the websocket connection to the hub.
//
// The application runs ReadPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) ReadPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	ctx := context.Background()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logger.Errorf(ctx, "error: %v", err)
			}
			break
		}

		if c.handler == nil {
			logger.Errorf(ctx, "logic handler isn't set for client", c.id)
			continue
		}

		var e event.Event
		if err := json.Unmarshal(message, &e); err != nil {
			logger.Error(ctx, err)
		} else {
			e.UserID = c.id
			if err := c.handler.Send(ctx, &e); err != nil {
				logger.Error(ctx, err)
			}
		}
	}
}

// WritePump pumps messages from the hub to the websocket connection.
//
// A goroutine running WritePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	ctx := context.Background()

	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}

		case message, ok := <-c.proxy.Channel():
			if !ok {
				c.hub.unregister <- c
			} else {
				c.send <- []byte(message.Payload)
			}

		case message, ok := <-c.control.Channel():
			if !ok {
				logger.Error(ctx, "control chan closed")
				continue
			}
			var e event.Event
			if err := json.Unmarshal([]byte(message.Payload), &e); err != nil {
				logger.Error(ctx, err)
			} else {
				if err := c.handleControl(&e); err != nil {
					logger.Error(ctx, err)
				}
			}
		}
	}
}

func (c *Client) handleControl(e *event.Event) error {
	switch e.Type {
	case event.TypeControlDisconnect:
		c.hub.unregister <- c

	case event.TypeControlLinkGameHandler:
		m, _ := json.Marshal(e.Data)
		var startedEvent event.ControlLinkGameHandler
		if err := json.Unmarshal(m, &startedEvent); err != nil {
			return err
		}

		var h = NewProxy(c.redis, GameProxyChanName(startedEvent.ID))
		c.handler = h

		c.closer.Add(func() error {
			return h.Send(context.Background(), event.NewUnexpectedDisconnect(c.id))
		})
	}

	return nil
}

func (c *Client) Close(ctx context.Context) {
	c.closer.Close(ctx)
}
