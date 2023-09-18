package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
)

func GetCorrespondPath(ts int64) (string, error) {
	publicKeyPem := []byte(`
-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDLgd2OAkcGVtoE3ThUREbio0Eg
Uc/prcajMKXvkCKFCWhJYJcLkcM2DKKcSeFpD/j6Boy538YXnR6VhcuUJOhH2x71
nzPjfdTcqMz7djHum0qSZA0AyCBDABUqCrfNgCiJ00Ra7GmRj+YCK1NJEuewlb40
JNrRuoEUXpabUzGB8QIDAQAB
-----END PUBLIC KEY-----
`)

	block, _ := pem.Decode(publicKeyPem)
	if block == nil {
		return "", fmt.Errorf("Failed to decode public key")
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", err
	}

	rsaPublicKey, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		return "", fmt.Errorf("Not an RSA public key")
	}

	label := []byte("")
	hash := sha256.New()
	hash.Write([]byte(fmt.Sprintf("refresh_%d", ts)))
	hashed := hash.Sum(nil)

	encrypted, err := rsa.EncryptOAEP(
		hash,
		rand.Reader,
		rsaPublicKey,
		hashed,
		label,
	)

	if err != nil {
		return "", err
	}

	return hex.EncodeToString(encrypted), nil
}
