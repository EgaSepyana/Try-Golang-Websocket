package service

import (
	"bytes"
	"encoding/json"
	"gosock/src/models"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write massage to peer
	writeWait = 10 * time.Second

	// Time allowed to read the next pong massage form the peer
	pongWait = 60 * time.Second

	// Send pings to peer with this periode. Must be less than pong wait
	pingPeriode = (pongWait * 9) / 10

	// Maximum massage size allow from peer
	maxMessageSize = 512
)

var (
	newLine = []byte{'\n'}
	space   = []byte{' '}
)

type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte

	models.User
}

func (o *Client) ReadPump() {
	defer func() {
		o.hub.unregister <- o
		o.conn.Close()
	}()

	o.conn.SetReadLimit(maxMessageSize)
	o.conn.SetReadDeadline(time.Now().Add(pongWait))
	o.conn.SetPongHandler(func(string) error {
		o.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, massage, err := o.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error : %v", err)
			}

			break
		}

		massage = bytes.TrimSpace(bytes.Replace(massage, newLine, space, -1))
		data := map[string][]byte{
			"messege": massage,
			"id":      []byte(o.ID),
		}
		usermassage, _ := json.Marshal(data)
		o.hub.broadcasts <- usermassage
	}
}

func (o *Client) WritePump() {
	tinker := time.NewTicker(pingPeriode)
	defer func() {
		tinker.Stop()
		o.conn.Close()
	}()

	for {
		select {
		case messege, ok := <-o.send:
			if !ok {
				o.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			writer, err := o.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Println(err.Error())
				return
			}

			writer.Write(messege)

			// Add queued chat messages to the current websocket message.
			n := len(o.send)
			for i := 0; i < n; i++ {
				writer.Write(newLine)
				writer.Write(<-o.send)
			}

			if err := writer.Close(); err != nil {
				log.Printf("close writer Error : %v", err.Error())
				return
			}

		case <-tinker.C:
			o.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := o.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("write Messege in deadline Error : %v", err.Error())
				return
			}
		}
	}
}
