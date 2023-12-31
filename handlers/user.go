package handlers

import (
	"context"
	"log"
	"net/http"
	"os"
	"pert-api/models"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
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

// ListUsersHandler
// @Summary Get users list
// @Produce json
// @Success	200 {object} models.User
// @Failure 401 {string} Unauthorized
// @Failure 403 {string} Not logged in
// @Failure 500 {string} Cant list users
// @Router /users [get]
func (handler *UsersHandler) ListUsersHandler(c *gin.Context) {
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
	// Proceed to list users
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

// NewUserHandler
// @Summary Add a new user
// @Accept json
// @Produce json
// @Param request body models.User true "JSON for new user"
// @Success	200 {string} User added
// @Failure 400 {string} BadRequest
// @Failure 401 {string} Unauthorized
// @Failure 403 {string} Not logged in
// @Failure 500 {string} InternalServerError Bad Company
// @Router /users [post]
func (handler *UsersHandler) NewUserHandler(c *gin.Context) {
	// Get session values
	session := sessions.Default(c)
	sessionType := session.Get("type")

	if sessionType != "admin" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not authorized"})
		return
	}
	if userVerified(c, session) == false {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authorized"})
		return
	}
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate all fields
	var hashedPass []byte
	var fEmail, fPassword, fFirstName, fLastName, fType, fCompanyId bool
	if len(strings.TrimSpace(user.Email)) > 0 && govalidator.IsEmail(user.Email) {
		fEmail = true
	}
	if len(strings.TrimSpace(user.FirstName)) > 0 && govalidator.IsAlpha(user.FirstName) {
		fFirstName = true
	}
	if len(strings.TrimSpace(user.LastName)) > 0 && govalidator.IsAlpha(user.LastName) {
		fLastName = true
	}
	if len(strings.TrimSpace(user.Password)) > 0 {
		// Hash password
		hashedPass, _ = bcrypt.GenerateFromPassword([]byte(user.Password), 8)
		user.Password = string(hashedPass)
		fPassword = true
	}
	if len(strings.TrimSpace(user.Type)) > 0 && govalidator.IsAlpha(user.Type) {
		if user.Type == "admin" || user.Type == "pentester" || user.Type == "client" {
			fType = true
		}
	}
	// CompanyId
	if user.Type == "client" {
		if len(strings.TrimSpace(user.CompanyId)) > 0 {
			// Check if company really exists
			client, err := mongo.Connect(c, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
			if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
				log.Fatal(err)
			}
			companiesCollection := client.Database(os.Getenv("MONGO_DATABASE")).Collection("companies")
			companiesHandler := NewCompaniesHandler(c, companiesCollection)
			companyId, _ := primitive.ObjectIDFromHex(user.CompanyId)
			companyCursor := companiesHandler.collection.FindOne(c, bson.M{
				"_id": companyId,
			})
			// If query OK
			if companyCursor.Err() == nil {
				fCompanyId = true
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid Company"})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Company must be selected"})
			return

		}
	}
	// If all fields OK, proceed to insert
	if fEmail && fPassword && fFirstName && fLastName && fType && fCompanyId {
		user.ID = primitive.NewObjectID()
		_, err := handler.collection.InsertOne(c, user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while inserting a new user"})
			return
		}
		c.JSON(http.StatusOK, user)
	}
	c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while inserting a new user"})
	return
}

// UpdateUserHandler
// @Summary Update a user
// @Accept json
// @Produce json
// @Param request body models.User true "All fields are optional"
// @Param  id query string true "User.ID"
// @Success	200 {string} User updated
// @Failure 400 {string} BadRequest
// @Failure 401 {string} Unauthorized
// @Failure 403 {string} Not logged in
// @Failure 500 {string} InternalServerError Data malformed
// @Router /users/:id [put]
func (handler *UsersHandler) UpdateUserHandler(c *gin.Context) {
	// Get session values
	session := sessions.Default(c)
	sessionType := session.Get("type")

	if sessionType != "admin" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not authorized"})
		return
	}
	if userVerified(c, session) == false {
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
	filter := bson.D{{"_id", objectId}}
	var fieldUpdated int8

	if len(strings.TrimSpace(user.Email)) > 0 && govalidator.IsEmail(user.Email) {
		update := bson.D{{"$set", bson.D{{"email", user.Email}}}}
		_, err := handler.collection.UpdateOne(c, filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		fieldUpdated++
	}
	if len(strings.TrimSpace(user.FirstName)) > 0 && govalidator.IsAlpha(user.FirstName) {
		update := bson.D{{"$set", bson.D{{"firstName", user.FirstName}}}}
		_, err := handler.collection.UpdateOne(c, filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		fieldUpdated++
	}
	if len(strings.TrimSpace(user.LastName)) > 0 && govalidator.IsAlpha(user.LastName) {
		update := bson.D{{"$set", bson.D{{"lastName", user.LastName}}}}
		_, err := handler.collection.UpdateOne(c, filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		fieldUpdated++
	}
	if len(strings.TrimSpace(user.Password)) > 0 {
		// Hash password
		hashedPass, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 8)
		update := bson.D{{"$set", bson.D{{"password", string(hashedPass)}}}}
		_, err := handler.collection.UpdateOne(c, filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		fieldUpdated++
	}
	// Is UserType informed and valid?
	if len(strings.TrimSpace(user.Type)) > 0 && govalidator.IsAlpha(user.Type) {
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

	// CompanyId
	if user.Type == "client" {
		if len(strings.TrimSpace(user.CompanyId)) > 0 {
			// Check if company really exists
			client, err := mongo.Connect(c, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
			if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
				log.Fatal(err)
			}
			log.Println("Connected to MongoDB")
			companiesCollection := client.Database(os.Getenv("MONGO_DATABASE")).Collection("companies")
			companiesHandler := NewCompaniesHandler(c, companiesCollection)
			companyId, _ := primitive.ObjectIDFromHex(user.CompanyId)
			companyCursor := companiesHandler.collection.FindOne(c, bson.M{
				"_id": companyId,
			})
			if companyCursor.Err() != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Company"})
				return
			}

			// If Company exists, proceed to update
			update := bson.D{{"$set", bson.D{{"companyId", user.CompanyId}}}}
			_, err = handler.collection.UpdateOne(c, filter, update)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			fieldUpdated++
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Company must be selected"})
			return
		}
	}

	// If one or more field updated
	if fieldUpdated > 0 {
		c.JSON(http.StatusOK, gin.H{"message": "User has been updated"})
		return
	}
	c.JSON(http.StatusBadRequest, gin.H{"message": "User NOT updated"})
	return
}

// DeleteUserHandler
// @Summary Delete a user
// @Accept json
// @Produce json
// @Param  id query string true "User.ID"
// @Success	200 {string} User deleted
// @Failure 401 {string} Unauthorized
// @Failure 403 {string} Not logged in
// @Failure 404 {string} User not found
// @Router /users/:id [delete]
func (handler *UsersHandler) DeleteUserHandler(c *gin.Context) {
	// Get session values
	session := sessions.Default(c)
	sessionType := session.Get("type")

	if sessionType != "admin" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authorized"})
		return
	} else {
		// Compare email and usertype of Session with Database values
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
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "User has beem deleted"})
	}
}

// SearchUserHandler
// @Summary Search user by FirstName
// @Accept json
// @Produce json
// @Param request body string true "For search" SchemaExample({ "firstName": "UserFirstName" })
// @Success	200 {object} models.User
// @Failure 401 {string} Unauthorized
// @Failure 403 {string} Not logged in
// @Failure 404 {string} User not found
// @Router /users/search/ [get]
func (handler *UsersHandler) SearchUserHandler(c *gin.Context) {
	// Get session values
	session := sessions.Default(c)
	sessionType := session.Get("type")
	// Session Type
	if sessionType == "client" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not authorized"})
		return
	}
	// Verify if database says the same...
	if userVerified(c, session) == false {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not authorized"})
		return
	}
	// If Ok proceed to search
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
