package handlers

import (
	"context"
	"log"
	"net/http"
	"pert-api/models"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	collection *mongo.Collection
	ctx        context.Context
}
type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

type JWTOutput struct {
	Token   string    `json:"token"`
	Expires time.Time `json:"expires"`
}

func NewAuthHandler(ctx context.Context, collection *mongo.Collection) *AuthHandler {
	return &AuthHandler{
		collection: collection,
		ctx:        ctx,
	}
}

func (handler *AuthHandler) SignInHandler(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cur := handler.collection.FindOne(handler.ctx, bson.M{
		"email": user.Email,
	})
	if cur.Err() != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}
	// Get "User Type(admin|pentester|client)" from database
	var actualUser models.User
	cur.Decode(&actualUser)

	log.Println("PASSWORDS DB/Informed:", actualUser.Password, user.Password)
	// Compare password
	err := bcrypt.CompareHashAndPassword([]byte(actualUser.Password), []byte(user.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Session
	sessionToken := xid.New().String()
	session := sessions.Default(c)
	session.Set("email", user.Email)
	session.Set("token", sessionToken)
	session.Set("type", actualUser.Type)
	log.Println("Token: ", sessionToken)
	log.Println("Type: ", actualUser.Type)
	session.Save()
	msg := "User signed in: " + sessionToken
	c.JSON(http.StatusOK, gin.H{"message": msg})
}

func (handler *AuthHandler) RefreshHandler(c *gin.Context) {
	session := sessions.Default(c)
	sessionToken := session.Get("token")
	sessionEmail := session.Get("email")
	if sessionToken == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid session cookie"})
		return
	}

	// Session
	sessionToken = xid.New().String()
	session.Set("email", sessionEmail.(string))
	session.Set("token", sessionToken)
	session.Save()

	c.JSON(http.StatusOK, gin.H{"message": "New session issued"})
}

func (handler *AuthHandler) SignOutHandler(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.JSON(http.StatusOK, gin.H{"message": "Signed out..."})
}

func (handler *AuthHandler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		sessionToken := session.Get("token")
		if sessionToken == nil {
			c.JSON(http.StatusForbidden, gin.H{"message": "Not logged"})
			c.Abort()
		}
		c.Next()
	}
}
