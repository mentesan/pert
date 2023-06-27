package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Project struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Name        string             `json:"name" bson:"name"`
	Type        string             `json:"type" bson:"type"`
	CompanyId   string             `json:"companyId" bson:"companyId"`
	Description string             `json:"description" bson:"description"`
	Status      string             `json:"status" bson:"status"`
	Targets     []Target           `json:"targets" bson:"targets"`
	Framework   Framework          `json:"framework" bson:"framework"`
	CreatedAt   time.Time          `json:"createdAt" bson:"createdAt"`
}

type Target struct {
	Name        string `json:"name" bson:"name"`
	Url         string `json:"url" bson:"url"`
	IP          string `json:"ip" bson:"ip"`
	Description string `json:"description" bson:"description"`
	Network     string `json:"network" bson:"network"`
	Owner       string `json:"owner" bson:"owner"`
}
