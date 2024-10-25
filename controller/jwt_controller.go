package controller

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"jwt/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateJWT generates a JWT token for the user using RSA private key
func GenerateJWT(user models.User, rsa models.RSAKeyPair) (string, error) {
	privateKeyBlock, _ := pem.Decode([]byte(rsa.PrivateKey))
	if privateKeyBlock == nil {
		return "", errors.New("failed to decode PEM block containing private key")
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
	if err != nil {
		return "", err
	}
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Name,
		"email":    user.Email,
		"roles":    user.Roles,
		"groups":   user.Groups,
		"exp":      time.Now().Add(time.Hour * 72).Unix(), // Token expires in 72 hours
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// ParseJWT parses a JWT token and verifies it using the provided RSA public key.
func ParseJWT(tokenString string, publicKeyPEM string) (jwt.MapClaims, error) {
	// Parse the public key
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(publicKeyPEM))
	if err != nil {
		return nil, err
	}
	// Parse the JWT token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is RS256
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("unexpected signing method")
		}
		// Return the public key for verification
		return publicKey, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

// ValidateJWT validates the JWT token using the provided public key (in PEM format).
func ValidateJWT(tokenString string, publicKey string) (jwt.MapClaims, error) {
	// Decode the public key from PEM format
	block, _ := pem.Decode([]byte(publicKey))
	if block == nil {
		fmt.Println("Failed to decode PEM block: block is nil") // Debugging log
		return nil, errors.New("failed to decode PEM block containing public key: block is nil")
	}
	if block.Type != "RSA PUBLIC KEY" && block.Type != "PUBLIC KEY" {
		fmt.Println("Invalid PEM block type:", block.Type) // Debugging log
		return nil, errors.New("failed to decode PEM block containing public key: invalid block type")
	}
	// Try to parse the key as PKCS#1 (RSA PUBLIC KEY format)
	var rsaPubKey *rsa.PublicKey
	var err error

	rsaPubKey, err = x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		// If PKCS#1 parsing fails, attempt PKIX parsing (default)
		pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			fmt.Println("Error parsing public key:", err) // Debugging log
			return nil, errors.New("failed to parse public key: " + err.Error())
		}
		var ok bool
		rsaPubKey, ok = pubKey.(*rsa.PublicKey)
		if !ok {
			fmt.Println("Failed to cast public key to RSA public key") // Debugging log
			return nil, errors.New("failed to cast public key to RSA public key")
		}
	}
	// Parse the JWT token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method is RSA
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			fmt.Println("Unexpected signing method:", token.Method) // Debugging log
			return nil, errors.New("unexpected signing method")
		}
		return rsaPubKey, nil
	})
	if err != nil {
		fmt.Println("JWT parsing error:", err) // Debugging log
		return nil, errors.New("failed to parse JWT: " + err.Error())
	}
	// Extract the claims from the token
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Check if the token has expired
		if exp, ok := claims["exp"].(float64); ok {
			expirationTime := time.Unix(int64(exp), 0)
			fmt.Println("Token expiration time:", expirationTime) // Debugging log
			if time.Now().After(expirationTime) {
				fmt.Println("Token has expired") // Debugging log
				return nil, errors.New("token has expired")
			}
		} else {
			fmt.Println("exp claim is missing or invalid") // Debugging log
			return nil, errors.New("exp claim is missing or invalid")
		}
		return claims, nil
	}
	fmt.Println("Invalid token claims") // Debugging log
	return nil, errors.New("invalid token claims")
}
