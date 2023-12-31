package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Question struct {
	Id           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Group        string             `json:"group" bson:"group" binding:"required"`
	Text         string             `json:"text" bson:"text" binding:"required"`
	Answer       string             `json:"answer" bson:"answer" binding:"required"`
	AnswerSource string             `json:"answerSource" bson:"answerSource" binding:"required"`
}
