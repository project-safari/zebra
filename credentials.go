package zebra

import (
	"errors"
	"strings"
	"unicode"
)

var (
	ErrPassLen           = errors.New("password is less than 12 characters long")
	ErrPassCase          = errors.New("password does not contain both upper and lowercase")
	ErrPassNum           = errors.New("password does not contain a number")
	ErrPassSpecial       = errors.New("password does not contain a special character")
	ErrUnknownCredential = errors.New("unknown credential type")
	ErrNoKeys            = errors.New("keys is nil")
)

// Credentials represents a named resource that has a set of keys (where each key is
// an authentication method) with corresponding values (where each value is the
// information to store for the authentication method).
type Credentials struct {
	LoginID string            `json:"loginId"`
	Keys    map[string]string `json:"keys,omitempty"`
}

// Function that returns a credentials struct with login ID and a pair of keys.
func NewCredentials(loginID string) Credentials {
	return Credentials{
		LoginID: loginID,
		Keys:    map[string]string{},
	}
}

// Operation function on the credentials - add.
//
// This function adds a new set of key-value pairs with
// password and corresponding ssh key.
//
// It returns an error or nil in the absence thereof.
func (c Credentials) Add(key string, value string) error {
	validKeys := map[string]string{"password": "", "ssh-key": ""}
	if _, ok := validKeys[key]; !ok {
		return ErrUnknownCredential
	}

	c.Keys[key] = value

	return nil
}

// Validate returns an error if the given Credentials object has incorrect values.
// Else, it returns nil.
func (c Credentials) Validate() error {
	if c.LoginID == "" {
		return ErrIDEmpty
	}

	if len(c.Keys) == 0 {
		return ErrNoKeys
	}

	keyValidators := map[string]func(string) error{
		"password": ValidatePassword,
		"ssh-key":  ValidateSSHKey,
	}

	for keyType, key := range c.Keys {
		v := keyValidators[keyType]
		if err := v(key); err != nil {
			return err
		}
	}

	return nil
}

// Check to make sure password follows rules.
// 1. At least 12 characters long.
// 2. Contains upper and lowercase letters.
// 3. Contains at least 1 number.
// 4. Contains at least 1 special character.
func ValidatePassword(password string) error { //nolint:cyclop
	minLength := 12

	var okLen, okLower, okUpper, okNum, okSpec bool

	for _, char := range password {
		switch {
		case unicode.IsLower(char) && unicode.IsLetter(char):
			okLower = true
		case unicode.IsUpper(char) && unicode.IsLetter(char):
			okUpper = true
		case unicode.IsDigit(char):
			okNum = true
		case strings.Contains("!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~", string(char)): //nolint:gocritic
			okSpec = true
		}
	}

	okLen = len(password) >= minLength

	switch {
	case !okLen:
		return ErrPassLen
	case (!okLower || !okUpper):
		return ErrPassCase
	case !okNum:
		return ErrPassNum
	case !okSpec:
		return ErrPassSpecial
	}

	return nil
}

// TO BE IMPLEMENTED AT A LATER TIME.
// Check to make sure SSH key follows specified rules.
func ValidateSSHKey(key string) error {
	return nil
}
