package updater

import (
	"crypto/rand"
	"crypto/rsa"
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
