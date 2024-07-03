package main

import (
	"github.com/gin-gonic/gin"
	"log"
)

// @title Transaction Management Server
// @version 1.0.0
// @description This is the server for any Transaction Management system
// @termsOfService http://swagger.io/terms/

// @contact.name Devashish Gupta
// @contact.email devagpta@gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @schemes https
func main() {
	r := gin.Default()

	r.GET("/ping", Ping)

	err := r.Run()
	if err != nil {
		log.Println("Error in starting server")
		return
	}
}

func Ping(c *gin.Context) {
	response := PongResponse{Message: "pong"}
	c.JSON(200, &response)
}

type PongResponse struct {
	Message string `json:"message" example:"pong"`
}
