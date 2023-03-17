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

// Errors that occur when something is wrong about a resource.
var (
	// ErrNotFound happens if the resource does not exist in the store.
	ErrNotFound = errors.New("resource not found in store")
	// ErrInvalidResource happens if the resource is invalid.
	ErrInvalidResource = errors.New("create/delete on invalid resource")
	// ErrInvalidQuery happens if the query is invalid.
	ErrInvalidQuery = errors.New("invalid query")
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
}

// Validate function for the query.
// It returns an error or nil in the absence thereof.
func (q *Query) Validate() error {
	if (q.Op == MatchEqual || q.Op == MatchNotEqual) && len(q.Values) != 1 {
		return ErrInvalidQuery
	}

	if q.Op > MatchNotIn {
		return ErrInvalidQuery
	}

	return nil
}

// Function on a pointer to the Operator struct to marshal.
//
// It returns a byte array and an error or nil, in the absence thereof.
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

// Function on a pointer to the Operation struct to unMarshal.
//
// It takes in a byte array.
// It returns an error or nil, in the absence thereof.
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
