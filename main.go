// @title PERT API
// @version 1.0
// @description Pentest Execution and Reporting Tool
//
// @contact.name Fabio Almeida
// @contact.email mentesan@gmail.com
// @license MIT
// @license.url https://github.com/mentesan/pert/blob/main/LICENSE

// @Basepath /
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	handlers "pert/handlers"

	"github.com/gin-contrib/sessions"
	redisStore "github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var authHandler *handlers.AuthHandler

// var contactsHandler *handlers.ContactsHandler
var companiesHandler *handlers.CompaniesHandler
var projectsHandler *handlers.ProjectsHandler

func init() {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB")
	// contactsCollection := client.Database(os.Getenv("MONGO_DATABASE")).Collection("contacts")
	usersCollection := client.Database(os.Getenv("MONGO_DATABASE")).Collection("users")
	companiesCollection := client.Database(os.Getenv("MONGO_DATABASE")).Collection("companies")
	projectsCollection := client.Database(os.Getenv("MONGO_DATABASE")).Collection("projects")

	// Authentication handler
	authHandler = handlers.NewAuthHandler(ctx, usersCollection)
	//	contactsHandler = handlers.NewContactsHandler(ctx, contactsCollection )
	companiesHandler = handlers.NewCompaniesHandler(ctx, companiesCollection)
	projectsHandler = handlers.NewProjectsHandler(ctx, projectsCollection)

	/*
		redisClient := redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		})
	*/

	fmt.Println("END INIT")
}

func main() {
	router := gin.Default()

	store, _ := redisStore.NewStore(10, "tcp", "localhost:6379", "", []byte("secret"))
	router.Use(sessions.Sessions("pert", store))

	// Public endpoints
	router.POST("/signin", authHandler.SignInHandler)
	router.POST("/refresh", authHandler.RefreshHandler)
	router.POST("/signout", authHandler.SignOutHandler)

	authorized := router.Group("/")
	authorized.Use(authHandler.AuthMiddleware())
	{
		authorized.GET("/companies", companiesHandler.ListCompaniesHandler)
		authorized.POST("/companies", companiesHandler.NewCompanyHandler)
		authorized.PUT("/companies/:id", companiesHandler.UpdateCompanyHandler)
		authorized.DELETE("/companies/:id", companiesHandler.DeleteCompanyHandler)
		authorized.GET("/companies/search", companiesHandler.SearchCompanyHandler)

		authorized.GET("/projects", projectsHandler.ListProjectsHandler)
		authorized.GET("/projects/:companyId", projectsHandler.ListProjectsHandler)
		authorized.POST("/projects", projectsHandler.NewProjectHandler)
		authorized.PUT("/projects/:id", projectsHandler.UpdateProjectHandler)
		authorized.DELETE("/projects/:id", projectsHandler.DeleteProjectHandler)
		authorized.GET("/projects/search", projectsHandler.SearchProjectHandler)
	}
	router.Run()
}
