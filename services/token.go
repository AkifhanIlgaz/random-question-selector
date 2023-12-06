package services

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/AkifhanIlgaz/random-question-selector/cfg"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt"
	"github.com/thanhpk/randstr"
)

type TokenService struct {
	redisClient *redis.Client
	config      *cfg.Config
	ctx         context.Context
}

func NewTokenService(ctx context.Context, redisClient *redis.Client, config *cfg.Config) *TokenService {
	return &TokenService{
		redisClient: redisClient,
		config:      config,
		ctx:         ctx,
	}
}

func (service *TokenService) GenerateAccessToken(uid string) (string, error) {
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
	claims["sub"] = uid
	claims["exp"] = now.Add(service.config.AccessTokenExpiresIn).Unix()
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(key)
	if err != nil {
		return "", fmt.Errorf("create: sign token: %w", err)
	}

	return token, nil

}

func (service *TokenService) ParseAccessToken(token string) (string, error) {
	decodedPublicKey, err := base64.StdEncoding.DecodeString(service.config.AccessTokenPublicKey)
	if err != nil {
		return "", fmt.Errorf("could not decode: %w", err)
	}

	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return "", fmt.Errorf("unexpected method: %s", t.Header["alg"])
		}

		return jwt.ParseRSAPublicKeyFromPEM(decodedPublicKey)
	})
	if err != nil {
		return "", fmt.Errorf("validate: %w", err)
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return "", fmt.Errorf("validate: invalid token")
	}

	return claims["sub"].(string), nil
}

func (service *TokenService) GenerateRefreshToken(uid string) (string, error) {
	refreshToken := randstr.String(32)

	err := service.redisClient.Set(service.ctx, refreshToken, uid, 0).Err()
	if err != nil {
		return "", fmt.Errorf("generate refresh token: %w", err)
	}

	return refreshToken, nil
}

func (service *TokenService) DeleteRefreshToken(refreshToken string) error {
	if err := service.redisClient.Del(service.ctx, refreshToken).Err(); err != nil {
		return fmt.Errorf("delete refresh token: %w", err)
	}

	return nil
}

func (service *TokenService) GetUidByRefreshToken(refreshToken string) (string, error) {
	uid, err := service.redisClient.Get(service.ctx, refreshToken).Result()
	if err != nil {
		if err == redis.Nil {
			return "", errors.New("refresh token doesn't exist")

		}
		return "", fmt.Errorf("get uid by refresh token: %w", err)
	}

	return uid, nil
}
