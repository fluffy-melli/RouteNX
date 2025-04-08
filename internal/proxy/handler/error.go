package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func IPBlock(c *gin.Context) {
	c.HTML(http.StatusForbidden, "error.tmpl", gin.H{
		"error": "CIDR IP Block",
	})
}

func NoRoute(c *gin.Context) {
	c.HTML(http.StatusBadGateway, "error.tmpl", gin.H{
		"error": "No Domain Route",
	})
}

func InternalError(c *gin.Context, err error) {
	c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{
		"error": err.Error(),
	})
}
