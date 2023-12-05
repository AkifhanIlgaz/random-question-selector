package services

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/AkifhanIlgaz/random-question-selector/cfg"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt"
)

type TokenService struct {
	redisClient *redis.Client
	config      *cfg.Config
}

func NewTokenService(redisClient *redis.Client, config *cfg.Config) *TokenService {
	return &TokenService{
		redisClient: redisClient,
		config:      config,
	}
}

func (service *TokenService) GenerateAccessToken(payload any) (string, error) {
	decodedPrivateKey, err := base64.StdEncoding.DecodeString(service.config.AccessTokenPrivateKey)
	if err != nil {
		return "", fmt.Errorf("could not decode key: %w", err)
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(decodedPrivateKey)
	if err != nil {
		return "", fmt.Errorf("create: parse key: %w", err)
	}

	now := time.Now().UTC()

	claims := jwt.MapClaims{}
	claims["sub"] = payload
	claims["exp"] = now.Add(service.config.AccessTokenExpiresIn).Unix()
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(key)
	if err != nil {
		return "", fmt.Errorf("create: sign token: %w", err)
	}

	return token, nil

}

func (service *TokenService) ParseAccessToken(token string) (any, error) {
	decodedPublicKey, err := base64.StdEncoding.DecodeString(service.config.AccessTokenPublicKey)
	if err != nil {
		return nil, fmt.Errorf("could not decode: %w", err)
	}

	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected method: %s", t.Header["alg"])
		}

		return jwt.ParseRSAPublicKeyFromPEM(decodedPublicKey)
	})
	if err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return nil, fmt.Errorf("validate: invalid token")
	}

	return claims["sub"], nil
}

func (service *TokenService) GenerateRefreshToken(uid string) (string, error) {
	return "", nil
}
