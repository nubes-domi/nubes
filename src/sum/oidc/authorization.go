package oidc

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"nubes/sum/db"
	"nubes/sum/utils"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lestrrat-go/jwx/jwt"
)

func CreateAuthorizationRequest(c *gin.Context) {
	request := buildAuthorizationRequest(c)
	db.DB.OidcAuthorizationRequests().Create(&request)

	client, err := db.DB.OidcClients().FindById(request.ClientID)
	if err != nil {
		c.HTML(http.StatusOK, "error.html", gin.H{
			"error": "Invalid Client ID.",
		})
		return
	}

	if !utils.Contains(client.RedirectURIs, request.RedirectURI) {
		c.HTML(http.StatusOK, "error.html", gin.H{
			"error": "The application specified a Redirect URI which has not been whitelisted.",
		})
		return
	}

	if request.ResponseType == "" {
		c.Redirect(302, buildRedirection(request.RedirectURI, map[string]string{
			"error": "invalid_request",
			"state": request.State,
		}))
	} else {
		c.Redirect(302, "/openid/authorization/"+request.ID)
	}
}

func ShowAuthorizationRequest(c *gin.Context) {
	// TODO: verify granted permissions, ask if needed, then redirect
	ConfirmAuthorizationRequest(c)

	// id := c.Param("id")
	// auth, _ := db.DB.OidcAuthorizationRequests().FindById(id)
	// client, _ := db.DB.OidcClients().FindById(auth.ClientID)

	// logo := client.GetLogoURI("")
	// policyUri := client.GetPolicyURI("")
	// tosUri := client.GetTosURI("")
	// c.HTML(http.StatusOK, "new_session.html", gin.H{
	// 	"client":  client,
	// 	"request": auth,
	// 	"logo":    logo,
	// 	"policy":  policyUri,
	// 	"tos":     tosUri,
	// })
}

func ConfirmAuthorizationRequest(c *gin.Context) {
	id := c.Param("id")
	auth, _ := db.DB.OidcAuthorizationRequests().FindById(id)
	redirectURI := auth.RedirectURI

	if !strings.Contains(redirectURI, "?") {
		redirectURI += "?"
	} else {
		redirectURI += "&"
	}

	db.DB.OidcAuthorizationRequests().Delete(auth)

	c.Redirect(http.StatusFound, buildRedirection(redirectURI, map[string]string{
		"state": auth.State,
		"code":  generateAuthorizationCode("123", auth.Nonce, auth.ClientID),
	}))
}

func buildAuthorizationRequest(c *gin.Context) db.OidcAuthorizationRequest {
	authRequest := db.OidcAuthorizationRequest{
		ID: uuid.New().String(),
	}
	c.Bind(&authRequest)

	var requestObject []byte
	if authRequest.RequestURI != "" {
		if authRequest.Request != "" {
			log.Panicf("Cannot use both request_uri and request")
		}

		resp, err := http.Get(authRequest.RequestURI)
		if err != nil {
			log.Panicf("%v", err)
		}

		requestObject, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Panicf("%v", err)
		}
	}

	if authRequest.Request != "" {
		requestObject = []byte(authRequest.Request)
	}

	if len(requestObject) > 0 {
		token, err := jwt.Parse(requestObject) //, jwt.WithVerify(jwa.RS256, RSAKey.PublicKey))
		if err != nil {
			log.Panicf("%v", err)
		}

		if val, ok := token.Get("response_type"); ok {
			authRequest.ResponseType = val.(string)
		}

		if val, ok := token.Get("redirect_uri"); ok {
			authRequest.RedirectURI = val.(string)
		}

		if val, ok := token.Get("state"); ok {
			authRequest.State = val.(string)
		}

		if val, ok := token.Get("scope"); ok {
			authRequest.Scope = val.(string)
		}

		if val, ok := token.Get("client_id"); ok {
			authRequest.ClientID = val.(string)
		}

		if val, ok := token.Get("response_mode"); ok {
			authRequest.ResponseMode = val.(string)
		}

		if val, ok := token.Get("nonce"); ok {
			authRequest.Nonce = val.(string)
		}

		if val, ok := token.Get("display"); ok {
			authRequest.Display = val.(string)
		}

		if val, ok := token.Get("prompt"); ok {
			authRequest.Prompt = val.(string)
		}

		if val, ok := token.Get("max_age"); ok {
			authRequest.MaxAge = val.(uint)
		}

		if val, ok := token.Get("ui_locales"); ok {
			authRequest.UILocales = val.(string)
		}

		if val, ok := token.Get("id_token_hint"); ok {
			authRequest.IDTokenHint = val.(string)
		}

		if val, ok := token.Get("login_hint"); ok {
			authRequest.LoginHint = val.(string)
		}

		if val, ok := token.Get("acr_values"); ok {
			authRequest.ACRValues = val.(string)
		}

		// TODO: claims
	}

	return authRequest
}

func generateAuthorizationCode(subject, nonce, clientId string) string {
	token, err := jwt.NewBuilder().
		Expiration(time.Now().Add(time.Minute*10)).
		JwtID(uuid.New().String()).
		Subject(subject).
		Audience([]string{"token"}).
		Issuer("authorization").
		Claim("nonce", nonce).
		Claim("client_id", clientId).
		Build()

	if err != nil {
		log.Panicf("Could not build a JWT: %v", err)
	}

	return utils.JwtSign(token, "RS256")
}

func buildRedirection(baseURI string, additionalParams map[string]string) string {
	if len(additionalParams) == 0 {
		return baseURI
	}

	if !strings.Contains(baseURI, "?") {
		baseURI += "?"
	} else {
		baseURI += "&"
	}

	for k, v := range additionalParams {
		baseURI += url.QueryEscape(k) + "=" + url.QueryEscape(v) + "&"
	}

	return baseURI
}
