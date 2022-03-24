package api

import (
	"log"
	"net/http"
	"nubes/sum/db"
	"nubes/sum/utils"

	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/jwt"
)

func SessionsIndex(c *gin.Context) {
	currentUser := currentUser(c)
	sessions := db.DB.UserSessions().ListForUserID(currentUser.ID)

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

	user, err := db.DB.Users().FindByCredentials(request.Username, request.Password)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "invalid_credentials",
		})
	} else {
		// Genearate and save the new session
		session := user.NewSession(c)
		db.DB.UserSessions().Create(session)

		// Prepare an access token
		token, err := jwt.NewBuilder().
			JwtID(session.ID).
			Subject(user.ID).
			Expiration(session.ExpiresAt).
			Audience([]string{"sessions"}).
			Build()
		if err != nil {
			log.Panicf("Could not build the JWT token: %v", err)
		}

		c.JSON(http.StatusOK, gin.H{
			"access_token": utils.JwtSign(token, "RS256"),
		})
	}
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
