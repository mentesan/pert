package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// swagger:model Project
type Project struct {
	ID primitive.ObjectID `json:"id" bson:"_id" swwaggerigrnore:"true"`
	// Project Name
	Name string `json:"name" bson:"name"`
	// Project type: BlackHat|GrayHat|WhiteHat
	Type string `json:"type" bson:"type"`
	// Company ID
	CompanyId string `json:"companyId" bson:"companyId"`
	// Description
	Description string `json:"description" bson:"description"`
	// Status: Approval|Open|Executing|Reporting|Retesting|Done
	Status string `json:"status" bson:"status"`
	// Targets: URL and/or IP
	Targets []Target `json:"targets" bson:"targets"`
	// Framework to use for the project
	Framework Framework `json:"framework" bson:"framework"`
	// Create Time
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
}

// swagger: Project Targets
type Target struct {
	// Target name
	Name string `json:"name" bson:"name"`
	// URL
	Url string `json:"url" bson:"url"`
	// Ip address
	IP string `json:"ip" bson:"ip"`
	// Description
	Description string `json:"description" bson:"description"`
	// Network/AS
	Network string `json:"network" bson:"network"`
	// Whois Owner
	Owner string `json:"owner" bson:"owner"`
}
