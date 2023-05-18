package token

import (
	"errors"
	"fmt"

	"github.com/o1egl/paseto"
)

// TokenMaker is a PASETO token maker
type TokenMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
	Value        string
}

// errors
var (
	ErrInvalidToken = errors.New("token is invalid")
	KeySize         = int(32)
)

// NewTokenMaker creates a new TokenMaker.
//
// Symmetric key must be 32 character string, often
// stored in an env variable or secret in each service.
func NewTokenMaker(symmetricKey string) (*TokenMaker, error) {
	if len(symmetricKey) != KeySize {
		return nil, fmt.Errorf("invalid key size: must be exactly %d characters", KeySize)
	}
	maker := &TokenMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}
	return maker, nil
}

// CreateToken creates a new token and returns it, with payload p,
// and any errors.
//
// The last token generated can be retreived with Value.
func (maker *TokenMaker) CreateToken(payload map[string]interface{}) (token string, p map[string]interface{}, err error) {
	token, err = maker.paseto.Encrypt(maker.symmetricKey, payload, nil)
	if err == nil {
		maker.Value = token
	}
	return token, payload, err
}

// VerifyToken checks if the token is valid or not.
//
// Will return payload and no error if valid.
func (maker *TokenMaker) VerifyToken(token string) (map[string]interface{}, error) {
	var payload map[string]interface{}

	err := maker.paseto.Decrypt(token, maker.symmetricKey, &payload, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}
	return payload, nil
}
