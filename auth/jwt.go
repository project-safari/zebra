package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// ErrInvalidToken is an error that occurs if the jwt token is invalid.
var ErrInvalidToken = errors.New("invalid jwt token")

// TokenDuration is a constant that sets the duration of a token.
const TokenDuration = time.Minute * 10

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

// Operation function on a pointer to Claims - Create.
//
// The function takes in a resorce, as a string.
// It returns a boolean value.
func (claims *Claims) Create(resource string) bool {
	return claims.Role.Create(resource)
}

// Operation function on a pointer to Claims - Read.
//
// The function takes in a resorce, as a string.
// It returns a boolean value.
func (claims *Claims) Read(resource string) bool {
	return claims.Role.Read(resource)
}

// Operation function on a pointer to Claims - Write.
//
// The function takes in a resorce, as a string.
// It returns a boolean value.
func (claims *Claims) Write(resource string) bool {
	return claims.Role.Write(resource)
}

// Operation function on a pointer to Claims - Delete.
//
// The function takes in a resorce, as a string.
// It returns a boolean value.
func (claims *Claims) Delete(resource string) bool {
	return claims.Role.Delete(resource)
}

// Operation function on a pointer to Claims - Update.
//
// The function takes in a resorce, as a string.
// It returns a boolean value.
func (claims *Claims) Update(resource string) bool {
	return claims.Role.Update(resource)
}

// Operation function on a pointer to Claims - JWT.
//
// It createsa new JWT token.
//
// The function takes in a key, as a string.
// It returns a string.
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
