package v1alpha1

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Subscription struct {
	Base
	TenantID uuid.UUID
	Title    string
	Username string
	URL      string
}

func (s *Subscription) notify() error {
	go func() {
		// TODO call a configured webhook URL
		log.Println("Pretending to call webhook")
	}()

	return nil
}

func (s *Subscription) AfterCreate(tx *gorm.DB) error {
	return s.notify()
}

func (s *Subscription) AfterUpdate(tx *gorm.DB) error {
	return s.notify()
}

func (s *Subscription) AfterDelete(tx *gorm.DB) error {
	return s.notify()
}

type SubscriptionAPI struct {
	DB *gorm.DB
}

func (a *SubscriptionAPI) AddRoutes(router *gin.Engine) {
	router.GET("/v1alpha1/tenants/:tid/subscriptions/", a.list)
	router.GET("/v1alpha1/tenants/:tid/subscriptions/:id", a.get)
	router.DELETE("/v1alpha1/tenants/:tid/subscriptions/:id", a.delete)
	router.POST("/v1alpha1/tenants/:tid/subscriptions/", a.post)
}

func (a *SubscriptionAPI) list(c *gin.Context) {
	tid, ok := parseUUID(c, "tid")
	if !ok {
		return
	}

	subscriptions := []Subscription{}
	a.DB.Where(&Subscription{TenantID: tid}).Find(&subscriptions)
	c.IndentedJSON(http.StatusOK, &subscriptions)
}

func (a *SubscriptionAPI) post(c *gin.Context) {
	tid, ok := parseUUID(c, "tid")
	if !ok {
		return
	}

	subscription := Subscription{}
	err := c.BindJSON(&subscription)
	// TODO improve error handling here and in other POST handlers to detect
	// JSON parse failure
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	subscription.TenantID = tid
	result := a.DB.Create(&subscription)
	handlePostResult(c, result, subscription)
}

func (a *SubscriptionAPI) get(c *gin.Context) {
	pk, ok := parsePK(c)
	if !ok {
		return
	}

	subscription := Subscription{Base: Base{ID: pk}}

	err := a.DB.First(&subscription).Error

	handleGetResult(c, err, subscription)
}

func (a *SubscriptionAPI) delete(c *gin.Context) {
	pk, ok := parsePK(c)
	if !ok {
		return
	}

	err := a.DB.Delete(&Subscription{Base: Base{ID: pk}}).Error

	handleDeleteResult(c, err)
}
