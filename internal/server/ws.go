package server

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type wsHandler struct {
	wsConnection *websocket.Conn
	messageChan  chan string
}

var (
	_ http.Handler = &wsHandler{}

	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
)

// NewWSHandler создает хэндлер для WebSocket соединения.
func newWSHandler(messageChan chan string) (http.Handler, error) {
	wsHandler := &wsHandler{
		messageChan: messageChan,
	}

	go wsHandler.proxyMessages()

	return wsHandler, nil
}

func (wsh *wsHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	// upgrade HTTP connection to a WebSocket connection
	wsConnection, err := upgrader.Upgrade(response, request, nil)
	if err != nil {
		panic(err)
	}

	wsh.wsConnection = wsConnection

	messageType, _, err := wsh.wsConnection.ReadMessage()
	if err != nil || messageType == websocket.CloseMessage {
		log.Printf("ws connection lost\n")
		return
	}

	// log.Printf("new '%s' client\n", settings.Name)
}

// Он смотрит в канал сообщений и переправляет его клиенту
func (wsh *wsHandler) proxyMessages() {
	for {
		message := <-wsh.messageChan

		if wsh.wsConnection == nil {
			log.Printf("message not sended via WS because of no connection\n")
			continue
		}

		if err := wsh.wsConnection.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
			panic(err)
		}
	}
}
