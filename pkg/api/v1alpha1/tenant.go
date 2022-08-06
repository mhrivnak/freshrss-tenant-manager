package v1alpha1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Tenant struct {
	UUID string `gorm:"primaryKey"`
	Name string `gorm:"unique"`
}

type TenantAPI struct {
	DB *gorm.DB
}

func (t *TenantAPI) AddRoutes(router *gin.Engine) {
	router.GET("/v1alpha1/tenants/", t.list)
	router.GET("/v1alpha1/tenants/:uuid", t.get)
	router.DELETE("/v1alpha1/tenants/:uuid", t.delete)
	router.POST("/v1alpha1/tenants/", t.post)
}

func (t *TenantAPI) list(c *gin.Context) {
	levels := []Tenant{}
	t.DB.Find(&levels)
	c.IndentedJSON(http.StatusOK, &levels)
}

func (t *TenantAPI) post(c *gin.Context) {
	level := Tenant{}
	err := c.BindJSON(&level)
	if err != nil {
		return
	}

	level.UUID = uuid.New().String()

	result := t.DB.Create(&level)
	if result.Error != nil {
		return
	}

	c.IndentedJSON(http.StatusCreated, level)
}

func (t *TenantAPI) get(c *gin.Context) {
	pk, ok := parsePK(c)
	if !ok {
		return
	}

	level := Tenant{}

	err := t.DB.First(&level, "uuid = ?", pk).Error

	handleGetResult(c, err, level)
}

func (t *TenantAPI) delete(c *gin.Context) {
	pk, ok := parsePK(c)
	if !ok {
		return
	}

	err := t.DB.Delete(&Tenant{UUID: pk}).Error

	handleDeleteResult(c, err)
}
