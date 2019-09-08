package updater

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"testing"

	"src/github.com/stretchr/testify/assert"
)

func TestSigner(t *testing.T) {
	rng := rand.Reader
	privKey, e := rsa.GenerateKey(rng, 2048)
	assert.Nil(t, e)

	// wyUpdate stores the public key as XML with base64 encoded modulus and exponent
	b64Mod := base64.StdEncoding.EncodeToString(privKey.PublicKey.N.Bytes())

	exp := make([]byte, 4)
	binary.LittleEndian.PutUint32(exp, uint32(privKey.PublicKey.E))
	b64Exp := base64.StdEncoding.EncodeToString(exp)

	ks := fmt.Sprintf("<RSAKeyValue><Modulus>%s</Modulus><Exponent>%s</Exponent></RSAKeyValue>", b64Mod, b64Exp)
	key, err := ParsePublicKey(ks)
	assert.Nil(t, err)

	assert.Equal(t, b64Mod, key.ModulusString)
	assert.Equal(t, privKey.PublicKey.N, key.Modulus)
	assert.Equal(t, privKey.PublicKey.E, key.Exponent)
}

func TestSigner_VerifyUpdate(t *testing.T) {
	rng := rand.Reader
	privKey, e := rsa.GenerateKey(rng, 2048)
	assert.Nil(t, e)

	message := []byte("message to be signed")

	hashed := sha1.Sum(message)

	signature, err := rsa.SignPKCS1v15(rng, privKey, crypto.SHA1, hashed[:])
	assert.Nil(t, err)

	// validated
	err = VerifyUpdate(&privKey.PublicKey, hashed[:], signature)
	assert.Nil(t, err)

	// not validated, different signing key
	privKey2, e := rsa.GenerateKey(rng, 2048)
	assert.Nil(t, e)
	err = VerifyUpdate(&privKey2.PublicKey, hashed[:], signature)
	assert.NotNil(t, err)
}
