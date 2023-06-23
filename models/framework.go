package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Framework struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	Phases      []Phase            `json:"phases" bson:"phases"`
}

type Phase struct {
	Name        string `json:"name" bson:"name"`
	Description string `json:"description" bson:"description"`
	Tools       []Tool `json:"tools" bson:"tools"`
}

type Tool struct {
	Name        string `json:"name" bson:"name"`
	Description string `json:"description" bson:"description"`
	Synopsis    string `json:"synopsis" bson:"synopsis"`
	Url         string `json:"url" bson:"url"`
}
