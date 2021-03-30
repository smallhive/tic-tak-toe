package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"

	"github.com/smallhive/tic-tak-toe/internal/tic-tak-toe/game"
	"github.com/smallhive/tic-tak-toe/internal/tic-tak-toe/network"
)

// Отправка в закрытый канал. Проверить что заканчиваются горутины в хабах

var addr = flag.String("addr", ":8080", "http service address")

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "web/index.html")
}

func main() {
	flag.Parse()

	gm := game.NewManager()

	rdb := redis.NewClient(&redis.Options{
		Addr:         "lan_ip:6379",
		MinIdleConns: 10,
		DB:           0,
	})

	hub := network.NewHub()
	go hub.Run()
	q := game.NewQueue(rdb, gm)
	q.Reset(context.Background())

	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, q, rdb, w, r)
	})

	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
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

	var proxyConfig = network.NewPlayerProxyConfig(id)

	var playerPubSub = redisClient.Subscribe(ctx, proxyConfig.UserChanName)
	var controlPubSub = redisClient.Subscribe(ctx, proxyConfig.ControlChanName)

	client := network.NewClient(id, h, conn, redisClient, playerPubSub, controlPubSub)

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

	go client.WritePump()
	go client.ReadPump()

	if amount > 1 {
		q.StartGame(ctx)
	}
}
