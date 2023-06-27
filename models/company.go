package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Company model info
// @Description Company
type Company struct {
	ID primitive.ObjectID `json:"id" bson:"_id" swwaggerigrnore:"true"`
	// Company Name
	Name string `json:"name" bson:"name"`
	// Company Address
	Address string `json:"address" bson:"address"`
	// Company Site
	Site string `json:"site" bson:"site"`
	// Company Info
	Info string `json:"info" bson:"info"`
	// Some important details
	Description string `json:"description" bson:"description"`
	//swagger:ignore
	CreatedAt time.Time `json:"publishedAt" bson:"publishedAt"`
	// Date of last update
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
	// Contacts
	Contacts []Contact `json:"contacts" bson:"contacts"`
}

// Contact model info
// @Description Company Contacts
type Contact struct {
	// First name
	FirstName string `json:"firstName" bson:"firstName"`
	// Last name
	LastName string `json:"lastName" bson:"lastName"`
	// Email
	Email string `json:"email" bson:"email"`
	// Phone/cel
	Phone string `json:"phone" bson:"phone"`
	// Info
	Info string `json:"info" bson:"info"`
}
