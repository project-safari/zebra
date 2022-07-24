package zebra_test

import (
	"testing"

	"github.com/project-safari/zebra"
	"github.com/stretchr/testify/assert"
)

func TestLabelsValidation(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	// first test - with a correct default label
	resOne := new(zebra.BaseResource)
	mapOne := map[string]string{
		"group": "Americas",
		"color": "red",
	}

	resOne.Labels = mapOne
	assert.Nil(resOne.Labels.Validate())

	// second test - with an incorrect default label
	resTwo := new(zebra.BaseResource)
	mapTwo := map[string]string{
		"letter": "alpha",
		"color":  "blue",
	}

	resTwo.Labels = mapTwo

	assert.NotNil(resTwo.Labels.Validate())
	assert.Equal(zebra.ErrLabel, resTwo.Labels.Validate())
}
