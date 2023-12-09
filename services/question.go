package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/AkifhanIlgaz/random-question-selector/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type QuestionService struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewQuestionService(ctx context.Context, collection *mongo.Collection) *QuestionService {
	return &QuestionService{
		collection: collection,
		ctx:        ctx,
	}
}

func (service *QuestionService) AddQuestion(question models.Question) error {
	question.Id = primitive.NewObjectID()

	_, err := service.collection.InsertOne(service.ctx, question, options.InsertOne())
	if err != nil {
		return fmt.Errorf("add question: %w", err)
	}

	return nil
}

func (service *QuestionService) DeleteQuestion(questionId string) error {
	query, err := buildIdQuery(questionId)
	if err != nil {
		return fmt.Errorf("delete question: %w", err)
	}

	_, err = service.collection.DeleteOne(service.ctx, query)
	if err != nil {
		return fmt.Errorf("delete question: %w", err)
	}

	return nil
}

func (service *QuestionService) GetQuestion(questionId string) (*models.Question, error) {

	query, err := buildIdQuery(questionId)
	if err != nil {
		return nil, fmt.Errorf("get question: %w", err)
	}

	res := service.collection.FindOne(service.ctx, query)
	if res.Err() == mongo.ErrNoDocuments {
		return nil, errors.New(fmt.Sprintf("No question found with the id: %v", questionId))
	}

	var question models.Question
	res.Decode(&question)

	return &question, nil
}

func (service *QuestionService) EditQuestion(questionId string, updatedQuestion models.Question) error {
	return nil
}

func (service *QuestionService) GetAllQuestionsOfGroup(group string) ([]models.Question, error) {
	return nil, nil
}

func (service *QuestionService) GetRandomQuestionsByGroup(group string, count int) ([]models.Question, error) {
	return nil, nil
}

func buildIdQuery(objectId string) (primitive.M, error) {
	id, err := primitive.ObjectIDFromHex(objectId)
	if err != nil {
		return nil, fmt.Errorf("delete question: %w", err)
	}

	query := bson.M{
		"_id": id,
	}

	return query, nil
}
