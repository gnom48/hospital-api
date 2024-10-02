package tokens

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
	models "github.com/gnom48/hospital-api-lib"
)

type TokenSigner interface {
	GenerateRegularToken(*models.User) (string, error)
	GenerateCreationToken(*models.User) (string, error)
	ValidateCreationToken(string) (*CreationClaims, error)
	ValidateRegularToken(string) (*RegularClaims, error)
}

type TokenSign struct{}

type RegularClaims struct {
	UserId string `json:"user_id"`
	jwt.StandardClaims
}

type CreationClaims struct {
	UserId   string `json:"user_id"`
	Username string `json:"username"`
	Password string `json:"password"`
	jwt.StandardClaims
}

func (t *TokenSign) GenerateRegularToken(user *models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &RegularClaims{
		UserId: user.Id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Add(24 * 7 * time.Hour)),
			IssuedAt:  jwt.At(time.Now()),
		},
	})
	tokenString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (t *TokenSign) GenerateCreationToken(user *models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &CreationClaims{
		UserId:   user.Id,
		Username: user.Username,
		Password: user.Password,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.At(time.Now()),
		},
	})
	tokenString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (t *TokenSign) ValidateCreationToken(tokenString string) (*CreationClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CreationClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*CreationClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}

func (t *TokenSign) ValidateRegularToken(tokenString string) (*RegularClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&RegularClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(SecretKey), nil
		})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*RegularClaims); ok && token.Valid {
		if claims.ExpiresAt.Unix() < time.Now().Unix() {
			return nil, fmt.Errorf("Token has expired, refresh it")
		}
		return claims, nil
	}

	return nil, err
}
