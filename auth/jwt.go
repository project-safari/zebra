package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var ErrInvalidToken = errors.New("invalid jwt token")

const TokenDuration = time.Minute * 10

type Claims struct {
	jwt.StandardClaims
	Role *Role `json:"role"`
}

func NewClaims(issuer string, subject string, role *Role) *Claims {
	claims := new(Claims)

	claims.Issuer = issuer
	claims.Subject = subject
	claims.ExpiresAt = time.Now().Add(TokenDuration).Unix()
	claims.Role = role

	return claims
}

func (claims *Claims) Create(resource string) bool {
	return claims.Role.Create(resource)
}

func (claims *Claims) Read(resource string) bool {
	return claims.Role.Read(resource)
}

func (claims *Claims) Write(resource string) bool {
	return claims.Role.Write(resource)
}

func (claims *Claims) Delete(resource string) bool {
	return claims.Role.Delete(resource)
}

func (claims *Claims) Update(resource string) bool {
	return claims.Role.Update(resource)
}

func (claims *Claims) JWT(key string) string {
	tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtStr, _ := tkn.SignedString([]byte(key))

	return jwtStr
}

func FromJWT(token string, key string) (*Claims, error) {
	claims := new(Claims)
	tkn, err := jwt.ParseWithClaims(token, claims,
		func(token *jwt.Token) (interface{}, error) { return []byte(key), nil })

	if !tkn.Valid {
		return nil, ErrInvalidToken
	}

	return claims, err
}
