package jwt

import (
	"context"
	"fmt"
	"time"

	"github.com/go-chi/jwtauth/v5"
)

const (
	USER_ID_CLAIM = "user_id"
	EXP_CLAIM     = "exp"
	IAT_CLAIM     = "iat"
)

type Claims struct {
	UserID int64     `json:"user_id"`
	Exp    time.Time `json:"exp"`
	Iat    time.Time `json:"iat"`
}

func NewClaims(claims Claims) map[string]any {
	return map[string]any{
		USER_ID_CLAIM: claims.UserID,
		EXP_CLAIM:     claims.Exp,
		IAT_CLAIM:     claims.Iat,
	}
}

func getClaimInt64(claims map[string]interface{}, key string) (int64, error) {
	value, exists := claims[key]
	if !exists {
		return 0, fmt.Errorf("claim %s not found", key)
	}

	switch v := value.(type) {
	case float64:
		return int64(v), nil
	case int:
		return int64(v), nil
	case int64:
		return v, nil
	case time.Time:
		return v.Unix(), nil
	default:
		return 0, fmt.Errorf("invalid %s type: %T", key, value)
	}
}

func GetClaims(ctx context.Context) (Claims, error) {
	_, claimsMap, err := jwtauth.FromContext(ctx)
	if err != nil {
		return Claims{}, err
	}

	userID, err := getClaimInt64(claimsMap, USER_ID_CLAIM)
	if err != nil {
		return Claims{}, err
	}

	expUnix, err := getClaimInt64(claimsMap, EXP_CLAIM)
	if err != nil {
		return Claims{}, err
	}

	iatUnix, err := getClaimInt64(claimsMap, IAT_CLAIM)
	if err != nil {
		return Claims{}, err
	}

	return Claims{
		UserID: userID,
		Exp:    time.Unix(expUnix, 0),
		Iat:    time.Unix(iatUnix, 0),
	}, nil
}
