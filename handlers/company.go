package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"pert/models"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CompaniesHandler struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewCompaniesHandler(ctx context.Context, collection *mongo.Collection) *CompaniesHandler {

	return &CompaniesHandler{
		collection: collection,
		ctx:        ctx,
	}
}

func (handler *CompaniesHandler) ListCompaniesHandler(c *gin.Context) {
	log.Printf("Request to MongoDB")
	cur, err := handler.collection.Find(handler.ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "err.Error()"})
		return
	}
	defer cur.Close(handler.ctx)

	companies := make([]models.Company, 0)
	for cur.Next(handler.ctx) {
		var company models.Company
		cur.Decode(&company)
		companies = append(companies, company)
	}
	c.JSON(http.StatusOK, companies)
}

func (handler *CompaniesHandler) NewCompanyHandler(c *gin.Context) {
	var company models.Company
	if err := c.ShouldBindJSON(&company); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return

	}
	company.ID = primitive.NewObjectID()
	company.CreatedAt = time.Now()
	_, err := handler.collection.InsertOne(c, company)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while inserting a new company"})
		return
	}
	c.JSON(http.StatusOK, company)
}

func (handler *CompaniesHandler) UpdateCompanyHandler(c *gin.Context) {
	id := c.Param("id")
	var company models.Company
	if err := c.ShouldBindJSON(&company); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	objectId, _ := primitive.ObjectIDFromHex(id)
	fmt.Println("OBJECT ID:", objectId)

	company.UpdatedAt = time.Now()
	filter := bson.D{{"_id", objectId}}
	update := bson.D{{"$set", bson.D{
		{"name", company.Name},
		{"address", company.Address},
		{"site", company.Site},
		{"info", company.Info},
		{"description", company.Description},
		{"updatedAt", company.UpdatedAt},
		{"contacts", company.Contacts},
	}}}
	_, err := handler.collection.UpdateOne(c, filter, update)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Company has been updated"})
}

func (handler *CompaniesHandler) DeleteCompanyHandler(c *gin.Context) {
	id := c.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)

	filter := bson.D{{"_id", objectId}}
	_, err := handler.collection.DeleteOne(c, filter)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Company not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Company has beem deleted"})
}

func (handler *CompaniesHandler) SearchCompanyHandler(c *gin.Context) {
	firstName := c.Query("firstName")

	filter := bson.D{{"firstName", bson.D{{"$eq", firstName}}}}
	cur, err := handler.collection.Find(c, filter)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Company not found"})
		return
	}
	defer cur.Close(c)

	var companies []models.Company
	if err = cur.All(c, &companies); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Error"})
		return
	}

	c.JSON(http.StatusOK, companies)
}
