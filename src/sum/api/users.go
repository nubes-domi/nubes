package api

import (
	"net/http"
	"nubes/sum/db"
	"nubes/sum/utils"

	"github.com/gin-gonic/gin"
)

func currentUser(c *gin.Context) *db.User {
	return utils.CtxMustGet[*db.User](c, "currentUser")
}

func UsersIndex(c *gin.Context) {
	currentUser := currentUser(c)
	if !currentUser.IsAdmin {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "access_denied",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users": db.DB.Users().List("username ASC"),
	})
}

func UsersCreate(c *gin.Context) {
	currentUser := currentUser(c)
	if !currentUser.IsAdmin {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "access_denied",
		})
		return
	}

	user := db.DB.Users().New()
	c.BindJSON(user)
	db.DB.Users().Create(user)

	c.JSON(http.StatusOK, &user)
}

func UsersUpdate(c *gin.Context) {
	user, err := db.DB.Users().FindById(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "not_found",
		})
		return
	}

	currentUser := currentUser(c)
	if user.ID != currentUser.ID && !currentUser.IsAdmin {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "access_denied",
		})
		return
	}

	c.BindJSON(user)
	db.DB.Users().Update(user)

	c.JSON(http.StatusOK, &user)
}

func UsersDelete(c *gin.Context) {
	user, err := db.DB.Users().FindById(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "not_found",
		})
		return
	}

	currentUser := currentUser(c)
	if !currentUser.IsAdmin {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "access_denied",
		})
		return
	}

	if currentUser.ID == user.ID {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":             "bad_request",
			"error_description": "Cannot delete self",
		})
	}
}

func UsersShow(c *gin.Context) {
	user, err := db.DB.Users().FindById(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "not_found",
		})
		return
	}

	currentUser := currentUser(c)
	if user.ID != currentUser.ID && !currentUser.IsAdmin {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "access_denied",
		})
		return
	}

	c.JSON(http.StatusOK, &user)
}
