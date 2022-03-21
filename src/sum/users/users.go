package users

import (
	"nubes/sum/db"

	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {
	db.DB.Users().List("username ASC")
}

func Create(c *gin.Context) {

}

func Update(c *gin.Context) {

}

func Delete(c *gin.Context) {

}

func Show(c *gin.Context) {

}
