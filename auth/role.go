package auth

import (
	"bytes"
	"errors"
	"regexp"
	"strings"
)

// Potential errors concerning resource keys and priviledges.
var (
	ErrResourceKeyEmpty  = errors.New("resource key is empty")
	ErrInvalidPrivileges = errors.New("atleast one or atmost four privileges must be set")
)

// Possible crud priviledges.
type Priv struct {
	c bool
	r bool
	u bool
	d bool
	k *ResourceKey
}

// Create / generate a new privilege.
func NewPriv(k string, c bool, r bool, u bool, d bool) (*Priv, error) {
	rk, e := NewKey(k)
	if e != nil {
		return nil, e
	}

	if !c && !r && !u && !d {
		return nil, ErrInvalidPrivileges
	}

	return &Priv{
		c: c,
		r: r,
		u: u,
		d: d,
		k: rk,
	}, nil
}

func (p *Priv) String() string {
	crud := ":c,r,u,d"
	buf := bytes.NewBuffer(make([]byte, 0, len(p.k.key)+len(crud)))
	buf.WriteString(p.k.key)
	buf.WriteString(":")

	prev := false
	prevWrite := func(b bool, l string) {
		if prev && b {
			buf.WriteString(",")
			buf.WriteString(l)
		} else if b {
			buf.WriteString(l)
			prev = true
		}
	}
	prevWrite(p.c, "c")
	prevWrite(p.r, "r")
	prevWrite(p.u, "u")
	prevWrite(p.d, "d")

	return buf.String()
}

func (p *Priv) MarshalText() ([]byte, error) {
	return []byte(p.String()), nil
}

//nolint:cyclop
func (p *Priv) UnmarshalText(text []byte) error {
	t := string(text)
	t = strings.TrimSpace(t)

	if !strings.Contains(t, ":") {
		return ErrResourceKeyEmpty
	}

	s := strings.Split(t, ":")

	privs := strings.Split(s[1], ",")
	if len(privs) == 0 || len(privs) > 4 {
		return ErrInvalidPrivileges
	}

	k, err := NewKey(s[0])
	if err != nil {
		return err
	}

	p.k = k

	for _, priv := range privs {
		switch priv {
		case "c":
			p.c = true
		case "r":
			p.r = true
		case "u":
			p.u = true
		case "d":
			p.d = true
		default:
			return ErrInvalidPrivileges
		}
	}

	return nil
}

// Crud operation function for priv: read, returns boolean.
func (p *Priv) Read(key string) bool {
	return p.k.Match(key) && p.r
}

// Operation function for priv: write, returns boolean.
func (p *Priv) Write(key string) bool {
	return p.k.Match(key) && p.c && p.u && p.d
}

// Crud operation function for priv: update, returns boolean.
func (p *Priv) Update(key string) bool {
	return p.k.Match(key) && p.u
}

// Crud operation function for priv: update, returns boolean.
func (p *Priv) Create(key string) bool {
	return p.k.Match(key) && p.c
}

// Crud operation function for priv: delete, returns boolean.
func (p *Priv) Delete(key string) bool {
	return p.k.Match(key) && p.d
}

type ResourceKey struct {
	key string
	re  *regexp.Regexp
}

func NewKey(key string) (*ResourceKey, error) {
	re, err := regexp.Compile(key)
	if err != nil {
		return nil, err
	}

	return &ResourceKey{
		key: key,
		re:  re,
	}, nil
}

func (r *ResourceKey) Match(key string) bool {
	return r.re.MatchString(key)
}

type Role struct {
	Name       string  `json:"name"`
	Privileges []*Priv `json:"privileges"`
}

func (r *Role) Read(key string) bool {
	for _, priv := range r.Privileges {
		if priv.Read(key) {
			return true
		}
	}

	return false
}

// Operation function for role: write, returns boolean.
func (r *Role) Write(key string) bool {
	for _, priv := range r.Privileges {
		if priv.Write(key) {
			return true
		}
	}

	return false
}

// Crud operation function for role: create, returns boolean.
func (r *Role) Create(key string) bool {
	for _, priv := range r.Privileges {
		if priv.Create(key) {
			return true
		}
	}

	return false
}

// Crud operation function for role: update, returns boolean.
func (r *Role) Update(key string) bool {
	for _, priv := range r.Privileges {
		if priv.Update(key) {
			return true
		}
	}

	return false
}

// Crud operation function for role: delete, returns boolean.
func (r *Role) Delete(key string) bool {
	for _, priv := range r.Privileges {
		if priv.Delete(key) {
			return true
		}
	}

	return false
}
