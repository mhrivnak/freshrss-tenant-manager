package v1alpha1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func parsePK(c *gin.Context) (string, bool) {
	pk := c.Param("uuid")
	_, err := uuid.Parse(pk)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "invalid UUID"})
		return "", false
	}
	return pk, true
}

func handleGetResult(c *gin.Context, err error, model interface{}) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "not found"})
		return
	} else if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "error"})
		return
	}

	c.IndentedJSON(http.StatusOK, &model)
}

func handlePostResult(c *gin.Context, result *gorm.DB, model interface{}) {
	err := result.Error
	if err != nil {
		// TODO set a better response code and message
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusCreated, &model)
}

func handleDeleteResult(c *gin.Context, err error) {
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "error"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"message": "deleted"})
}
