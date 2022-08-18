package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// ErrInvalidToken is a potential error concerning tokens.
// ErrInvalidToken returns an error if the token is not valid.
var ErrInvalidToken = errors.New("invalid jwt token")

const TokenDuration = time.Minute * 10

// Struct that helps manage claims.
// It contains StandardClaims, a Role, and an email address.
type Claims struct {
	jwt.StandardClaims
	Role  *Role  `json:"role"`
	Email string `json:"email"`
}

func NewClaims(issuer string, subject string, role *Role, email string) *Claims {
	claims := new(Claims)

	claims.Issuer = issuer
	claims.Subject = subject
	claims.ExpiresAt = time.Now().Add(TokenDuration).Unix()
	claims.Role = role
	claims.Email = email

	return claims
}

// Create claims, return boolean.
func (claims *Claims) Create(resource string) bool {
	return claims.Role.Create(resource)
}

// Read claims, return boolean.
func (claims *Claims) Read(resource string) bool {
	return claims.Role.Read(resource)
}

// Write claims, return boolean.
func (claims *Claims) Write(resource string) bool {
	return claims.Role.Write(resource)
}

// Delete existing claims, return boolean.
func (claims *Claims) Delete(resource string) bool {
	return claims.Role.Delete(resource)
}

// Update existing claims, return boolean.
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
		func(token *jwt.Token) (interface{}, error) {
			return []byte(key), nil
		})
	// Error parsing token
	if err != nil {
		return nil, err
	}

	if !tkn.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
