package zebra

import (
	"github.com/google/uuid"
	"github.com/project-safari/zebra/status"
)

func NewBaseResource(resType string, labels Labels) *BaseResource {
	id := GenerateID()

	if resType == "" {
		resType = "BaseResource"
	}

	return &BaseResource{
		ID:     id,
		Type:   resType,
		Labels: labels,
		Status: status.DefaultStatus(),
	}
}

func NewNamedResource(name string, resType string, labels Labels) *NamedResource {
	if resType == "" {
		resType = "NamedResource"
	}

	return &NamedResource{
		BaseResource: *NewBaseResource(resType, labels),
		Name:         name,
	}
}

func NewCredentials(name string, keys map[string]string, labels Labels) *Credentials {
	// Ensure name is set, and returned resource will be valid
	if name == "" {
		name = "unknown"
	}

	// Ensure keys are not nil
	if keys == nil {
		keys = map[string]string{}
	}

	ret := &Credentials{
		NamedResource: *NewNamedResource(name, "Credentials", labels),
		Keys:          keys,
	}

	return ret
}

func GenerateID() string {
	return uuid.New().String()
}

// Return if val is in string list.
func IsIn(val string, list []string) bool {
	for _, v := range list {
		if val == v {
			return true
		}
	}

	return false
}
