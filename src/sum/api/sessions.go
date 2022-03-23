package api

import (
	"net/http"
	"nubes/sum/db"

	"github.com/gin-gonic/gin"
)

func SessionsIndex(c *gin.Context) {
	currentUser := currentUser(c)
	sessions := db.DB.UserSessions().ListForUserID(currentUser.ID)

	c.JSON(http.StatusOK, gin.H{
		"sessions": sessions,
	})
}

func SessionsDelete(c *gin.Context) {
	currentUser := currentUser(c)
	if err := db.DB.UserSessions().DeleteFor(c.Param("id"), currentUser.ID); err != nil {
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "not_found",
			})
			return
		}
	} else {
		c.JSON(http.StatusOK, gin.H{
			"ok": true,
		})
	}
}
