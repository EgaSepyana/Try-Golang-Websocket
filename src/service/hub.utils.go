package service

import (
	"encoding/json"
	"fmt"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Resgitered Client
	clients map[*Client]bool

	// Inbound Massage from the client
	broadcasts chan []byte

	// Register request from the client
	register chan *Client

	// unregister request from client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcasts: make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (o *Hub) Run() {
	for {
		select {
		case client := <-o.register:
			clientId := client.ID
			for client := range o.clients {
				msg := []byte(fmt.Sprintf("Ada yang join coy ( ID : %s) ", clientId))
				client.send <- msg
			}

			o.clients[client] = true

		case client := <-o.unregister:
			fmt.Println(client)
			clinetId := client.ID
			if _, ok := o.clients[client]; ok {
				delete(o.clients, client)
				close(client.send)
			}
			for client := range o.clients {
				msg := []byte(fmt.Sprintf("Ada yang Left Coy ( ID : %s)", clinetId))
				client.send <- msg
			}

		case userMassage := <-o.broadcasts:
			var data map[string][]byte
			json.Unmarshal(userMassage, &data)

			for client := range o.clients {
				if client.ID == string(data["id"]) {
					continue
				}

				select {
				case client.send <- []byte(fmt.Sprintf("( ID : %s) : %s", string(data["id"]), string(data["messege"]))):
				default:
					close(client.send)
					delete(o.clients, client)
				}
			}
		}
	}
}
