package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"

	"github.com/smallhive/tic-tak-toe/app"
	"github.com/smallhive/tic-tak-toe/cmd/tic-tak-toe/web"
	"github.com/smallhive/tic-tak-toe/internal/tic-tak-toe/closer"
	"github.com/smallhive/tic-tak-toe/internal/tic-tak-toe/config"
	"github.com/smallhive/tic-tak-toe/internal/tic-tak-toe/game"
	"github.com/smallhive/tic-tak-toe/internal/tic-tak-toe/network"
)

func main() {
	fmt.Println(app.Name, app.Version, app.Commit)

	var c = config.Load()
	gm := game.NewManager()

	rdb := redis.NewClient(&redis.Options{
		Addr:         c.RedisAddr,
		MinIdleConns: 10,
		DB:           0,
	})

	hub := network.NewHub()
	go hub.Run()
	q := game.NewQueue(rdb, gm)
	q.Reset(context.Background())

	http.Handle("/", http.StripPrefix("/", http.FileServer(http.FS(web.Content))))
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, q, rdb, w, r)
	})

	err := http.ListenAndServe(c.Addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// serveWs handles websocket requests from the peer.
func serveWs(h *network.Hub, q *game.Queue, redisClient *redis.Client, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	ctx := context.Background()
	var id = strconv.FormatInt(time.Now().UnixNano(), 16)

	log.Println("user connected", id)

	var playerPubSub = redisClient.Subscribe(ctx, network.PlayerProxyChanName(id))
	var controlPubSub = redisClient.Subscribe(ctx, network.PlayerProxyCommandChanName(id))

	var cl = closer.NewCloser()
	client := network.NewClient(id, h, conn, redisClient, playerPubSub, controlPubSub, cl)

	if err := q.Add(ctx, id); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	amount, err := q.MemberAmount(ctx)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	cl.Add(func() error {
		return playerPubSub.Close()
	})
	cl.Add(func() error {
		return controlPubSub.Close()
	})

	go client.WritePump()
	go client.ReadPump()

	if amount > 1 {
		q.StartGame(ctx)
	}
}
