package updater

import (
	"crypto"
	"crypto/rsa"
	"encoding/base64"
	"encoding/binary"
	"encoding/xml"
	"math/big"
)

// ""

type Key struct {
	ModulusString  string        `xml:"Modulus"`
	ExponentString string        `xml:"Exponent"`
	Modulus        *big.Int      `xml:"-"`
	Exponent       int           `xml:"-"`
	PublicKey      rsa.PublicKey `xml:"-"`
}

// ParsePublicKey parses a string in the form of
// <RSAKeyValue><Modulus>%s</Modulus><Exponent>%s</Exponent></RSAKeyValue>
// returning a struct
func ParsePublicKey(s string) (Key, error) {
	var key Key
	err := xml.Unmarshal([]byte(s), &key)
	if nil != err {
		return key, err
	}

	// convert the base64 modules to a big.Int
	data, err := base64.StdEncoding.DecodeString(key.ModulusString)
	if nil != err {
		return key, err
	}
	z := new(big.Int)
	z.SetBytes(data)
	key.Modulus = z

	// convert the base64 exponent to an int
	data, err = base64.StdEncoding.DecodeString(key.ExponentString)
	if nil != err {
		return key, err
	}

	// sometimes there is not 4 bytes, so we make it 4 bytes
	// >>> binary = base64.b64decode('AQAB')
	// >>> binary
	// '\x01\x00\x01'
	// >>> int(binary.encode('hex'), 16)
	// 65537
	var b [4]byte
	copy(b[4-len(data):], data)
	i := binary.BigEndian.Uint32(b[:])
	key.Exponent = int(i)

	return key, nil
}

func VerifyUpdate(pub *rsa.PublicKey, hashed []byte, sig []byte) error {
	// func VerifyPKCS1v15(pub *PublicKey, hash crypto.Hash, hashed []byte, sig []byte) error
	err := rsa.VerifyPKCS1v15(pub, crypto.SHA1, hashed[:], sig)
	if err != nil {
		return err
	}
	return nil
}
