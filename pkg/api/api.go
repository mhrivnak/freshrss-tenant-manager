package api

import "github.com/gin-gonic/gin"

type RouteAdder interface {
	AddRoutes(*gin.Engine)
}
