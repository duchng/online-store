package jwttoken

import (
	"crypto/x509"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type Config struct {
	privateKey string
	publicKey  string
}

type SignParser interface {
	Scheme() x509.SignatureAlgorithm
	SignClaims(claims jwt.Claims) (string, error)
	ParseClaims(tokenString string) (*jwt.Token, error)
}

type NewOption func(*Config)

func WithPrivateKey(privateKey string) NewOption {
	return func(cfg *Config) {
		cfg.privateKey = privateKey
	}
}

func WithPublicKey(publicKey string) NewOption {
	return func(cfg *Config) {
		cfg.publicKey = publicKey
	}
}

func New(algorithm x509.SignatureAlgorithm, options ...NewOption) (SignParser, error) {
	cfg := &Config{}
	for _, option := range options {
		option(cfg)
	}
	switch algorithm {
	case x509.PureEd25519:
		return NewEd25519(cfg.privateKey, cfg.publicKey)
	default:
		return nil, fmt.Errorf("unsupported algorithm: %v", algorithm)
	}
}
