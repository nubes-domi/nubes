package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AppsIndex(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"apps": []string{},
	})
}

func AppsUpdate(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"ok": true,
	})
}

func AppsDelete(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"ok": true,
	})
}
