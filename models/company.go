package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// swagger:model Company
type Contact struct {
	FirstName string `json:"firstName" bson:"firstName"`
	LastName  string `json:"lastName" bson:"lastName"`
	Email     string `json:"email" bson:"email"`
	Phone     string `json:"phone" bson:"phone"`
	Info      string `json:"info" bson:"info"`
}

type Company struct {
	//swagger:ignore
	ID primitive.ObjectID `json:"id" bson:"_id"`
	// Company Name
	// example: White Enterprises
	Name string `json:"name" bson:"name"`
	// Company Address
	// example: Rua Jorca, 43 - Bauru/SP
	Address string `json:"address" bson:"address"`
	// Company Site
	// example: https://wwww.whiteenter.net
	Site string `json:"site" bson:"site"`
	// Company Info
	// example: Weekdays from 9 to 17hs
	Info string `json:"info" bson:"info"`
	// Some important details
	// example: A very imporant detail
	Description string `json:"description" bson:"description"`
	// Date of creation
	// example: 2023-06-26T00:53:16.535668Z
	CreatedAt time.Time `json:"publishedAt" bson:"publishedAt"`
	// Date of last update
	// example: 2023-06-28T02:53:16.535668Z
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
	// Company contacts
	Contacts []Contact `json:"contacts" bson:"contacts"`
}
