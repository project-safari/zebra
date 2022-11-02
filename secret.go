package zebra

type Secret struct {
	secret string
}

// Function to Marshal.
//
// Returns a byte array and an error, or nil in the absence thereof.
func (s *Secret) MarshalText() ([]byte, error) {
	return []byte("*****"), nil
}

// Function to unMarshal.
//
// Returns an error, or nil in the absence thereof.
func (s *Secret) UnmarshalText(text []byte) error {
	s.secret = string(text)

	return nil
}
