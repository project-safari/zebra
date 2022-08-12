package auth

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
)

const ReadOnly = 0o600

var (
	ErrIdentityEmpty   = errors.New("identity is empty")
	ErrBadPEMFile      = errors.New("bad PEM file")
	ErrUnknownPEMBlock = errors.New("unknown PEM block")
	ErrNoPrivateKey    = errors.New("no private key")
)

const RSAKeySize = 2048

// RsaIdentity is just a small struct that clearly differentiates between the
// private and public key of an RSA keypair.
type RsaIdentity struct {
	public  *rsa.PublicKey
	private *rsa.PrivateKey
}

// function for generation of rsa key, returns the rsa and a(n) (potential) error.

func Generate() (*RsaIdentity, error) {
	priv, err := rsa.GenerateKey(rand.Reader, RSAKeySize)

	return &RsaIdentity{
		private: priv,
		public:  &priv.PublicKey,
	}, err
}

// function to empty rsa.

func Empty() *RsaIdentity {
	return &RsaIdentity{
		private: nil,
		public:  nil,
	}
}

// function to load rsa from file, returns the rsa and a(n) (potential) error.

func Load(rsaFile string) (*RsaIdentity, error) {
	rsaText, err := ioutil.ReadFile(rsaFile)
	if err != nil {
		return nil, err
	}

	id := Empty()
	if err := id.UnmarshalText(rsaText); err != nil {
		return nil, err
	}

	return id, nil
}

// function to save rsa to file, returns an error.

func (r *RsaIdentity) Save(rsaFile string) error {
	data, err := r.MarshalText()
	if err != nil {
		return err
	}

	return ioutil.WriteFile(rsaFile, data, ReadOnly)
}

func (r *RsaIdentity) MarshalText() ([]byte, error) {
	if r.private != nil {
		return pem.EncodeToMemory(&pem.Block{
			Type:    "RSA PRIVATE KEY",
			Headers: nil,
			Bytes:   x509.MarshalPKCS1PrivateKey(r.private),
		}), nil
	} else if r.public != nil {
		return pem.EncodeToMemory(&pem.Block{
			Type:    "RSA PUBLIC KEY",
			Headers: nil,
			Bytes:   x509.MarshalPKCS1PublicKey(r.public),
		}), nil
	}

	return nil, ErrIdentityEmpty
}

func (r *RsaIdentity) UnmarshalText(text []byte) error {
	b, _ := pem.Decode(text)
	if b == nil {
		return ErrBadPEMFile
	}

	if b.Type == "RSA PRIVATE KEY" {
		p, e := x509.ParsePKCS1PrivateKey(b.Bytes)
		if e != nil {
			return e
		}

		r.private = p
		r.public = &p.PublicKey

		return nil
	} else if b.Type == "RSA PUBLIC KEY" {
		p, e := x509.ParsePKCS1PublicKey(b.Bytes)
		if e != nil {
			return e
		}

		r.public = p

		return nil
	}

	return ErrUnknownPEMBlock
}

// function to make rsa to string.

func (r *RsaIdentity) String() string {
	b, e := r.MarshalText()
	if e != nil {
		return ""
	}

	return string(b)
}

// NewRsaIdentity returns a new identity with spefied keys.
func NewRsaIdentity(pri *rsa.PrivateKey) *RsaIdentity {
	return &RsaIdentity{
		private: pri,
		public:  &pri.PublicKey,
	}
}

// RsaPubIdentity returns identity with public key. This identity object can
// only be used to verify messages.
func RsaPubIdentity(pub *rsa.PublicKey) *RsaIdentity {
	return &RsaIdentity{
		private: nil,
		public:  pub,
	}
}

// function for public rsa, returns the public key of rsa.

func (r *RsaIdentity) PublicKey() *rsa.PublicKey {
	return r.public
}

func (r *RsaIdentity) Public() *RsaIdentity {
	return &RsaIdentity{
		public:  r.public,
		private: nil,
	}
}

// Sign returns a signature made by combining the message and the signers private key
// With the r.Verify function, the signature can be checked.
func (r *RsaIdentity) Sign(msg []byte) ([]byte, error) {
	hs := r.getHashSum(msg)

	if r.private == nil {
		return nil, ErrNoPrivateKey
	}

	return rsa.SignPKCS1v15(rand.Reader, r.private, crypto.SHA256, hs)
}

// Verify checks if a message is signed by a given Public Key.
func (r *RsaIdentity) Verify(msg []byte, sig []byte, pubKey *rsa.PublicKey) error {
	hs := r.getHashSum(msg)

	if pubKey == nil {
		pubKey = r.PublicKey()
	}

	return rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, hs, sig)
}

// Encrypt's the message using EncryptOAEP which encrypts the given message with RSA-OAEP.
// https://en.wikipedia.org/wiki/Optimal_asymmetric_encryption_padding
// Returns the encrypted message and an error.
func (r *RsaIdentity) Encrypt(msg []byte, key *rsa.PublicKey) ([]byte, error) {
	label := []byte("")
	hash := sha256.New()

	if key == nil {
		key = r.PublicKey()
	}

	return rsa.EncryptOAEP(hash, rand.Reader, key, msg, label)
}

// Decrypt a message using your private key.
// A received message should be encrypted using the receivers public key.
func (r *RsaIdentity) Decrypt(msg []byte) ([]byte, error) {
	label := []byte("")
	hash := sha256.New()

	return rsa.DecryptOAEP(hash, rand.Reader, r.private, msg, label)
}

func (r *RsaIdentity) getHashSum(msg []byte) []byte {
	h := sha256.New()
	h.Write(msg)

	return h.Sum(nil)
}
