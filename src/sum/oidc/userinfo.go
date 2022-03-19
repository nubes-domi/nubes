package oidc

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserinfoResponse struct {
	Sub string `json:"sub"`
}

func Userinfo(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, UserinfoResponse{
		Sub: "123",
	})
}
