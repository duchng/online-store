package jwttoken

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/base64"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type schemeEd25519 struct {
	publicKey  ed25519.PublicKey
	privateKey ed25519.PrivateKey
}

func (e *schemeEd25519) Scheme() x509.SignatureAlgorithm {
	return x509.PureEd25519
}

func (e *schemeEd25519) SignClaims(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
	tokenString, err := token.SignedString(e.privateKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (e *schemeEd25519) ParseClaims(tokenString string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(
		tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodEd25519); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return e.publicKey, nil
		},
		jwt.WithValidMethods([]string{"EdDSA"}),
	)
}

func parsePrivateKeyEdDSA(base64Key string) (ed25519.PrivateKey, error) {
	key, err := base64.StdEncoding.DecodeString(base64Key)
	if err != nil {
		return nil, err
	}

	return key, nil
}

func parsePublicKeyEdDSA(base64Key string) (ed25519.PublicKey, error) {
	key, err := base64.StdEncoding.DecodeString(base64Key)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func NewEd25519(base64PrivateKey, base64PublicKey string) (SignParser, error) {
	scheme := &schemeEd25519{}
	if base64PrivateKey != "" {
		privateKey, err := parsePrivateKeyEdDSA(base64PrivateKey)
		if err != nil {
			return nil, err
		}
		scheme.privateKey = privateKey
	}
	if base64PublicKey != "" {
		publicKey, err := parsePublicKeyEdDSA(base64PublicKey)
		if err != nil {
			return nil, err
		}
		scheme.publicKey = publicKey
	}
	return scheme, nil
}
