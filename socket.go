package circle

import (
	"github.com/azer/logger"
	"code.google.com/p/go.net/websocket"
	"net/http"
	"time"
	"fmt"
)

var (
	log = logger.New("socket")
	Online int = 0
)

func OnOpen(ws *websocket.Conn) {
	Online = Online + 1

	log.Info("%d online users.", Online)

	key := fmt.Sprintf("Socket#%d", Online)
	timer := log.Timer()

	var (
		user []byte
		ch   chan string
	)

	for {
		if err := websocket.Message.Receive(ws, &user); err != nil {
			break
		}

		ch = make(chan string)
		go SubscribeTo(string(user), ch)
		go Receive(ws, ch)
	}

	Online = Online - 1
	ws.Close()

	timer.End("%s got closed.", key)
}

func Receive(ws *websocket.Conn, ch chan string) {
	for {
		websocket.Message.Send(ws, <-ch)
		time.Sleep(10 * time.Millisecond)
	}

	close(ch)
}

func Start(port string) {
	log.Info("Starting the server...")

	http.Handle("/", websocket.Handler(OnOpen))

	err := http.ListenAndServe(port, nil)

	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
