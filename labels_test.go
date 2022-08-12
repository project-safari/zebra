package zebra_test

import (
	"testing"

	"github.com/project-safari/zebra"
	"github.com/stretchr/testify/assert"
)

// various tests for labels.
// include tests for matching equality / inequality of labels.
// matching existence of a label key.
// tests are performed for ovarious quantities labels in the labels map.
func TestLabels(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)
	labels := zebra.Labels{
		"color": "red",
	}

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
