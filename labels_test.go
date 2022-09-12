package zebra_test

import (
	"context"
	"testing"

	"github.com/project-safari/zebra"
	"github.com/stretchr/testify/assert"
)

func TestLabels(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)
	labels := zebra.Labels{
		"color": "red",
	}
	assert.NotNil(labels.Validate())

	assert.True(labels.MatchEqual("color", "red"))
	assert.False(labels.MatchEqual("color", "blue"))
	assert.False(labels.MatchEqual("blah", "blee"))

	assert.True(labels.MatchNotEqual("color", "blue"))
	assert.False(labels.MatchNotEqual("color", "red"))
	assert.False(labels.MatchNotEqual("blah", "blee"))

	assert.True(labels.MatchIn("color", "red", "blue"))
	assert.True(labels.MatchIn("color", "red"))
	assert.False(labels.MatchIn("color", "blue", "green"))
	assert.False(labels.MatchIn("blah"))

	assert.False(labels.MatchNotIn("color", "red", "blue"))
	assert.False(labels.MatchNotIn("color", "red"))
	assert.True(labels.MatchNotIn("color", "blue", "green"))
	assert.False(labels.MatchNotIn("blah"))

	// Add a new label and make sure that it doesnt effect existing label
	labels.Add("shape", "sphere").Add("planet", "mars")

	assert.True(labels.MatchEqual("color", "red"))
	assert.False(labels.MatchEqual("color", "blue"))
	assert.False(labels.MatchEqual("blah", "blee"))

	assert.True(labels.MatchNotEqual("color", "blue"))
	assert.False(labels.MatchNotEqual("color", "red"))
	assert.False(labels.MatchNotEqual("blah", "blee"))

	assert.True(labels.MatchIn("color", "red", "blue"))
	assert.True(labels.MatchIn("color", "red"))
	assert.False(labels.MatchIn("color", "blue", "green"))
	assert.False(labels.MatchIn("blah"))

	assert.False(labels.MatchNotIn("color", "red", "blue"))
	assert.False(labels.MatchNotIn("color", "red"))
	assert.True(labels.MatchNotIn("color", "blue", "green"))
	assert.False(labels.MatchNotIn("blah"))
}

func TestIsIn(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	list := []string{"hi", "hello", "goodbye"}

	assert.True(zebra.IsIn("hello", list))
	assert.False(zebra.IsIn("hey", list))
}

func TestLabelsValidation(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	dummy, _ := dummyType()
	resOne := zebra.NewBaseResource(dummy, "dummy", "dummy", "dummy")
	assert.Nil(resOne.Validate(context.Background()))
}
