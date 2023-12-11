package services

import (
	"context"
	"fmt"

	"github.com/AkifhanIlgaz/random-question-selector/models"
	"github.com/AkifhanIlgaz/random-question-selector/utils"
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

	res, err := service.collection.DeleteOne(service.ctx, query)
	if err != nil {
		return fmt.Errorf("delete question: %w", err)
	}

	if res.DeletedCount == 0 {
		return utils.ErrNoQuestion
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
		return nil, fmt.Errorf("no question found with the id: %v", questionId)
	}

	var question models.Question
	res.Decode(&question)

	return &question, nil
}

func (service *QuestionService) UpdateQuestion(questionId string, updatedQuestion models.Question) error {
	query, err := buildIdQuery(questionId)
	if err != nil {
		return fmt.Errorf("update question: %w", err)
	}

	res, err := service.collection.ReplaceOne(service.ctx, query, updatedQuestion)
	if err != nil {
		return fmt.Errorf("update question: %w", err)
	}

	if res.MatchedCount == 0 {
		return fmt.Errorf("no question found with the id: %v", questionId)
	}

	return nil
}

func (service *QuestionService) GetAllQuestionsOfGroup(group string) ([]models.Question, error) {
	query := bson.M{
		"group": group,
	}

	cursor, err := service.collection.Find(service.ctx, query)
	if err != nil {
		return nil, fmt.Errorf("get all questions of group: %w", err)
	}

	var questions []models.Question
	err = cursor.All(service.ctx, &questions)
	if err != nil {
		return nil, fmt.Errorf("get all questions of group: %w", err)
	}

	return questions, nil
}

func (service *QuestionService) GetRandomQuestionsByGroup(group string, count int) ([]models.Question, error) {
	matchStage := bson.D{{Key: "$match", Value: bson.D{{Key: "group", Value: group}}}}
	randomStage := bson.D{{Key: "$sample", Value: bson.D{{Key: "size", Value: count}}}}

	cursor, err := service.collection.Aggregate(service.ctx, mongo.Pipeline{
		matchStage,
		randomStage,
	})
	if err != nil {
		return nil, fmt.Errorf("get random questions by group: %w", err)
	}

	var questions []models.Question
	err = cursor.All(service.ctx, &questions)
	if err != nil {
		return nil, fmt.Errorf("get random questions by group: %w", err)
	}

	return questions, nil
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
