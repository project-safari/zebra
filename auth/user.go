package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/project-safari/zebra"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrKeyEmpty      = errors.New("ssh key is empty")
	ErrPasswordEmpty = errors.New("password hash is empty")
	ErrRoleEmpty     = errors.New("role is empty")
)

func UserType() zebra.Type {
	return zebra.Type{
		Name:        "User",
		Description: "zebra user",
		Constructor: func() zebra.Resource { return new(User) },
	}
}

type User struct {
	zebra.NamedResource
	Key          *RsaIdentity `json:"key"`
	PasswordHash string       `json:"passwordHash"`
	Role         *Role        `json:"role"`
}

// Validate returns an error if the given Datacenter object has incorrect values.
// Else, it returns nil.
func (u *User) Validate(ctx context.Context) error {
	if u.Key == nil {
		return ErrKeyEmpty
	}

	if u.Role == nil {
		return ErrRoleEmpty
	}

	if u.PasswordHash == "" {
		return ErrPasswordEmpty
	}

	return u.NamedResource.Validate(ctx)
}

const SharedSecret = "खुल जा सिम सिम" //nolint:gosec

func (u *User) Authenticate(token string) error {
	return u.Key.Verify([]byte(SharedSecret), []byte(token), nil)
}

func (u *User) AuthenticatePassword(password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return fmt.Errorf("bad password: %w", err)
	}

	return nil
}

func (u *User) Create(resource string) bool {
	return u.Role.Create(resource)
}

func (u *User) Read(resource string) bool {
	return u.Role.Read(resource)
}

func (u *User) Write(resource string) bool {
	return u.Role.Write(resource)
}

func (u *User) Delete(resource string) bool {
	return u.Role.Delete(resource)
}

func (u *User) Update(resource string) bool {
	return u.Role.Update(resource)
}

func HashPassword(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	return string(hash)
}
