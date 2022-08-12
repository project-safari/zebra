package pkg_test

import (
	"testing"

	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/cmd/herd/pkg"
	"github.com/stretchr/testify/assert"
)

// tests for generation of system.group labels.
//
// the tested function generates system.group labels if the given resource does not have any such labels.
func TestGroupLabel(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)

	labels := pkg.CreateLabels()

	// test for generating group label based on address.
	assert.NotNil(pkg.GroupLabels(labels, pkg.Addresses()))

	// test to see if group is created for given address.
	groupTest := pkg.GroupLabels(labels, "Mexico")
	assert.True(groupTest.MatchEqual("system.group", "Mexico"))
}

func TestLabelGeneration(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)

	resource := new(zebra.BaseResource)

	resource.Type = "VM"

	grouped := pkg.GroupVal(resource)

	assert.NotNil(grouped)
}
