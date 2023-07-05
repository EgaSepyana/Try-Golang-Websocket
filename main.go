package main

import (
	"gosock/src/config"
	"gosock/src/controller"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/subosito/gotenv"
)

func init() {
	if err := gotenv.Load(".env"); err != nil {
		log.Println(err)
	}
}

func main() {
	appRouter := ":" + os.Getenv(config.SERVER_PORT)

	router := gin.Default()

	BasePath := "/api/v1"

	apiV1 := router.Group(BasePath)

	controller.NewWebsocketController(apiV1)

	log.Fatalln(router.Run(appRouter))
}
