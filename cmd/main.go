package main

import (
	"github.com/gin-gonic/gin"
	"go-jwt/handler"
)

const url = "localhost:8000"

func main() {
	//var g *gin.Engine
	g := gin.New()
	// we will implement these handlers in the next sections
	g.POST("/sign-in", handler.SignIn)
	g.GET("/welcome", handler.Welcome)
	g.GET("/refresh", handler.Refresh)
	g.GET("/logout", handler.Logout)

	// start the server on port 8000
	if err := g.Run(url); err != nil {
		return
	}
}
