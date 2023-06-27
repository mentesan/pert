package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// swagger: model Vuln
type Vuln struct {
	ID primitive.ObjectID `json:"id" bson:"_id" swaggerignore:"true"`
	// Vulnerability name
	Name string `json:"name" bson:"name"`
	// Risk
	Risk Risk `json:"risk" bson:"risk"`
	// Description
	Description string `json:"description" bson:"description"`
	// swagger: ignore
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	// Status
	Status string `json:"status" bson:"status"`
	// Poc
	Poc Poc `json:"pocs" bson:"pocs"`
}

// swagger: Poc
type Poc struct {
	// Poc name
	Name string `json:"name" bson:"name"`
	// Description
	Description string `json:"description" bson:"description"`
	// Images
	Images []Image `json:"images" bson:"images"`
}

// swagger: Image
type Image struct {
	// Image description/name
	Description string `json:"name" bson:"name"`
	// Path
	Path string `json:"path" bson:"path"`
}

// swagger: Risk
type Risk struct {
	// Risk name
	Name string `json:"name" bson:"name"`
	// Severity
	Severity string `json:"severity" bson:"severity"`
	// Description
	Description string `json:"description" bson:"description"`
	// Color
	Color string `json:"color" bson:"color"`
}
