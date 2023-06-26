package handlers

import (
	"context"
	"log"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func userVerified(c *gin.Context, session sessions.Session) bool {
	client, err := mongo.Connect(c, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	usersCollection := client.Database(os.Getenv("MONGO_DATABASE")).Collection("users")
	handler := NewUsersHandler(c, usersCollection)

	sessionEmail := session.Get("email")
	sessionType := session.Get("type")
	cur := handler.collection.FindOne(c, bson.M{
		"email": sessionEmail,
		"type":  sessionType,
	})
	if cur.Err() != nil {
		return false
	}
	return true
}
