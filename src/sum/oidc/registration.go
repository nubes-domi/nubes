package oidc

import (
	"log"
	"net/http"
	"nubes/sum/db"

	"github.com/gin-gonic/gin"
)

type RegistrationResponse struct {
	ClientID                string `json:"client_id"`
	ClientSecret            string `json:"client_secret,omitempty"`
	RegistrationAccessToken string `json:"registration_access_token,omitempty"`
	RegistrationClientURI   string `json:"registration_client_uri,omitempty"`
	ClientIDIssuedAt        string `json:"client_id_issued_at,omitempty"`
	ClientSecretExpiresAt   int    `json:"client_secret_expires_at"`
}

func Registration(c *gin.Context) {
	client := db.BuildOpenIDClient(c)
	err := client.Validate()
	if err != nil {
		result := make(map[string]string)
		if err.Error() == "invalid_redirect_uri" {
			result["error"] = err.Error()
		} else {
			result["error"] = "invalid_client_metadata"
			result["error_description"] = err.Error()
		}

		c.IndentedJSON(http.StatusBadRequest, result)
	} else {
		err := db.DB.OidcClients().Create(&client)
		if err != nil {
			log.Panicf("%v", err)
		}

		client.RegistrationClientURI = baseURI(c) + "/openid/registration/" + client.ID
		c.IndentedJSON(http.StatusCreated, client)
	}
}

func GetClient(c *gin.Context) {
	id := c.Param("id")
	client, err := db.DB.OidcClients().FindById(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, map[string]string{
			"error": "client_not_found",
		})
	}

	c.IndentedJSON(http.StatusOK, client)
}

func DeleteClient(c *gin.Context) {
	id := c.Param("id")
	client, err := db.DB.OidcClients().FindById(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, map[string]string{
			"error": "invalid client id or access token",
		})
	}

	if client.VerifyRegistrationToken(getBearer(c)) {
		db.DB.OidcClients().Delete(client)
		c.Writer.WriteHeader(http.StatusNoContent)
	} else {
		c.IndentedJSON(http.StatusNotFound, map[string]string{
			"error": "invalid client id or access token",
		})
	}
}

func getBearer(c *gin.Context) string {
	if len(c.Request.Header["Authorization"]) > 0 {
		val := c.Request.Header["Authorization"][0]
		if val[:7] == "Bearer " {
			return val[7:]
		}
	}

	return ""
}
