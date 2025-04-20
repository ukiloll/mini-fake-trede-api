package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Id       string    `json:"id"`
	Username string    `json:"username"`
	jwt.RegisteredClaims//get all type in this struct
}

type Cookie struct{
	name string
	time time.Duration 
}

