package controller

import (
	"fmt"
	"log"
	"net/http"

	"gosock/src/service"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var userActive int = 0

type WebsocketController struct {
	router   *gin.RouterGroup
	upgrader websocket.Upgrader
	service  *service.WebsocketService
	hub      *service.Hub
}

func NewWebsocketController(router *gin.RouterGroup) *WebsocketController {
	o := &WebsocketController{
		router: router,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				// origin := r.Header.Get("Origin")
				// fmt.Println(origin)
				return true
			},
		},
		service: service.NewWebsocketService(),
		hub:     service.NewHub(),
	}

	go o.hub.Run()

	websocket := o.router.Group("/ws")
	websocket.GET("/", o.SendMassage)
	websocket.GET("/broadcast", o.Broadcast)

	return o
}

func (o *WebsocketController) Connect(w http.ResponseWriter, r *http.Request, h http.Header) (con *websocket.Conn) {
	con, err := o.upgrader.Upgrade(w, r, h)
	if err != nil {
		log.Printf("failed Connect error : %s", err.Error())
		return
	}
	userActive++
	fmt.Printf("Clinet Connected %v \n", userActive)
	return
}

func (o *WebsocketController) Broadcast(c *gin.Context) {
	ws := o.Connect(c.Writer, c.Request, nil)

	o.service.Broadcast(ws, o.hub)
}

func (o *WebsocketController) SendMassage(c *gin.Context) {
	//Upgrade GET request to WebSocket protocol
	ws := o.Connect(c.Writer, c.Request, nil)
	defer ws.Close()

	o.service.SendMassage(ws, c.Query("username"))
}
