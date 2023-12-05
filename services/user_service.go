package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/AkifhanIlgaz/random-question-selector/models"
	"github.com/AkifhanIlgaz/random-question-selector/utils"
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

func (service *UserService) CreateUser(user *models.SignUpInput) (*models.User, error) {
	user.CreatedAt = time.Now()
	user.UpdatedAt = user.CreatedAt
	user.Email = strings.ToLower(user.Email)
	user.Verified = true
	user.Role = "user"

	hashedPassword, _ := utils.HashPassword(user.Password)
	user.Password = hashedPassword

	res, err := service.collection.InsertOne(service.ctx, &user)

	if err != nil {
		if er, ok := err.(mongo.WriteException); ok && er.WriteErrors[0].Code == 11000 {
			return nil, errors.New("user with that email already exist")
		}
		return nil, err
	}

	var newUser *models.User
	query := bson.M{"_id": res.InsertedID}

	err = service.collection.FindOne(service.ctx, query).Decode(&newUser)
	if err != nil {
		return nil, err
	}

	return newUser, nil
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

func (service *UserService) UpdateUserById(id string, field string, value string) (*models.User, error) {
	uid, _ := primitive.ObjectIDFromHex(id)
	query := bson.M{
		"_id": uid,
	}
	updateQuery := bson.M{
		"$set": bson.M{
			field: value,
		},
	}

	_, err := service.collection.UpdateOne(service.ctx, query, updateQuery)
	if err != nil {
		return nil, fmt.Errorf("update user by id: %w", err)
	}

	return nil, nil
}

func (service *UserService) UpdateOne(field string, value interface{}) (*models.User, error) {
	query := bson.M{field: value}
	update := bson.M{"$set": bson.M{field: value}}
	result, err := service.collection.UpdateOne(service.ctx, query, update)

	fmt.Print(result.ModifiedCount)
	if err != nil {
		fmt.Print(err)
		return &models.User{}, err
	}

	return &models.User{}, nil
}
