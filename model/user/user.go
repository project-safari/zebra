package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/auth"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrKeyEmpty      = errors.New("ssh key is empty")
	ErrPasswordEmpty = errors.New("password hash is empty")
	ErrRoleEmpty     = errors.New("role is empty")
)

func Type() zebra.Type {
	return zebra.Type{Name: "system.user", Description: "zebra system user"}
}

func Empty() zebra.Resource {
	user := new(User)
	user.Meta.Type = Type()

	return user
}

type User struct {
	zebra.BaseResource
	Key          *auth.RsaIdentity `json:"key"`
	PasswordHash string            `json:"passwordHash"`
	Role         *auth.Role        `json:"role"`
	Email        string            `json:"email"`
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

	return u.BaseResource.Validate(ctx)
}

func (u *User) Authenticate(token string) error {
	return u.Key.Verify([]byte(u.Email), []byte(token), nil)
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

// create new user data.
func NewUser(name string, email string, pwd string, key *auth.RsaIdentity, role *auth.Role) *User {
	return &User{
		BaseResource: *zebra.NewBaseResource(Type(), name, name, "users"),
		Key:          key,
		PasswordHash: HashPassword(pwd),
		Role:         role,
		Email:        email,
	}
}

func MockUser(num int) []zebra.Resource {
	rs := make([]zebra.Resource, 0, num)

	for i := 1; i <= num; i++ {
		k, _ := auth.Generate()
		p, _ := auth.NewPriv("", true, true, true, true)
		r := NewUser(
			fmt.Sprintf("mock-user-%d", i),
			fmt.Sprintf("mock-user-%d@zebra.local", i),
			fmt.Sprintf("UserSecret%d!!!", i),
			k.Public(),
			&auth.Role{
				Name:       "admin",
				Privileges: []*auth.Priv{p},
			},
		)

		rs = append(rs, r)
	}

	return rs
}
