package api

import (
	"net/http"
	"nubes/sum/services/sessions"

	"github.com/gin-gonic/gin"
)

func SessionsIndex(c *gin.Context) {
	currentUser := currentUser(c)
	sessions, _ := sessions.List(currentUser, "")

	c.JSON(http.StatusOK, gin.H{
		"sessions": sessions,
	})
}

type SessionCreateRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func SessionsCreate(c *gin.Context) {
	request := SessionCreateRequest{}
	err := c.BindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "bad_request",
			"details": err.Error(),
		})
		return
	}

	session, err := sessions.Create(request.Username, request.Password, "", "")
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": err,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"access_token": session.SignedToken,
		})
	}
}

func SessionsDelete(c *gin.Context) {
	currentUser := currentUser(c)
	sessions.Delete(currentUser, c.Param("id"))

	c.JSON(http.StatusOK, gin.H{
		"ok": true,
	})
}
