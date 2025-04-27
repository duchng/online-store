package jwttoken

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestIssueJWT(t *testing.T) {
	privateKey := "E8fO3LakFD4FjD+hp0hX3bXRDKdMhwp9Dlr/mgatAFHrM781NjUcv8ZzKv4pVrnEs9UqZd27GoM96mWwyjSLsw=="
	publicKey := "6zO/NTY1HL/Gcyr+KVa5xLPVKmXduxqDPeplsMo0i7M="
	signParser, _ := New(
		x509.PureEd25519, WithPrivateKey(privateKey), WithPublicKey(publicKey),
	)

	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "duc",
			Subject:   "ai_do@duchng.com.vn",
			Audience:  []string{"*.duchng.com.vn"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 4)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		Roles: []string{"admin"},
	}

	signedClaims, err := signParser.SignClaims(claims)
	if err != nil {
		t.Errorf("signParser.SignClaims() error = %v", err)
		return
	}

	t.Logf("signedClaims: %+v", signedClaims)

	cl, err := signParser.ParseClaims(signedClaims)
	if err != nil {
		t.Errorf("signParser.ParseClaims() error = %v", err)
		return
	}

	t.Logf("claims: %+v", cl.Claims.(*Claims))
}

func TestGenerateKey(t *testing.T) {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}

	// Encode the public key
	publicKeyBase64 := base64.StdEncoding.EncodeToString(publicKey)
	fmt.Println(publicKeyBase64)

	// Encode the private key
	privateKeyBase64 := base64.StdEncoding.EncodeToString(privateKey)
	fmt.Println(privateKeyBase64)
}
