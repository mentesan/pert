package models

import (
	"time"
)

type Vuln struct {
	Name        string    `json:"name" bson:"name"`
	Risk        Risk      `json:"risk" bson:"risk"`
	Description string    `json:"description" bson:"description"`
	CreatedAt   time.Time `json:"createdAt" bson:"createdAt"`
	Status      string    `json:"status" bson:"status"`
	Poc         Poc       `json:"pocs" bson:"pocs"`
}

type Poc struct {
	Name        string  `json:"name" bson:"name"`
	Description string  `json:"description" bson:"description"`
	Images      []Image `json:"images" bson:"images"`
}

type Image struct {
	Description string `json:"name" bson:"name"`
	Path        string `json:"path" bson:"path"`
}

type Risk struct {
	Name        string `json:"name" bson:"name"`
	Severity    string `json:"severity" bson:"severity"`
	Description string `json:"description" bson:"description"`
	Color       string `json:"color" bson:"color"`
}
