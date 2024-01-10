package handler

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

// Create the JWT key used to create the signature
var jwtKey = []byte("my_secret_key")

// For simplification, we're storing the users information as an in-memory map in our code
var users = map[string]string{
	"user1": "password1",
	"user2": "password2",
}

// Create a struct to read the username and password from the request body
type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

// Create a struct that will be encoded to a JWT.
// We add jwt.RegisteredClaims as an embedded type, to provide fields like expiry time
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// Create the Signin handler
func SignIn(g *gin.Context) {
	var creds Credentials
	// Get the JSON body and decode into credentials
	err := json.NewDecoder(g.Request.Body).Decode(&creds)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		g.Writer.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get the expected password from our in memory map
	expectedPassword, ok := users[creds.Username]

	// If a password exists for the given user
	// AND, if it is the same as the password we received, the we can move ahead
	// if NOT, then we return an "Unauthorized" status
	if !ok || expectedPassword != creds.Password {
		g.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Declare the expiration time of the tk
	// here, we have kept it as 5 minutes
	expirationTime := time.Now().Add(4000 * time.Minute)
	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{
		Username: creds.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// Declare the tk with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		g.Writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Finally, we set the client cookie for "tk" as the JWT we just generated
	// we also set an expiry time which is the same as the tk itself
	http.SetCookie(g.Writer, &http.Cookie{
		Name:     tk,
		Value:    tokenString,
		Expires:  expirationTime,
		HttpOnly: false,
	})
}
