package zebra

import (
	"github.com/google/uuid"
	"github.com/project-safari/zebra/status"
)

const DefaultMaxDuration = 4

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

func NewCredentials(name string, labels Labels) *Credentials {
	namedRes := new(NamedResource)

	namedRes.BaseResource = *NewBaseResource("Credentials", labels)

	// Ensure name is set, and returned resource will be valid
	if name == "" {
		name = "unknown"
	}

	namedRes.Name = name

	ret := &Credentials{
		NamedResource: *namedRes,
		// some labels.
		Keys: labels,
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
