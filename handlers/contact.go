package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"pert-api/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ContactsHandler struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewContactsHandler(ctx context.Context, collection *mongo.Collection) *ContactsHandler {

	return &ContactsHandler{
		collection: collection,
		ctx:        ctx,
	}
}

// ListContactsHandler
// @Summary Get contacts list
// @Produce json
// @Success 200 {object} models.Contact
// @Failure 401 {string} Unauthorized
// @Failure 403 {string} Not logged in
// @Failure 500 {string} Cant list contacts
// @Router /contacts [get]
func (handler *ContactsHandler) ListContactsHandler(c *gin.Context) {
	log.Printf("Request to MongoDB")
	cur, err := handler.collection.Find(handler.ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "err.Error()"})
		return
	}
	defer cur.Close(handler.ctx)

	contacts := make([]models.Contact, 0)
	for cur.Next(handler.ctx) {
		var contact models.Contact
		cur.Decode(&contact)
		contacts = append(contacts, contact)
	}
	c.JSON(http.StatusOK, contacts)
}

// NewContactsHandler
// @Summary Add new contact
// @Produce json
// @Param request body models.Contact true "JSON for new contact"
// @Success 200 {string} Contact added
// @Failure 401 {string} Unauthorized
// @Failure 403 {string} Not logged in
// @Failure 500 {string} Cant list users
// @Router /contacts [post]
func (handler *ContactsHandler) NewContactHandler(c *gin.Context) {
	var contact models.Contact
	if err := c.ShouldBindJSON(&contact); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return

	}
	_, err := handler.collection.InsertOne(c, contact)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while inserting a new recipe"})
		return
	}
	log.Println("Remove data from Redis")
	c.JSON(http.StatusOK, contact)
}

// UpdateContactHandler
// @Summary Update contact information
// @Produce json
// @Param id query string true "Contact.ID"
// @Param request body models.Contact true "All fields are optional"
// @Success 200 {string} Contact has been updated
// @Failure 401 {string} Unauthorized
// @Failure 403 {string} Not logged in
// @Router /contacts [put]
func (handler *ContactsHandler) UpdateContactHandler(c *gin.Context) {
	id := c.Param("id")
	var contact models.Contact
	if err := c.ShouldBindJSON(&contact); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	objectId, _ := primitive.ObjectIDFromHex(id)
	fmt.Println("OBJECT ID:", objectId)

	filter := bson.D{{"_id", objectId}}
	update := bson.D{{"$set", bson.D{
		{"firstName", contact.FirstName},
		{"lastName", contact.LastName},
		{"email", contact.Email},
		{"phone", contact.Phone},
		{"info", contact.Info},
	}}}
	_, err := handler.collection.UpdateOne(c, filter, update)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Contact has been updated"})
}

// DeleteContactHandler
// @Summary Delete contact
// @Produce json
// @Param id query string true "Contact.ID"
// @Success 200 {string} Contact has been deleted
// @Failure 401 {string} Not authorized
// @Failure 403 {string} Not logged in
// @Failure 404 {string} Contact not found
// @Router /contacts [delete]
func (handler *ContactsHandler) DeleteContactHandler(c *gin.Context) {
	id := c.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)

	filter := bson.D{{"_id", objectId}}
	_, err := handler.collection.DeleteOne(c, filter)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Contact not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Contact has beem deleted"})
}

// SearchContactHandler
// @Summary Search contact by FirstName
// @Produce json
// @Param request body string true "For search" SchemaExample({ "name": "ContactFirstName" })
// @Success 200 {string} Contact has been deleted
// @Failure 401 {string} Not authorized
// @Failure 403 {string} Not logged in
// @Failure 404 {string} Contact not found
// @Router /contacts/search [get]
func (handler *ContactsHandler) SearchContactHandler(c *gin.Context) {
	firstName := c.Query("firstName")

	filter := bson.D{{"firstName", bson.D{{"$eq", firstName}}}}
	cur, err := handler.collection.Find(c, filter)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Contact not found"})
		return
	}
	defer cur.Close(c)

	var contacts []models.Contact
	if err = cur.All(c, &contacts); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Error"})
		return
	}

	c.JSON(http.StatusOK, contacts)
}
