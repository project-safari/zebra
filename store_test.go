package zebra_test

import (
	"encoding/json"
	"testing"

	"github.com/project-safari/zebra"
	"github.com/stretchr/testify/assert"
)

// test for store validation.
func TestValidateStore(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	q := new(zebra.Query)
	assert.NotNil(q)
	assert.NotNil(q.Validate())

	q.Op = 8
	assert.NotNil(q.Validate())

	q.Op = zebra.MatchEqual
	q.Key = "key"
	q.Values = []string{"value"}
	assert.Nil(q.Validate())
}

func TestMarshalQuery(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	q := new(zebra.Query)

	q.Op = zebra.MatchEqual
	q.Key = "key"
	q.Values = []string{"value1", "value2"}

	ret, err := json.Marshal(q)
	assert.Nil(err)
	assert.Equal(`{"key":"key","op":"==","values":["value1","value2"]}`, string(ret))

	// bad marshal
	q.Op = 9
	_, err = json.Marshal(q)
	assert.NotNil(err)
}

func TestUnmarshalQuery(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	// bad unmarshal
	bytes := []byte(`{"key":"","op":"blahblah","values":[]}`)
	q := new(zebra.Query)

	err := json.Unmarshal(bytes, &q)
	assert.NotNil(err)

	// good unmarshal
	bytes = []byte(`{"key":"key","op":"==","values":["value1","value2"]}`)
	q = new(zebra.Query)

	assert.Nil(json.Unmarshal(bytes, &q))
	assert.Equal("key", q.Key)
	assert.Equal(zebra.MatchEqual, q.Op)
	assert.Equal(2, len(q.Values))
	assert.True(q.Values[0] == "value1" || q.Values[1] == "value1")
	assert.True(q.Values[0] == "value2" || q.Values[1] == "value2")
}
