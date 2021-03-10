package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/smallhive/tic-tak-toe/internal/tic-tak-toe/event"
	"github.com/smallhive/tic-tak-toe/internal/tic-tak-toe/game"
	"github.com/smallhive/tic-tak-toe/internal/tic-tak-toe/network"
)

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

	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(gm, w, r)
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
func serveWs(gm *game.Manager, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	session := gm.Session()

	client := network.NewClient(session.Hub(), conn)
	player := session.AddPlayer(client)
	player.Send(event.NewInit(player.Label, session.ID()))

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.WritePump()
	go client.ReadPump(session)

	if session.IsFull() {
		session.Start()
	}
}
