package jwttoken

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const Issuer = "duchng"

const (
	DefaultAccessTokenLifetime = 5 * time.Hour
)

type TokenType string

const (
	TypeDefault  TokenType = "default"
	TypePassword TokenType = "password"
	TypeSignIn   TokenType = "signin"
)

type Claims struct {
	jwt.RegisteredClaims

	FullName string `json:"fullName,omitempty"`
	UserId   int    `json:"userId,omitempty"`
	Email    string `json:"email,omitempty"`
	Role     string `json:"roles,omitempty"`
}
