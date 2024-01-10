package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func Refresh(g *gin.Context) {
	// (BEGIN) The code until this point is the same as the first part of the `Welcome` route
	c, err := g.Request.Cookie(tk)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			g.Writer.WriteHeader(http.StatusUnauthorized)
			return
		}
		g.Writer.WriteHeader(http.StatusBadRequest)
		return
	}
	tknStr := c.Value
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (any, error) {
		return jwtKey, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			g.Writer.WriteHeader(http.StatusUnauthorized)
			return
		}
		g.Writer.WriteHeader(http.StatusBadRequest)
		return
	}
	if !tkn.Valid {
		g.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}
	// (END) The code until this point is the same as the first part of the `Welcome` route

	// We ensure that a new tk is not issued until enough time has elapsed
	// In this case, a new tk will only be issued if the old tk is within
	// 30 seconds of expiry. Otherwise, return a bad request status
	if time.Until(claims.ExpiresAt.Time) > 30*time.Second {
		g.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	// Now, create a new tk for the current use, with a renewed expiration time
	expirationTime := time.Now().Add(5 * time.Minute)
	claims.ExpiresAt = jwt.NewNumericDate(expirationTime)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		g.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set the new tk as the users `tk` cookie
	http.SetCookie(g.Writer, &http.Cookie{
		Name:    tk,
		Value:   tokenString,
		Expires: expirationTime,
	})
}
