package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type WebsocketService struct {
	ctx context.Context
}

func NewWebsocketService() *WebsocketService {
	return &WebsocketService{
		ctx: context.Background(),
	}
}

func (o *WebsocketService) SendMassage(ws *websocket.Conn, param string) {
	for {
		// Read Massage From Client
		mt, massage, err := ws.ReadMessage()
		if err != nil {
			log.Printf("Failed Read Massage %s", err.Error())
			break
		}
		// If Clinet Massage is ping will return pong
		if string(massage) == "ping" {
			massage = []byte("pong")
		}

		if string(massage) == "hai" {
			massage = []byte(fmt.Sprintf("Hai , %v", param))
		}

		if string(massage) == "loop" {
			for _, i := range []string{"1", "2", "3"} {
				if err := ws.WriteMessage(mt, []byte(fmt.Sprintf("massage ke : %v", i))); err != nil {
					log.Printf("Failed Sent Massage error : %v", err)
					break
				}
			}
		}
		//Response massage to Client
		if err := ws.WriteMessage(mt, massage); err != nil {
			log.Printf("Failed Sent Massage error : %v", err)
			break
		}
	}
}

func (o *WebsocketService) Broadcast(ws *websocket.Conn, hub *Hub) {
	client := &Client{hub: hub, conn: ws, send: make(chan []byte)}
	client.hub.register <- client
	client.ID = GenerateID()
	client.Addr = client.conn.RemoteAddr().String()
	client.EnterAt = time.Now()

	go client.WritePump()
	go client.ReadPump()

	client.send <- []byte("Welcome")
}
