package oidc

import (
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
	db.DB.OidcClients().Create(&client)

	c.IndentedJSON(http.StatusCreated, RegistrationResponse{
		ClientID:              client.ID,
		ClientSecret:          "a-fake-secret",
		ClientSecretExpiresAt: 0,
	})
}
