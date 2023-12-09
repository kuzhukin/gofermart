package cryptographer

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
)

const cryptoKey = "7WXJz8RIFqv8zBsL67vtVw=="
const cryptoVec = "WCwiTM0s8JuF7eFs"

type Key string

type Cryptographer interface {
	Encrypt(data string) (string, error)
}

type AesCryptographer struct {
	key []byte
	vec []byte
}

func NewAesCryptographer() (*AesCryptographer, error) {
	key, err := base64.StdEncoding.DecodeString(cryptoKey)
	if err != nil {
		return nil, fmt.Errorf("convert crypto key from base64, err=%w", err)
	}

	vec, err := base64.StdEncoding.DecodeString(cryptoVec)
	if err != nil {
		return nil, fmt.Errorf("convert crypto vec from base64, err=%w", err)
	}

	return &AesCryptographer{
		key: key,
		vec: vec,
	}, nil
}

func (c *AesCryptographer) Encrypt(data string) (string, error) {
	aesblock, err := aes.NewCipher(c.key)
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return "", err
	}

	aesData := aesgcm.Seal(nil, c.vec, []byte(data), nil)
	encrypted := base64.StdEncoding.EncodeToString(aesData)

	return encrypted, nil
}
