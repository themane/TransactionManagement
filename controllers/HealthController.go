package controllers

import "github.com/gin-gonic/gin"

// Ping godoc
// @Summary Pings the server
// @Description Pings the server for checking the health of the server
// @Tags root
// @Accept */*
// @Produce json
// @Success 200 {object} PongResponse
// @Router /ping [get]
func Ping(c *gin.Context) {
	response := PongResponse{Message: "pong"}
	c.JSON(200, &response)
}

type PongResponse struct {
	Message string `json:"message" example:"pong"`
}
