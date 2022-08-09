package v1alpha1

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Subscription struct {
	Base
	TenantID uuid.UUID
	Service  Service `binding:"required,oneof=enabled disabled purged"`
	Title    string  `binding:"required"`
	Username string  `binding:"required"`
	URL      string
	Links    SubscriptionLinks `gorm:"-" json:"links"` // ignores this field
}

type Service string

var Enabled Service = "enabled"
var Disabled Service = "disabled"
var Purged Service = "purged"

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

type SubscriptionLinks struct {
	Self   string `json:"self"`
	Tenant string `json:"tenant"`
}

func (s *Subscription) AddLinks() {
	s.Links = SubscriptionLinks{
		Self:   fmt.Sprintf("/v1alpha1/tenants/%s/subscriptions/%s", s.TenantID.String(), s.ID.String()),
		Tenant: fmt.Sprintf("/v1alpha1/tenants/%s", s.TenantID.String()),
	}
}

func (s *Subscription) SelfLink() string {
	return s.Links.Self
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
	result := a.DB.Where(&Subscription{TenantID: tid}).Find(&subscriptions)
	models := make([]LinkAdder, len(subscriptions))
	for i, _ := range subscriptions {
		models[i] = &subscriptions[i]
	}
	handleListResult(c, result.Error, models)
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
	handlePostResult(c, result, &subscription)
}

func (a *SubscriptionAPI) get(c *gin.Context) {
	pk, ok := parsePK(c)
	if !ok {
		return
	}

	subscription := Subscription{Base: Base{ID: pk}}

	err := a.DB.First(&subscription).Error

	handleGetResult(c, err, &subscription)
}

func (a *SubscriptionAPI) delete(c *gin.Context) {
	pk, ok := parsePK(c)
	if !ok {
		return
	}

	err := a.DB.Delete(&Subscription{Base: Base{ID: pk}}).Error

	handleDeleteResult(c, err)
}
