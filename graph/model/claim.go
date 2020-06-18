package model

import jwt "github.com/dgrijalva/jwt-go"

type Claim struct {
	PrivateKey string `json:"private_key"`
	jwt.StandardClaims
}
