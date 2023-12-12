package services

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/AkifhanIlgaz/random-question-selector/cfg"
	"github.com/AkifhanIlgaz/random-question-selector/utils"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt"
	"github.com/thanhpk/randstr"
)

type TokenService struct {
	redisClient *redis.Client
	config      *cfg.Config
	ctx         context.Context
}

type RedisClaims struct {
	Subject string   `json:"subject"`
	Groups  []string `json:"groups"`
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

	claims := jwt.StandardClaims{}
	claims.Subject = uid
	claims.ExpiresAt = now.Add(service.config.AccessTokenExpiresIn).Unix()
	claims.IssuedAt = now.Unix()
	claims.NotBefore = now.Unix()

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(key)
	if err != nil {
		return "", fmt.Errorf("create: sign token: %w", err)
	}

	return token, nil
}

func (service *TokenService) ParseAccessToken(token string) (*jwt.StandardClaims, error) {
	decodedPublicKey, err := base64.StdEncoding.DecodeString(service.config.AccessTokenPublicKey)
	if err != nil {
		return nil, fmt.Errorf("could not decode: %w", err)
	}

	parsedToken, err := jwt.ParseWithClaims(token, &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return jwt.StandardClaims{}, fmt.Errorf("unexpected method: %s", t.Header["alg"])
		}

		return jwt.ParseRSAPublicKeyFromPEM(decodedPublicKey)
	})
	if err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}

	claims, ok := parsedToken.Claims.(*jwt.StandardClaims)
	if !ok || !parsedToken.Valid {
		return nil, fmt.Errorf("validate: invalid token")
	}

	return claims, nil
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

func (service *TokenService) GetSub(refreshToken string) (string, error) {
	sub, err := service.redisClient.Get(service.ctx, refreshToken).Result()
	if err != nil {
		if err == redis.Nil {
			return "", utils.ErrNoRefreshToken
		}
		return "", utils.ErrSomethingWentWrong
	}

	return sub, nil
}
