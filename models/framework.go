package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// swagger: model Framework
type Framework struct {
	ID primitive.ObjectID `json:"id" bson:"_id" swwaggerigrnore:"true"`
	// Framework name
	Name string `json:"name" bson:"name"`
	// Description
	Description string `json:"description" bson:"description"`
	// Pentest phases
	Phases []Phase `json:"phases" bson:"phases"`
}

// swagger: Phases
type Phase struct {
	// Phase name
	Name string `json:"name" bson:"name"`
	// Description
	Description string `json:"description" bson:"description"`
	// Information
	Info string `json:"info" bson:"info"`
	// URL
	Url string `json:"url" bson:"url"`
	// Phase common used tools
	Tools []Tool `json:"tools" bson:"tools"`
	// Status of completion
	// example: TODO|DOING|DONE
	Status string `json:"status" bson:"status"`
}

// swagger: Tool
type Tool struct {
	// Tool name
	Name string `json:"name" bson:"name"`
	// Description
	Description string `json:"description" bson:"description"`
	// Synopsis and examples
	Synopsis string `json:"synopsis" bson:"synopsis"`
	// URL for reference
	Url string `json:"url" bson:"url"`
}
