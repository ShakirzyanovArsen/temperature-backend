package service

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
)

type tokenGenerator interface {
	getToken() (string, error)
}

type tokenGeneratorImpl struct{}

func (tokenGeneratorImpl) getToken() (string, error) {
	b := make([]byte, 256)
	_, e := rand.Read(b)
	if e != nil {
		return "", e
	}
	hash := sha1.New()
	hash.Write(b)
	return hex.EncodeToString(hash.Sum(nil)), nil
}
