package handlers

import (
	"context"
	"fmt"
	"net/http"
	"pert-api/models"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gin-contrib/sessions"
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
	// Get session values
	session := sessions.Default(c)
	sessionType := session.Get("type")
	// Session Type
	if sessionType != "admin" || sessionType != "pentester" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not authorized"})
		return
	}
	// Verify if database says the same...
	if userVerified(c, session) == false {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authorized"})
		return
	}

	// Proceed to list
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
	// Get session values
	session := sessions.Default(c)
	sessionType := session.Get("type")
	// Session Type
	if sessionType != "admin" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not authorized"})
		return
	}
	// Verify if database says the same...
	if userVerified(c, session) == false {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authorized"})
		return
	}

	var company models.Company
	if err := c.ShouldBindJSON(&company); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Validate all fields
	var fName, fAddress, fSite bool
	if len(strings.TrimSpace(company.Name)) > 0 && govalidator.IsAlpha(company.Name) {
		fName = true
	}
	if len(strings.TrimSpace(company.Address)) > 0 && govalidator.IsAlpha(company.Address) {
		fAddress = true
	}
	if len(strings.TrimSpace(company.Site)) > 0 && govalidator.IsURL(company.Site) {
		fSite = true
	}
	// If ok, insert company
	if fName && fAddress && fSite {
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
	c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while inserting a new company"})
	return
}

func (handler *CompaniesHandler) UpdateCompanyHandler(c *gin.Context) {
	// Get session values
	session := sessions.Default(c)
	sessionType := session.Get("type")
	// Session Type
	if sessionType != "admin" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not authorized"})
		return
	}
	// Verify if database says the same...
	if userVerified(c, session) == false {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authorized"})
		return
	}

	// Proceed to update
	id := c.Param("id")
	var company models.Company
	if err := c.ShouldBindJSON(&company); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	objectId, _ := primitive.ObjectIDFromHex(id)
	company.UpdatedAt = time.Now()
	filter := bson.D{{"_id", objectId}}
	var fieldUpdated int8

	// Validate all relevant fields
	if len(strings.TrimSpace(company.Name)) > 0 && govalidator.IsAlpha(company.Name) {
		update := bson.D{{"$set", bson.D{{"name", company.Name}}}}
		_, err := handler.collection.UpdateOne(c, filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		fieldUpdated++
	}
	if len(strings.TrimSpace(company.Address)) > 0 && govalidator.IsAlpha(company.Address) {
		update := bson.D{{"$set", bson.D{{"address", company.Address}}}}
		_, err := handler.collection.UpdateOne(c, filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		fieldUpdated++
	}
	if len(strings.TrimSpace(company.Site)) > 0 && govalidator.IsURL(company.Site) {
		update := bson.D{{"$set", bson.D{{"site", company.Site}}}}
		_, err := handler.collection.UpdateOne(c, filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		fieldUpdated++
	}
	if len(strings.TrimSpace(company.Contacts[0].FirstName)) > 0 &&
		govalidator.IsAlpha(company.Contacts[0].FirstName) {

		update := bson.D{{"$set", bson.D{{"contacts", company.Contacts}}}}
		_, err := handler.collection.UpdateOne(c, filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		fieldUpdated++
	}

	if fieldUpdated > 0 {
		c.JSON(http.StatusOK, gin.H{"message": "Company has been updated"})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"error": "Error"})
	return
}

func (handler *CompaniesHandler) DeleteCompanyHandler(c *gin.Context) {
	// Get session values
	session := sessions.Default(c)
	sessionType := session.Get("type")
	// Session Type
	if sessionType != "admin" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not authorized"})
		return
	}
	// Verify if database says the same...
	if userVerified(c, session) == false {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authorized"})
		return
	}

	// Proceed to Delete
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
	// Get session values
	session := sessions.Default(c)
	sessionType := session.Get("type")
	// Session Type
	if sessionType != "admin" || sessionType != "pentester" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not authorized"})
		return
	}
	// Verify if database says the same...
	if userVerified(c, session) == false {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authorized"})
		return
	}

	// Proceed to search
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
