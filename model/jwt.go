package model

import "github.com/golang-jwt/jwt/v5"

type JwtToken struct {
	UserID       int    `json:"user_id"`
	Username     string `json:"username"`
	Enable       bool   `json:"enable"`
	Admin        bool   `json:"admin"`
	TokenVersion string `json:"token_version"`
	jwt.RegisteredClaims
}
