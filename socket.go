package circle

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"net/http"
	"time"
)

var (
	Online int = 0
)

func OnOpen(ws *websocket.Conn) {
	Online = Online + 1
	Log()

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

		for {
			websocket.Message.Send(ws, <-ch)
			time.Sleep(10 * time.Millisecond)
		}
	}

	Online = Online - 1
	ws.Close()
	Log()
}

func Log() {
	fmt.Println(fmt.Sprintf("%d open connections.", Online))
}

func Start(port string) {
	fmt.Println("Starting...")

	http.Handle("/", websocket.Handler(OnOpen))

	err := http.ListenAndServe(port, nil)

	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
