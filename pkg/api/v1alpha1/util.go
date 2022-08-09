package v1alpha1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func parseUUID(c *gin.Context, field string) (uuid.UUID, bool) {
	pk := c.Param(field)
	ret, err := uuid.Parse(pk)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "invalid UUID"})
		return ret, false
	}
	return ret, true
}

func parsePK(c *gin.Context) (uuid.UUID, bool) {
	return parseUUID(c, "id")
}

func handleGetResult(c *gin.Context, err error, model LinkAdder) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "not found"})
		return
	} else if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "error"})
		return
	}

	model.AddLinks()
	c.IndentedJSON(http.StatusOK, &model)
}

func handleListResult(c *gin.Context, err error, models []LinkAdder) {
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "error getting data"})
		return
	}
	for i, _ := range models {
		models[i].AddLinks()
	}
	c.IndentedJSON(http.StatusOK, &models)
}

func handlePostResult(c *gin.Context, result *gorm.DB, model LinkAdder) {
	err := result.Error
	if err != nil {
		// TODO set a better response code and message
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	model.AddLinks()
	c.Header("Location", model.SelfLink())
	c.IndentedJSON(http.StatusCreated, &model)
}

func handleDeleteResult(c *gin.Context, err error) {
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "error"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "deleted"})
}
