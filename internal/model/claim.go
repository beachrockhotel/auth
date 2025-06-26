package model

import (
	"github.com/beachrockhotel/auth/pkg/auth_v1"
	"github.com/dgrijalva/jwt-go"
)

const (
	ExamplePath = "/auth_v1.AuthV1/Get"
)

type UserClaims struct {
	jwt.StandardClaims
	Name string       `json:"name"`
	Role auth_v1.Role `json:"role"`
}
