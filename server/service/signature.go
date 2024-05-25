package service

import (
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

func VerifyBentoSignature(hashed, signature, pubKeyBytes []byte) error {
	block, _ := pem.Decode(pubKeyBytes)
	if block == nil || block.Type != "RSA PUBLIC KEY" {
		return errors.New("Invalid public key provided")
	}

	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return err
	}

	rsaPubKey, ok := pubKey.(*rsa.PublicKey)
	if !ok {
		return errors.New("Not an RSA public key")
	}

	err = rsa.VerifyPKCS1v15(rsaPubKey, crypto.SHA256, hashed, signature)
	return err
}
