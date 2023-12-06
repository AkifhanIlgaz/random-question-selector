package services

import (
	"context"

	"github.com/AkifhanIlgaz/random-question-selector/models"
	"go.mongodb.org/mongo-driver/mongo"
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
	return nil
}

func (service *QuestionService) DeleteQuestion(questionId string) error {
	return nil
}

func (service *QuestionService) GetQuestion(questionId string) (models.Question, error) {
	return models.Question{}, nil
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
