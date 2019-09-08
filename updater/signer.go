package updater

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/xml"
	"math/big"
)

// ""

type RSAKey struct {
	ModulusString  string   `xml:"Modulus"`
	ExponentString string   `xml:"Exponent"`
	Modulus        *big.Int `xml:"-"`
	Exponent       int      `xml:"-"`
}

// ParsePublicKey parses a string in the form of
// <RSAKeyValue><Modulus>%s</Modulus><Exponent>%s</Exponent></RSAKeyValue>
// returning a struct
func ParsePublicKey(s string) (RSAKey, error) {
	var key RSAKey
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
	e := binary.LittleEndian.Uint32(data)
	key.Exponent = int(e)

	return key, nil
}
