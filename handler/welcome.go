package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func Welcome(g *gin.Context) {
	userName, err := Parse(g)
	if err != nil {
		return
	}
	// Finally, return the welcome message to the user, along with their
	// username given in the tk
	_, err = g.Writer.Write([]byte(fmt.Sprintf("Welcome %s!", userName)))
	if err != nil {
		return
	}
}
