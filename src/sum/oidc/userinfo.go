package oidc

import (
	"fmt"
	"net/http"
	"nubes/sum/db"
	"nubes/sum/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/jwt"
)

func Userinfo(c *gin.Context) {
	accessTokenString := retrieveAccessToken(c)
	accessToken, err := db.DB.OidcAccessTokens().Find(accessTokenString)
	if err != nil {
		c.Header("WWW-Authenticate", "Bearer, error=\"invalid_token\"")
		c.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	client, _ := db.DB.OidcClients().FindById(accessToken.ClientID)
	user, _ := db.DB.Users().FindById(accessToken.UserID)

	response := map[string]interface{}{
		"iss": baseURI(c),
		"aud": client.ID,
		"sub": fmt.Sprintf("%d", user.ID),
	}

	if client.UserinfoSignedResponseAlg == "none" && client.UserinfoEncryptedResponseAlg == "none" {
		c.IndentedJSON(http.StatusOK, response)
	} else {
		c.Header("Content-Type", "application/jwt")

		token := jwt.New()
		for k, v := range response {
			token.Set(k, v)
		}

		signed := utils.JwtSign(token, client.UserinfoSignedResponseAlg)
		c.Writer.Write([]byte(signed))
	}
}

func retrieveAccessToken(c *gin.Context) string {
	bearer := getBearer(c)
	if bearer != "" {
		return bearer
	}

	if c.Request.Method != "GET" {
		if len(c.Request.Header["Content-Type"]) > 0 && strings.HasPrefix(c.Request.Header["Content-Type"][0], "application/x-www-form-urlencoded") {
			form := c.PostForm("access_token")
			if form != "" {
				return form
			}
		}
	}

	return c.Query("access_token")
}
