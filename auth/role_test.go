package auth_test

import (
	"testing"

	"github.com/project-safari/zebra/auth"
	"github.com/stretchr/testify/assert"
)

// Test for keys.
func TestKey(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)
	k, e := auth.NewKey("a/b/c")
	assert.Nil(e)
	assert.NotNil(k)
	assert.True(k.Match("a/b/c"))
	assert.True(k.Match("a/b/c/d"))
	assert.False(k.Match("/b/c/d"))

	k, e = auth.NewKey("*")
	assert.NotNil(e)
	assert.Nil(k)

	k, e = auth.NewKey("a/[a-z0-9]*/c")
	assert.Nil(e)
	assert.NotNil(k)
	assert.True(k.Match("a/b/c"))
	assert.True(k.Match("a/k/c"))
	assert.False(k.Match("/b/c/d"))
}

// Tests for crud operations on priv.
func TestPriv(t *testing.T) { //nolint:funlen
	t.Parallel()

	assert := assert.New(t)
	p, e := auth.NewPriv("*", false, false, false, false)
	assert.Nil(p)
	assert.NotNil(e)
	p, e = auth.NewPriv("a/b/c", false, false, false, false)
	assert.Nil(p)
	assert.NotNil(e)

	p, e = auth.NewPriv("a/b/c", true, false, false, false)
	assert.Nil(e)
	assert.NotNil(p)
	assert.Equal("a/b/c:c", p.String())
	b, e := p.MarshalText()
	assert.Nil(e)
	assert.Equal(p.String(), string(b))
	assert.True(p.Create("a/b/c"))
	assert.False(p.Read("a/b/c"))
	assert.False(p.Update("a/b/c"))
	assert.False(p.Delete("a/b/c"))
	assert.False(p.Write("a/b/c"))

	assert.Nil(p.UnmarshalText([]byte("c/d/e:c,r,u,d")))
	assert.Equal("c/d/e:c,r,u,d", p.String())
	assert.NotNil(p.UnmarshalText([]byte("*:c,r,u,d")))
	assert.NotNil(p.UnmarshalText([]byte("c,r,u,d")))
	assert.NotNil(p.UnmarshalText([]byte("a/b/c")))
	assert.NotNil(p.UnmarshalText([]byte("a/b/c:xxx")))
	assert.NotNil(p.UnmarshalText([]byte("a/b/c:c,r,u,d,e,f")))

	p, e = auth.NewPriv("a/b/c", true, true, false, false)
	assert.Nil(e)
	assert.NotNil(p)
	assert.Equal("a/b/c:c,r", p.String())
	assert.True(p.Create("a/b/c"))
	assert.True(p.Read("a/b/c"))
	assert.False(p.Update("a/b/c"))
	assert.False(p.Delete("a/b/c"))
	assert.False(p.Write("a/b/c"))

	p, e = auth.NewPriv("a/b/c", true, true, true, false)
	assert.Nil(e)
	assert.NotNil(p)
	assert.Equal("a/b/c:c,r,u", p.String())
	assert.True(p.Create("a/b/c"))
	assert.True(p.Read("a/b/c"))
	assert.True(p.Update("a/b/c"))
	assert.False(p.Delete("a/b/c"))
	assert.False(p.Write("a/b/c"))

	p, e = auth.NewPriv("a/b/c", true, true, true, true)
	assert.Nil(e)
	assert.NotNil(p)
	assert.Equal("a/b/c:c,r,u,d", p.String())
	assert.True(p.Create("a/b/c"))
	assert.True(p.Read("a/b/c"))
	assert.True(p.Update("a/b/c"))
	assert.True(p.Delete("a/b/c"))
	assert.True(p.Write("a/b/c"))
	assert.False(p.Create("e/f/g"))
	assert.False(p.Read("e/f/g"))
	assert.False(p.Update("e/f/g"))
	assert.False(p.Delete("e/f/g"))
	assert.False(p.Write("e/f/g"))
}
