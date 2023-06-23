package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"pert/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type UsersHandler struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewUsersHandler(ctx context.Context, collection *mongo.Collection) *UsersHandler {
	return &UsersHandler{
		collection: collection,
		ctx:        ctx,
	}
}

func (handler *UsersHandler) ListUsersHandler(c *gin.Context) {
	log.Printf("Request to MongoDB")
	cur, err := handler.collection.Find(handler.ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "err.Error()"})
		return
	}
	defer cur.Close(handler.ctx)

	users := make([]models.User, 0)
	for cur.Next(handler.ctx) {
		var user models.User
		cur.Decode(&user)
		users = append(users, user)
	}
	c.JSON(http.StatusOK, users)
}

func (handler *UsersHandler) NewUserHandler(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return

	}
	user.ID = primitive.NewObjectID()
	_, err := handler.collection.InsertOne(c, user)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while inserting a new user"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (handler *UsersHandler) UpdateUserHandler(c *gin.Context) {
	session := sessions.Default(c)
	// Get session values
	sessionUser := session.Get("username")
	sessionType := session.Get("type")

	if sessionType != "admin" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not authorized"})
		return
	} else {
		// Compare username and usertype of Session with Database values
		cur := handler.collection.FindOne(handler.ctx, bson.M{
			"username": sessionUser,
			"type":     sessionType,
		})
		if cur.Err() != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authorized"})
			return
		}
		// Proceed to update informed values
		id := c.Param("id")
		var user models.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		objectId, _ := primitive.ObjectIDFromHex(id)
		fmt.Println("OBJECT ID:", objectId)

		filter := bson.D{{"_id", objectId}}
		var fieldUpdated int8
		if len(user.FirstName) > 0 {
			log.Println("Entrou aqui")
			update := bson.D{{"$set", bson.D{{"firstName", user.FirstName}}}}
			_, err := handler.collection.UpdateOne(c, filter, update)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			fieldUpdated++
		}
		if len(user.LastName) > 0 {
			log.Println("Entrou aqui")
			update := bson.D{{"$set", bson.D{{"lastName", user.LastName}}}}
			_, err := handler.collection.UpdateOne(c, filter, update)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			fieldUpdated++
		}
		// todo: Check username availability
		if len(user.Username) > 0 {
			log.Println("Entrou aqui")
			update := bson.D{{"$set", bson.D{{"userName", user.Username}}}}
			_, err := handler.collection.UpdateOne(c, filter, update)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			fieldUpdated++
		}
		// todo: Convert password to hash
		if len(user.Password) > 0 {
			log.Println("Entrou aqui")
			// Hash password
			//h := sha256.New()
			hashedPass, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 8)
			update := bson.D{{"$set", bson.D{{"password", string(hashedPass)}}}}
			_, err := handler.collection.UpdateOne(c, filter, update)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			fieldUpdated++
		}
		if len(user.Email) > 0 {
			log.Println("Entrou aqui")
			update := bson.D{{"$set", bson.D{{"email", user.Email}}}}
			_, err := handler.collection.UpdateOne(c, filter, update)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			fieldUpdated++
		}
		// Is UserType informed and valid?
		if len(user.Type) > 0 {
			if user.Type == "admin" || user.Type == "pentester" || user.Type == "client" {
				update := bson.D{{"$set", bson.D{{"type", user.Type}}}}
				log.Println("USER TYPE:", user.Type)
				_, err := handler.collection.UpdateOne(c, filter, update)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				fieldUpdated++
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid"})
				return
			}
		}
		if len(user.CompanyId) > 0 {
			log.Println("Entrou aqui")
			update := bson.D{{"$set", bson.D{{"companyId", user.CompanyId}}}}
			_, err := handler.collection.UpdateOne(c, filter, update)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			fieldUpdated++
		}
		if len(user.FirstName) > 0 {
			log.Println("Entrou aqui")
			update := bson.D{{"$set", bson.D{{"firstName", user.FirstName}}}}
			_, err := handler.collection.UpdateOne(c, filter, update)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			fieldUpdated++
		}
		/*
				update := bson.D{{"$set", bson.D{
					{"firstName", user.FirstName},
					{"lastName", user.LastName},
					{"username", user.Username},
					{"password", user.Password},
					{"email", user.Email},
					{"type", user.Type},
					{"companyId", user.CompanyId},
				}}}
			_, err := handler.collection.UpdateOne(c, filter, update)

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error()})
				return
			}
		*/
		if fieldUpdated > 0 {
			c.JSON(http.StatusOK, gin.H{"message": "User has been updated"})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"message": "User NOT updated"})

		}

	}
}

func (handler *UsersHandler) DeleteUserHandler(c *gin.Context) {
	id := c.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)

	filter := bson.D{{"_id", objectId}}
	_, err := handler.collection.DeleteOne(c, filter)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User has beem deleted"})
}

func (handler *UsersHandler) SearchUserHandler(c *gin.Context) {
	firstName := c.Query("firstName")

	filter := bson.D{{"firstName", bson.D{{"$eq", firstName}}}}
	cur, err := handler.collection.Find(c, filter)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	defer cur.Close(c)

	var users []models.User
	if err = cur.All(c, &users); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Error"})
		return
	}

	c.JSON(http.StatusOK, users)
}
