package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes the password using bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash compares the hashed password with the plain text password
func CheckPasswordHash(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateRSAKeys generates RSA private and public keys with an expiration time
func GenerateRSAKeys() (privateKeyPEM, publicKeyPEM string, expiresAt time.Time, err error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", "", time.Time{}, err
	}

	// Convert private key to PEM format
	privASN1 := x509.MarshalPKCS1PrivateKey(privateKey)
	privBlock := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privASN1,
	}
	privateKeyPEM = string(pem.EncodeToMemory(privBlock))

	// Convert public key to PEM format
	pubASN1 := x509.MarshalPKCS1PublicKey(&privateKey.PublicKey)
	pubBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubASN1,
	}

	publicKeyPEM = string(pem.EncodeToMemory(pubBlock))

	// Set expiration time to 30 days from now (customizable)
	expiresAt = time.Now().Add(30 * 24 * time.Hour)

	return privateKeyPEM, publicKeyPEM, expiresAt, nil
}
