package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func Logout(g *gin.Context) {
	// immediately clear the tk cookie
	http.SetCookie(g.Writer, &http.Cookie{
		Name:     tk,
		Expires:  time.Now(),
		HttpOnly: false,
		MaxAge:   -1,
		Value:    "",
	})
}
