package v1alpha1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ServiceLevel struct {
	UUID        string `gorm:"primaryKey"`
	Name        string `gorm:"unique"`
	Description string
	Price       uint
}

type ServiceLevelAPI struct {
	DB *gorm.DB
}

func (s *ServiceLevelAPI) AddRoutes(router *gin.Engine) {
	router.GET("/v1alpha1/servicelevels/", s.list)
	router.GET("/v1alpha1/servicelevels/:uuid", s.get)
	router.DELETE("/v1alpha1/servicelevels/:uuid", s.delete)
	router.POST("/v1alpha1/servicelevels/", s.post)
}

func (s *ServiceLevelAPI) list(c *gin.Context) {
	levels := []ServiceLevel{}
	s.DB.Find(&levels)
	c.IndentedJSON(http.StatusOK, &levels)
}

func (s *ServiceLevelAPI) post(c *gin.Context) {
	level := ServiceLevel{}
	err := c.BindJSON(&level)
	if err != nil {
		return
	}

	level.UUID = uuid.New().String()

	result := s.DB.Create(&level)
	if result.Error != nil {
		return
	}

	c.IndentedJSON(http.StatusCreated, level)
}

func (s *ServiceLevelAPI) get(c *gin.Context) {
	pk, ok := parsePK(c)
	if !ok {
		return
	}

	level := ServiceLevel{}

	err := s.DB.First(&level, "uuid = ?", pk).Error

	handleGetResult(c, err, level)
}

func (s *ServiceLevelAPI) delete(c *gin.Context) {
	pk, ok := parsePK(c)
	if !ok {
		return
	}

	err := s.DB.Delete(&ServiceLevel{UUID: pk}).Error

	handleDeleteResult(c, err)
}
