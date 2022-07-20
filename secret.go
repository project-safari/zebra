package zebra

type Secret struct {
	secret string
}

func (s *Secret) MarshalText() ([]byte, error) {
	return []byte("*****"), nil
}

func (s *Secret) UnmarshalText(text []byte) error {
	s.secret = string(text)

	return nil
}
