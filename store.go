package zebra

import (
	"errors"
)

type Operator uint8

// Constants defined for QueryOperator type.
const (
	MatchEqual Operator = iota
	MatchNotEqual
	MatchIn
	MatchNotIn
)

// Command struct for label queries.
type Query struct {
	Key    string   `json:"key"`
	Op     Operator `json:"op"`
	Values []string `json:"values"`
}

var (
	ErrNotFound        = errors.New("resource not found in store")
	ErrInvalidResource = errors.New("create/delete on invalid resource")
	ErrInvalidQuery    = errors.New("invalid query")
)

// Store interface requires basic store functionalities.
type Store interface {
	Initialize() error
	Wipe() error
	Clear() error
	Load() (*ResourceMap, error)
	Create(res Resource) error
	Delete(res Resource) error
	Query() *ResourceMap
	QueryUUID(uuids []string) *ResourceMap
	QueryType(types []string) *ResourceMap
	QueryLabel(query Query) (*ResourceMap, error)
	QueryProperty(query Query) (*ResourceMap, error)
	FreeResources(reslist []string)
}

func (q *Query) Validate() error {
	if (q.Op == MatchEqual || q.Op == MatchNotEqual) && len(q.Values) != 1 {
		return ErrInvalidQuery
	}

	if q.Op > MatchNotIn {
		return ErrInvalidQuery
	}

	return nil
}

func (o *Operator) MarshalText() ([]byte, error) {
	opMap := map[Operator]string{
		MatchEqual:    "==",
		MatchNotEqual: "!=",
		MatchIn:       "in",
		MatchNotIn:    "notin",
	}

	opVal, ok := opMap[*o]
	if !ok {
		return []byte(""), ErrInvalidQuery
	}

	return []byte(opVal), nil
}

func (o *Operator) UnmarshalText(data []byte) error {
	opMap := map[string]Operator{
		"==":    MatchEqual,
		"!=":    MatchNotEqual,
		"in":    MatchIn,
		"notin": MatchNotIn,
	}

	op, ok := opMap[string(data)]
	if !ok {
		return ErrInvalidQuery
	}

	*o = op

	return nil
}
