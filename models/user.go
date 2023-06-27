package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// swagger: model User
type User struct {
	ID primitive.ObjectID `json:"id" bson:"_id" swaggerignore:"true"`
	// User email, used to login
	Email string `json:"email" bson:"email"`
	// Password
	Password string `json:"password" bson:"password"`
	// First name
	FirstName string `json:"firstName" bson:"firstName"`
	// Last name
	LastName string `json:"lastName" bson:"lastName"`
	// Type: admin|pentester|client
	Type string `json:"type" bson:"type"`
	// If it's a "client" type, it needs to have a company
	CompanyId string `json:"companyId" bson:"companyId"`
}
