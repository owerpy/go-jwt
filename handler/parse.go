package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
)

const tk string = "token"

func Parse(g *gin.Context) (string, error) {
	c, err := g.Request.Cookie(tk)
	//	fmt.Println(c.Expires.Before(time.Now()), c.Expires.After(time.Now()), c.Expires)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			g.Writer.WriteHeader(http.StatusUnauthorized)
			return "", err
		}

		g.Writer.WriteHeader(http.StatusBadRequest)
		return "", err
	}

	tknStr := c.Value

	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (any, error) {
		return jwtKey, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			g.Writer.WriteHeader(http.StatusUnauthorized)
			return "", err
		}

		g.Writer.WriteHeader(http.StatusBadRequest)
		return "", err
	}

	if !tkn.Valid || c.Value == "" {
		g.Writer.WriteHeader(http.StatusUnauthorized)
		return "", err
	}

	return claims.Username, nil
}
