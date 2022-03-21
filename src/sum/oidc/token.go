package oidc

import (
	"fmt"
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
	token, err := utils.JwtVerify(c.PostForm("code"))
	if err != nil {
		log.Panicf("Invalid authorization code: %v", err)
	}

	auth, err := db.DB.OidcAuthorizationRequests().FindByIdAndStage(token.JwtID(), "code")
	if err != nil {
		log.Panicf("Invalid authorization code: %v", err)
	}

	c.IndentedJSON(http.StatusCreated, TokenResponse{
		AccessToken:  generateAccessToken(auth.ClientID, auth.UserID, auth.Scope),
		TokenType:    "Bearer",
		RefreshToken: "",
		ExpiresIn:    3600,
		IdToken:      generateIdToken(c, auth, map[string]string{}),
	})
}

func generateIdToken(c *gin.Context, auth *db.OidcAuthorizationRequest, additionalClaims map[string]string) string {
	session, err := db.DB.UserSessions().FindById(auth.SessionID)
	if err != nil {
		log.Panicf("%v", err)
	}

	id := openid.New()
	id.Set(jwt.IssuerKey, baseURI(c))
	id.Set(jwt.ExpirationKey, time.Now().Add(time.Hour))
	id.Set(jwt.IssuedAtKey, time.Now())
	id.Set(jwt.SubjectKey, fmt.Sprintf("%d", auth.UserID))
	id.Set(jwt.AudienceKey, auth.ClientID)

	id.Set("nonce", auth.Nonce)
	id.Set("auth_time", session.UpdatedAt.Unix())

	id.Set("acr", "0")
	id.Set("amr", []string{"pwd"})

	for k, v := range additionalClaims {
		id.Set(k, v)
	}

	client := db.LoadClient(auth.ClientID)
	return utils.JwtSign(id, client.IDTokenSignedResponseAlg)
}
