package encrypt

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

var ErrInvalidBlock = errors.New("invalid block")

func Key(privateKey []byte) (key *rsa.PrivateKey, err error) {
	block, _ := pem.Decode(privateKey)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		err = ErrInvalidBlock
		return
	}
	if key, err = x509.ParsePKCS1PrivateKey(block.Bytes); err != nil {
		return
	}
	return
}

func Cert(publicCrt []byte) (cert *x509.Certificate, err error) {
	block, _ := pem.Decode(publicCrt)
	if block == nil || block.Type != "CERTIFICATE" {
		err = ErrInvalidBlock
		return
	}
	if cert, err = x509.ParseCertificate(block.Bytes); err != nil {
		return
	}
	return
}

func KeyCert(privateKey, publicCrt []byte) (confidential Confidential, err error) {
	confidential.Key, err = Key(privateKey)
	if err != nil {
		return
	}
	confidential.Cert, err = Cert(publicCrt)
	return
}

type Confidential struct {
	Key  *rsa.PrivateKey
	Cert *x509.Certificate
}
