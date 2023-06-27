package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"pert-api/models"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
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

// ListProjectsHandler
// @Summary Get projects list
// @Produce json
// @Success 200 {object} models.Project
// @Failure 401 {string} Unauthorized
// @Failure 403 {string} Not logged in
// @Failure 500 {string} Cant list projects
// @Router /projects [get]
func (handler *ProjectsHandler) ListProjectsHandler(c *gin.Context) {
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
	// Proceed to list
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

// NewProjectHandler
// @Summary Add new project
// @Produce json
// @Param request body models.Project true "JSON for new project"
// @Success 200 {string} Project added
// @Failure 401 {string} Unauthorized
// @Failure 403 {string} Not logged in
// @Router /projects [post]
func (handler *ProjectsHandler) NewProjectHandler(c *gin.Context) {
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

	// Proceed
	var project models.Project
	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Validate all fields
	var fName, fType, fTargets, fFramework bool
	if len(strings.TrimSpace(project.Name)) > 0 && govalidator.IsAlpha(project.Name) {
		fName = true
	}
	if len(strings.TrimSpace(project.Type)) > 0 && govalidator.IsAlpha(project.Type) {
		fType = true
	}
	if len(strings.TrimSpace(project.Targets[0].Name)) > 0 &&
		govalidator.IsAlpha(project.Targets[0].Name) {
		fTargets = true
	}
	if len(strings.TrimSpace(project.Framework.Name)) > 0 &&
		govalidator.IsAlpha(project.Framework.Name) {
		fFramework = true
	}
	// Copy Framework from Default

	// If ok, insert project
	if fName && fType && fTargets && fFramework {
		project.ID = primitive.NewObjectID()
		project.CreatedAt = time.Now()
		_, err := handler.collection.InsertOne(c, project)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while inserting a new project"})
			return
		}
		c.JSON(http.StatusOK, project)
	}
	c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while inserting a new project"})
	return
}

// UpdateProjectHandler
// @Summary Update project information
// @Produce json
// @Param id query string true "Project.ID"
// @Param request body models.Project true "All fields are optional"
// @Success 200 {string} Project has been updated
// @Failure 401 {string} Unauthorized
// @Failure 403 {string} Not logged in
// @Router /projects [put]
func (handler *ProjectsHandler) UpdateProjectHandler(c *gin.Context) {
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

	// Proceed to Update
	id := c.Param("id")
	var project models.Project
	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	objectId, _ := primitive.ObjectIDFromHex(id)
	project.CreatedAt = time.Now()
	filter := bson.D{{"_id", objectId}}
	var fieldUpdated int8

	if len(strings.TrimSpace(project.Name)) > 0 && govalidator.IsAlpha(project.Name) {
		update := bson.D{{"$set", bson.D{{"name", project.Name}}}}
		_, err := handler.collection.UpdateOne(c, filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		fieldUpdated++
	}
	if len(strings.TrimSpace(project.Type)) > 0 && govalidator.IsAlpha(project.Type) {
		update := bson.D{{"$set", bson.D{{"type", project.Type}}}}
		_, err := handler.collection.UpdateOne(c, filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		fieldUpdated++
	}
	// CompanyId
	if len(strings.TrimSpace(project.CompanyId)) > 0 {
		// Check if company really exists
		client, err := mongo.Connect(c, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
		if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
			log.Fatal(err)
		}
		companiesCollection := client.Database(os.Getenv("MONGO_DATABASE")).Collection("companies")
		companiesHandler := NewCompaniesHandler(c, companiesCollection)
		companyId, _ := primitive.ObjectIDFromHex(project.CompanyId)
		companyCursor := companiesHandler.collection.FindOne(c, bson.M{
			"_id": companyId,
		})
		if companyCursor.Err() != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Company does not exist"})
			return
		}

		// If Company exists, proceed to update
		update := bson.D{{"$set", bson.D{{"companyId", project.CompanyId}}}}
		_, err = handler.collection.UpdateOne(c, filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		fieldUpdated++
	}
	// Description
	if len(strings.TrimSpace(project.Description)) > 0 && govalidator.IsAlpha(project.Description) {
		update := bson.D{{"$set", bson.D{{"description", project.Description}}}}
		_, err := handler.collection.UpdateOne(c, filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		fieldUpdated++
	}
	// Status
	if len(strings.TrimSpace(project.Status)) > 0 && govalidator.IsAlpha(project.Status) {
		update := bson.D{{"$set", bson.D{{"status", project.Status}}}}
		_, err := handler.collection.UpdateOne(c, filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		fieldUpdated++
	}
	// Targets
	if len(strings.TrimSpace(project.Targets[0].Name)) > 0 &&
		govalidator.IsAlpha(project.Targets[0].Name) {

		update := bson.D{{"$set", bson.D{{"targets", project.Targets}}}}
		_, err := handler.collection.UpdateOne(c, filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		fieldUpdated++
	}
	// TODO Framework (mus get Framework from Framework Collection)
	/*
		if len(strings.TrimSpace(project.Status)) > 0 && govalidator.IsAlpha(project.Status) {
			update := bson.D{{"$set", bson.D{{"status", project.Status}}}}
			_, err := handler.collection.UpdateOne(c, filter, update)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			fieldUpdated++
		}
	*/

	if fieldUpdated > 0 {
		c.JSON(http.StatusOK, gin.H{"message": "Project has been updated"})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"error": "Error"})
	return

}

// DeleteProjectHandler
// @Summary Delete project
// @Produce json
// @Param id query string true "Project.ID"
// @Success 200 {string} Project has been deleted
// @Failure 401 {string} Not authorized
// @Failure 403 {string} Not logged in
// @Failure 404 {string} Project not found
// @Router /projects [delete]
func (handler *ProjectsHandler) DeleteProjectHandler(c *gin.Context) {
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

	// Proceed to Delete
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

// SearchProjectHandler
// @Summary Search project by Name
// @Produce json
// @Param request body string true "For search" SchemaExample({ "name": "ProjectName" })
// @Success 200 {string} Project has been deleted
// @Failure 401 {string} Not authorized
// @Failure 403 {string} Not logged in
// @Failure 404 {string} Project not found
// @Router /projects/search [get]
func (handler *ProjectsHandler) SearchProjectHandler(c *gin.Context) {
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

	// Proceed to Search
	companyId := c.Query("companyId")

	filter := bson.D{{"companyId", bson.D{{"$eq", companyId}}}}
	cur, err := handler.collection.Find(c, filter)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No Project found"})
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
