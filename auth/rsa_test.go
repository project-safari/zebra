package auth_test

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"encoding/pem"
	"os"
	"testing"

	"github.com/project-safari/zebra/auth"
	"github.com/stretchr/testify/assert"
)

func TestNewRsaIdentity(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)

	henk, err := auth.Generate()
	assert.Nil(err)

	pkSize := henk.PublicKey().Size()
	assert.Equal(256, pkSize)

	b, e := henk.MarshalText()
	assert.NotNil(b)
	assert.Nil(e)
	assert.Nil(henk.UnmarshalText(b))

	x := auth.Empty()
	b, e = x.MarshalText()
	assert.Nil(b)
	assert.NotNil(e)

	e = x.UnmarshalText([]byte("blah"))
	assert.NotNil(e)

	henkPub := auth.RsaPubIdentity(henk.PublicKey())

	b, e = henkPub.MarshalText()
	assert.Nil(e)
	assert.NotNil(b)
	assert.Nil(x.UnmarshalText(b))

	henkPub2 := henkPub.Public()

	b, e = henkPub2.MarshalText()
	assert.Nil(e)
	assert.NotNil(b)
	assert.Nil(x.UnmarshalText(b))
}

func TestUnmarshalErrors(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	badType := pem.EncodeToMemory(&pem.Block{
		Type:    "NULL KEY",
		Headers: nil,
		Bytes:   nil,
	})
	badPriv := pem.EncodeToMemory(&pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   nil,
	})
	badPub := pem.EncodeToMemory(&pem.Block{
		Type:    "RSA PUBLIC KEY",
		Headers: nil,
		Bytes:   nil,
	})

	henk, err := auth.Generate()
	assert.Nil(err)

	pkSize := henk.PublicKey().Size()
	assert.Equal(256, pkSize)

	x := auth.Empty()
	b, e := x.MarshalText()
	assert.Nil(b)
	assert.NotNil(e)

	henkPub := auth.RsaPubIdentity(henk.PublicKey())

	b, e = henkPub.MarshalText()
	assert.Nil(e)
	assert.NotNil(b)
	assert.NotNil(x.UnmarshalText(badPub))

	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.Nil(err)

	henkPriv := auth.NewRsaIdentity(priv)

	b, e = henkPriv.MarshalText()
	assert.Nil(e)
	assert.NotNil(b)
	assert.NotNil(x.UnmarshalText(badPriv))

	assert.NotNil(x.UnmarshalText(badType))
}

func TestEncrypt(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)

	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.Nil(err)

	henk := auth.NewRsaIdentity(priv)

	priv, err = rsa.GenerateKey(rand.Reader, 2048)
	assert.Nil(err)

	jaap := auth.NewRsaIdentity(priv)

	priv, err = rsa.GenerateKey(rand.Reader, 2048)
	assert.Nil(err)

	ingrid := auth.NewRsaIdentity(priv)

	msg := []byte("Arme mensen kunnen niet met geld omgaan: ze geven alles uit aan eten en kleren, " +
		"terwijl rijke mensen het heel verstandig op de bank zetten.")

	// Lets encrypt it using Ingrid's public key.
	henksMessage, err := henk.Encrypt(msg, ingrid.PublicKey())
	assert.Nil(err)

	jaapsMessage, err := jaap.Encrypt(msg, ingrid.PublicKey())
	assert.Nil(err)

	// Decrypt
	hm, _ := ingrid.Decrypt(henksMessage)
	jm, _ := ingrid.Decrypt(jaapsMessage)

	// Compare the messages of Henk and Jaap, and the original
	assert.True(bytes.Equal(hm, jm))
	assert.True(bytes.Equal(hm, msg))
}

func TestEncryptionNeverTheSame(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)

	// Even when using the same public key, the encrypted messages are never the same
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.Nil(err)

	henk := auth.NewRsaIdentity(priv)

	priv, err = rsa.GenerateKey(rand.Reader, 2048)
	assert.Nil(err)

	jaap := auth.NewRsaIdentity(priv)

	priv, err = rsa.GenerateKey(rand.Reader, 2048)
	assert.Nil(err)

	joop := auth.NewRsaIdentity(priv)

	priv, err = rsa.GenerateKey(rand.Reader, 2048)
	assert.Nil(err)

	koos := auth.NewRsaIdentity(priv)

	priv, err = rsa.GenerateKey(rand.Reader, 2048)
	assert.Nil(err)

	kees := auth.NewRsaIdentity(priv)

	priv, err = rsa.GenerateKey(rand.Reader, 2048)
	assert.Nil(err)

	erik := auth.NewRsaIdentity(priv)

	// Added a couple of the same Identities at the end, just to prove that the
	// encrypted outcome differs each time.
	identities := []*auth.RsaIdentity{henk, jaap, joop, koos, kees, erik, erik, erik, erik}

	priv, err = rsa.GenerateKey(rand.Reader, 2048)
	assert.Nil(err)

	ingrid := auth.NewRsaIdentity(priv)

	msg := []byte("Aan ons land geen polonaise.")
	msgs := make([][]byte, 0, len(identities))

	for _, id := range identities {
		// encrypt the message using Ingrid her public key
		e, _ := id.Encrypt(msg, ingrid.PublicKey())
		msgs = append(msgs, e)
	}

	s := []byte("start")
	for _, m := range msgs {
		assert.False(bytes.Equal(m, s))
	}
}

func TestEncryptDecrypt(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)

	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.Nil(err)

	henk := auth.NewRsaIdentity(priv)

	priv, err = rsa.GenerateKey(rand.Reader, 2048)
	assert.Nil(err)

	ingrid := auth.NewRsaIdentity(priv)

	// a message from Henk to Ingrid
	msg := []byte("Die uitkeringstrekkers pikken al onze banen in.")
	// Lets encrypt it, we want to sent it to Ingrid, thus, we use her public key.
	encryptedMessage, err := henk.Encrypt(msg, ingrid.PublicKey())
	assert.Nil(err)

	// Decrypt Message
	plainTextMessage, err := ingrid.Decrypt(encryptedMessage)
	assert.Nil(err)

	assert.True(bytes.Equal(plainTextMessage, msg))
}

func TestEncryptDecryptMyself(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)

	// If anyone, even you, encrypts (id.e. “locks”) something with your public-key,
	// only you can decrypt it (id.e. “unlock” it) with your secret, private key.
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.Nil(err)

	henk := auth.NewRsaIdentity(priv)

	// a message from Henk
	msg := []byte("Subsidized, dat is toch iets dat je krijgt als je eigenlijk niet goed genoeg bent?")

	// Lets encrypt it, we want to sent it to self, thus, we need our public key.
	encryptedMessage, err := henk.Encrypt(msg, nil)
	assert.Nil(err)

	// Decrypt Message
	plainTextMessage, err := henk.Decrypt(encryptedMessage)
	assert.Nil(err)

	assert.True(bytes.Equal(plainTextMessage, msg))
}

func TestSignVerify(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)

	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.Nil(err)

	henk := auth.NewRsaIdentity(priv)

	// A public message from Henk.
	// note that the message is a byte array, not just a string.
	msg := []byte("Wilders doet tenminste iets tegen de politiek.")

	// Henk signs the message with his private key. This will show the recipient
	// proof that this message is indeed from Henk
	sig, _ := henk.Sign(msg)

	// now, if the message msg is public, anyone can read it.
	// the signature sig however, proves this message is from Henk.
	priv, err = rsa.GenerateKey(rand.Reader, 2048)
	assert.Nil(err)

	ingrid := auth.NewRsaIdentity(priv)

	priv, err = rsa.GenerateKey(rand.Reader, 2048)
	assert.Nil(err)

	hans := auth.NewRsaIdentity(priv)

	err = ingrid.Verify(msg, sig, henk.PublicKey())
	assert.Nil(err)

	err = hans.Verify(msg, sig, henk.PublicKey())
	assert.Nil(err)

	// Let's see if we can break the signature verification
	// (1) changing the message
	err = hans.Verify([]byte("Wilders is een opruier"), sig, henk.PublicKey())
	assert.NotNil(err)

	// (2) changing the signature
	err = hans.Verify(msg, []byte("I am not the signature"), henk.PublicKey())
	assert.NotNil(err)

	// (3) changing the public key
	err = hans.Verify(msg, sig, ingrid.PublicKey())
	assert.NotNil(err)

	_, err = ingrid.Public().Sign([]byte("test"))
	assert.Equal(auth.ErrNoPrivateKey, err)
}

func TestLoad(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)

	id, err := auth.Load("")
	assert.Nil(id)
	assert.NotNil(err)

	id, err = auth.Load("rsa_test.go")
	assert.Nil(id)
	assert.NotNil(err)

	id, err = auth.Load("../simulator/user.key")
	assert.Nil(err)
	assert.NotNil(id)

	assert.Nil(id.Save("new_copy.key"))
	assert.Nil(os.RemoveAll("new_copy.key"))
}
