package token

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt"
	"github.com/xybor/todennus-backend/pkg/xcontext"
)

var _ Engine = (*JWTEngine)(nil)

type JWTEngine struct {
	// RSA Signing method
	rsaPrivateKey *rsa.PrivateKey
	rsaPublicKey  *rsa.PublicKey

	// HMAC Signing method
	hmacSecret []byte
}

func NewJWTEngine() *JWTEngine {
	return &JWTEngine{}
}

func (*JWTEngine) Type() string {
	return "Bearer"
}

func (engine *JWTEngine) WithHMAC(secret string) error {
	if secret == "" {
		return fmt.Errorf("%w: require non-empty hmac secret", ErrSigningKeyInvalid)
	}

	engine.hmacSecret = []byte(secret)
	return nil
}

func (engine *JWTEngine) WithRSA(priv, pub string) error {
	if priv == "" && pub == "" {
		return fmt.Errorf("%w: require non-empty rsa private key or public key", ErrSigningKeyInvalid)
	}

	var err error
	if priv != "" {
		engine.rsaPrivateKey, err = jwt.ParseRSAPrivateKeyFromPEM([]byte(priv))
		if err != nil {
			return err
		}
	}

	if pub != "" {
		engine.rsaPublicKey, err = jwt.ParseRSAPublicKeyFromPEM([]byte(pub))
		if err != nil {
			return err
		}
	}

	return nil
}

func (engine *JWTEngine) Generate(ctx context.Context, claims Claims) (string, error) {
	var token string
	var err error
	switch {
	case engine.rsaPrivateKey != nil:
		token, err = jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(engine.rsaPrivateKey)
	case engine.hmacSecret != nil:
		token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(engine.hmacSecret)
	default:
		err = errors.New("not found any signing method provided for jwt engine")
	}

	return token, err
}

func (engine *JWTEngine) Validate(ctx context.Context, token string, claims Claims) (bool, error) {
	parsedToken, err := jwt.ParseWithClaims(token, claims, engine.publicKeyFunc)
	if err != nil {
		xcontext.Logger(ctx).Debug("failed to parse token", "err", err)
		return false, nil
	}

	_, ok := parsedToken.Claims.(Claims)
	if !ok {
		return false, ErrTokenInvalidFormat
	}

	if !parsedToken.Valid {
		return false, nil
	}

	return true, nil
}

func (engine *JWTEngine) publicKeyFunc(t *jwt.Token) (interface{}, error) {
	switch t.Method.(type) {
	case *jwt.SigningMethodRSA:
		return engine.rsaPublicKey, nil
	case *jwt.SigningMethodHMAC:
		return engine.hmacSecret, nil
	default:
		return nil, ErrTokenSigningMethodNotSupport
	}
}
