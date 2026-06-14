package handlers

import (
	"net/http"

	"fd-api/openapi"

	"github.com/gin-gonic/gin"
)

// NewOpenAPIHandler returns the OpenAPI 3.1 JSON document.
func NewOpenAPIHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, openapi.Spec())
	}
}
