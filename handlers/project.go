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

type ProjectsHandler struct {
	collection *mongo.Collection
	ctx        context.Context
}

func NewProjectsHandler(ctx context.Context, collection *mongo.Collection) *ProjectsHandler {

	return &ProjectsHandler{
		collection: collection,
		ctx:        ctx,
	}
}

func (handler *ProjectsHandler) ListProjectsHandler(c *gin.Context) {
	companyId := c.Param("companyId")
	var cur *mongo.Cursor
	var err error
	// If companyId is passed
	if companyId != "" {
		cur, err = handler.collection.Find(handler.ctx, bson.M{"companyId": companyId})
	} else {
		cur, err = handler.collection.Find(handler.ctx, bson.M{})
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "err.Error()"})
		return
	}
	defer cur.Close(handler.ctx)

	projects := make([]models.Project, 0)
	for cur.Next(handler.ctx) {
		var project models.Project
		cur.Decode(&project)
		projects = append(projects, project)
	}
	c.JSON(http.StatusOK, projects)
}

func (handler *ProjectsHandler) NewProjectHandler(c *gin.Context) {
	var project models.Project
	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return

	}
	project.ID = primitive.NewObjectID()
	project.CreatedAt = time.Now()
	_, err := handler.collection.InsertOne(c, project)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while inserting a new project"})
		return
	}
	log.Println("Remove data from Redis")
	c.JSON(http.StatusOK, project)
}

func (handler *ProjectsHandler) UpdateProjectHandler(c *gin.Context) {
	id := c.Param("id")
	var project models.Project
	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	objectId, _ := primitive.ObjectIDFromHex(id)
	fmt.Println("OBJECT ID:", objectId)

	project.CreatedAt = time.Now()
	filter := bson.D{{"_id", objectId}}
	update := bson.D{{"$set", bson.D{
		{"name", project.Name},
		{"type", project.Type},
		{"description", project.Description},
		{"status", project.Status},
		{"contacts", project.Contacts},
		{"targets", project.Targets},
		{"vulns", project.Vulns},
		{"createdAt", project.CreatedAt},
	}}}
	_, err := handler.collection.UpdateOne(c, filter, update)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Project has been updated"})
}

func (handler *ProjectsHandler) DeleteProjectHandler(c *gin.Context) {
	id := c.Param("id")
	objectId, _ := primitive.ObjectIDFromHex(id)

	filter := bson.D{{"_id", objectId}}
	_, err := handler.collection.DeleteOne(c, filter)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Project has beem deleted"})
}

func (handler *ProjectsHandler) SearchProjectHandler(c *gin.Context) {
	firstName := c.Query("firstName")

	filter := bson.D{{"firstName", bson.D{{"$eq", firstName}}}}
	cur, err := handler.collection.Find(c, filter)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}
	defer cur.Close(c)

	var projects []models.Project
	if err = cur.All(c, &projects); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Error"})
		return
	}

	c.JSON(http.StatusOK, projects)
}
