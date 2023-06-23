package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"pert/models"

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
