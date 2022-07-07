package auth

import (
	"context"
	"errors"

	"github.com/project-safari/zebra"
)

var (
	ErrKeyEmpty  = errors.New("key is empty")
	ErrRoleEmpty = errors.New("role is empty")
)

type User struct {
	zebra.NamedResource
	Key  *RsaIdentity `json:"key"`
	Role *Role        `json:"role"`
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

	return u.NamedResource.Validate(ctx)
}

const SharedSecret = "खुल जा सिम सिम" //nolint:gosec

func (u *User) Authenticate(token string) error {
	return u.Key.Verify([]byte(SharedSecret), []byte(token), nil)
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
