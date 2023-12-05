package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/AkifhanIlgaz/random-question-selector/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserService struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewUserService(ctx context.Context, collection *mongo.Collection) *UserService {
	return &UserService{collection, ctx}
}

func (service *UserService) FindUserById(id string) (*models.User, error) {
	uid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("find user by id: %w", err)
	}

	var user models.User

	query := bson.M{
		"_id": uid,
	}
	err = service.collection.FindOne(service.ctx, query).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("find user by id: %w", err)
		}
		return nil, fmt.Errorf("find user by id: %w", err)
	}

	return &user, nil
}

func (service *UserService) FindUserByEmail(email string) (*models.User, error) {
	var user models.User

	query := bson.M{
		"email": strings.ToLower(email),
	}
	err := service.collection.FindOne(service.ctx, query).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("find user by email: %w", err)
		}
		return nil, fmt.Errorf("find user by email: %w", err)
	}

	return &user, nil
}
