package oidc

import (
	"log"
	"net/http"
	"nubes/sum/db"
	"nubes/sum/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/jwt"
	"github.com/lestrrat-go/jwx/jwt/openid"
)

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	ExpiresIn    int    `json:"expires_in,omitempty"`
	IdToken      string `json:"id_token"`
}

func Token(c *gin.Context) {
	id := openid.New()
	id.Set(jwt.IssuerKey, baseURI(c))
	id.Set(jwt.ExpirationKey, time.Now().Add(time.Hour))
	id.Set(jwt.IssuedAtKey, time.Now())

	token, err := utils.JwtVerify(c.PostForm("code"))
	if err != nil {
		log.Panicf("Invalid authorization code: %v", err)
	}

	id.Set(jwt.SubjectKey, token.Subject())
	clientId, _ := token.Get("client_id")
	id.Set(jwt.AudienceKey, clientId)
	nonce, _ := token.Get("nonce")
	id.Set("nonce", nonce)

	client := db.LoadClient(clientId.(string))
	signed := utils.JwtSign(id, client.IDTokenSignedResponseAlg)

	c.IndentedJSON(http.StatusCreated, TokenResponse{
		AccessToken:  utils.RandBase64(16),
		TokenType:    "Bearer",
		RefreshToken: "",
		ExpiresIn:    3600,
		IdToken:      string(signed[:]),
	})
}
