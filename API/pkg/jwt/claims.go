package jwt

import (
	"context"
	"errors"
	"time"

	"github.com/go-chi/jwtauth/v5"
)

const (
	USER_ID_CLAIM = "user_id"
	EXP_CLAIM     = "exp"
	IAT_CLAIM     = "iat"
)

type Claims struct {
	UserID int64 `json:"user_id"`
	Exp    int64 `json:"exp"`
	Iat    int64 `json:"iat"`
}

func NewClaims(claims Claims) map[string]any {
	now := time.Now()
	return map[string]any{
		USER_ID_CLAIM: claims.UserID,
		EXP_CLAIM:     now.Add(24 * time.Hour).Unix(),
		IAT_CLAIM:     now.Unix(),
	}
}

func GetClaims(ctx context.Context) (Claims, error) {
	_, claims, err := jwtauth.FromContext(ctx)
	if err != nil {
		return Claims{}, err
	}

	userID, ok := claims[USER_ID_CLAIM].(float64)
	if !ok {
		return Claims{}, errors.New("invalid user_id in claims")
	}

	exp, ok := claims[EXP_CLAIM].(time.Time)
	if !ok {
		return Claims{}, errors.New("invalid exp in claims")
	}

	iat, ok := claims[IAT_CLAIM].(time.Time)
	if !ok {
		return Claims{}, errors.New("invalid iat in claims")
	}

	return Claims{
		UserID: int64(userID),
		Exp:    exp.Unix(),
		Iat:    iat.Unix(),
	}, nil
}
