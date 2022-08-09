package v1alpha1

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ServiceLevel struct {
	Base
	Name        string `gorm:"unique" binding:"required"`
	Description string
	Price       uint
	Links       ServiceLevelLinks `gorm:"-" json:"links"` // ignores this field
}

type ServiceLevelLinks struct {
	Self string `json:"self"`
}

func (s *ServiceLevel) AddLinks(base string) {
	s.Links = ServiceLevelLinks{
		Self: fmt.Sprintf("%s/v1alpha1/servicelevels/%s", base, s.ID.String()),
	}
}

func (s *ServiceLevel) SelfLink() string {
	return s.Links.Self
}

type ServiceLevelAPI struct {
	DB *gorm.DB
}

func (s *ServiceLevelAPI) AddRoutes(router *gin.Engine) {
	router.GET("/v1alpha1/servicelevels/", s.list)
	router.GET("/v1alpha1/servicelevels/:id", s.get)
	router.DELETE("/v1alpha1/servicelevels/:id", s.delete)
	router.POST("/v1alpha1/servicelevels/", s.post)
}

func (s *ServiceLevelAPI) list(c *gin.Context) {
	levels := []ServiceLevel{}
	result := s.DB.Find(&levels)
	models := make([]LinkAdder, len(levels))
	for i, _ := range levels {
		models[i] = &levels[i]
	}
	handleListResult(c, result.Error, models)
}

func (s *ServiceLevelAPI) post(c *gin.Context) {
	level := ServiceLevel{}
	err := c.BindJSON(&level)
	if err != nil {
		return
	}

	result := s.DB.Create(&level)
	handlePostResult(c, result, &level)
}

func (s *ServiceLevelAPI) get(c *gin.Context) {
	pk, ok := parsePK(c)
	if !ok {
		return
	}

	level := ServiceLevel{Base: Base{ID: pk}}

	err := s.DB.First(&level).Error

	handleGetResult(c, err, &level)
}

func (s *ServiceLevelAPI) delete(c *gin.Context) {
	pk, ok := parsePK(c)
	if !ok {
		return
	}

	err := s.DB.Delete(&ServiceLevel{Base: Base{ID: pk}}).Error

	handleDeleteResult(c, err)
}
