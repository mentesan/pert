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
	Info        string `json:"info" bson:"info"`
	Url         string `json:"url" bson:"url"`
	Tools       []Tool `json:"tools" bson:"tools"`
	// TODO|DOING|DONE
	Status string `json:"status" bson:"status"`
}

type Tool struct {
	Name        string `json:"name" bson:"name"`
	Description string `json:"description" bson:"description"`
	Synopsis    string `json:"synopsis" bson:"synopsis"`
	Url         string `json:"url" bson:"url"`
}
