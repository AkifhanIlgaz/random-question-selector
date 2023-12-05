package services

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/AkifhanIlgaz/random-question-selector/models"
	"github.com/AkifhanIlgaz/random-question-selector/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AuthService struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewAuthService(ctx context.Context, collection *mongo.Collection) *AuthService {
	return &AuthService{
		collection: collection,
		ctx:        ctx,
	}
}

func (service *AuthService) SignUp(user *models.SignUpInput) (*models.User, error) {
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

	// TODO: Create index for user in MongoDB
	opt := options.Index()
	opt.SetUnique(true)
	index := mongo.IndexModel{Keys: bson.M{"email": 1}, Options: opt}

	if _, err := service.collection.Indexes().CreateOne(service.ctx, index); err != nil {
		return nil, errors.New("could not create index for email")
	}

	var newUser *models.User
	query := bson.M{"_id": res.InsertedID}

	err = service.collection.FindOne(service.ctx, query).Decode(&newUser)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}

func (service *AuthService) SignIn(user *models.SignInInput) (*models.User, error) {
	return nil, nil
}
