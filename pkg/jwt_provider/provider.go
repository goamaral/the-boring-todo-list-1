package jwt_provider

import (
	"crypto"
	"fmt"
	"io"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

type Provider struct {
	PrivKey crypto.PrivateKey
	PubKey  crypto.PublicKey
}

func NewProvider(privKeyReader io.Reader, pubKeyReader io.Reader) (Provider, error) {
	privKeyBytes, err := io.ReadAll(privKeyReader)
	if err != nil {
		return Provider{}, fmt.Errorf("failed to read private key: %w", err)
	}
	pubKeyBytes, err := io.ReadAll(pubKeyReader)
	if err != nil {
		return Provider{}, fmt.Errorf("failed to read public key: %w", err)
	}

	privKey, err := jwt.ParseECPrivateKeyFromPEM(privKeyBytes)
	if err != nil {
		return Provider{}, fmt.Errorf("failed to parse private key: %w", err)
	}

	pubKey, err := jwt.ParseECPublicKeyFromPEM(pubKeyBytes)
	if err != nil {
		return Provider{}, fmt.Errorf("failed to parse public key: %w", err)
	}

	return Provider{PrivKey: privKey, PubKey: pubKey}, nil
}

func NewTestProvider(t *testing.T) Provider {
	privKeyBytes := []byte(`-----BEGIN EC PRIVATE KEY-----
MHcCAQEEII4j2z6PFkP5r0XLwEi+UdZXtIqa4batV4JzNfLC2tdsoAoGCCqGSM49
AwEHoUQDQgAEEbzgh4oiwk0ul2KZaNTmgtC684EGNTtPd1awHJqTNQqKtmp1azzh
guMzpI+fC7lZoiyptkjfVkt7a6S3FLcaRA==
-----END EC PRIVATE KEY-----
`)
	pubKeyBytes := []byte(`-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEEbzgh4oiwk0ul2KZaNTmgtC684EG
NTtPd1awHJqTNQqKtmp1azzhguMzpI+fC7lZoiyptkjfVkt7a6S3FLcaRA==
-----END PUBLIC KEY-----
`)

	privKey, err := jwt.ParseECPrivateKeyFromPEM(privKeyBytes)
	require.NoError(t, err)

	pubKey, err := jwt.ParseECPublicKeyFromPEM(pubKeyBytes)
	require.NoError(t, err)

	return Provider{PrivKey: privKey, PubKey: pubKey}
}

func (p Provider) GenerateSignedToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	return token.SignedString(p.PrivKey)
}

func (p Provider) GetClaims(signedToken string) (jwt.Claims, error) {
	token, err := jwt.Parse(signedToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodEd25519); !ok {
			return nil, fmt.Errorf("unexpected signing method (alg: %v)", t.Header["alg"])
		}
		return p.PubKey, nil
	})
	if err != nil {
		return nil, err
	}
	return token.Claims, nil
}
