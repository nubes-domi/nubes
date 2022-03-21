package oidc

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"nubes/sum/db"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRegistrationCode(t *testing.T) {
	db.InitDatabase()

	r := gin.Default()
	r.POST("/openid/registration", Registration)

	req, _ := http.NewRequest("POST", "/openid/registration", strings.NewReader("{\"grant_types\":[\"authorization_code\"],\"jwks\":{\"keys\":[{\"kty\":\"RSA\",\"e\":\"AQAB\",\"use\":\"sig\",\"alg\":\"RS256\",\"n\":\"j98uGQ\"}]},\"response_types\":[\"code\"],\"redirect_uris\":[\"https://www.certification.openid.net/test/aXAWKO1B6OTRy1y/callback\"]}"))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	body := make(map[string]interface{})
	response, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Fail()
	}

	err = json.Unmarshal(response, &body)
	if err != nil {
		t.Fail()
	}

	clientId := body["client_id"].(string)
	client, err := db.DB.OidcClients().FindById(clientId)

	if client.Jwks.Set.Len() != 1 {
		t.Fail()
	}
}

func TestRegistrationHybrid(t *testing.T) {
	db.InitDatabase()

	r := gin.Default()
	r.POST("/openid/registration", Registration)

	req, _ := http.NewRequest("POST", "/openid/registration", strings.NewReader("{\"grant_types\":[\"implicit\",\"authorization_code\"],\"response_types\":[\"code id_token\"],\"redirect_uris\":[\"https://example.com/callback\"]}"))

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	body := make(map[string]interface{})
	response, err := ioutil.ReadAll(w.Body)
	if err != nil {
		t.Fail()
	}

	err = json.Unmarshal(response, &body)
	if err != nil {
		t.Fail()
	}

	clientId := body["client_id"].(string)
	_, err = db.DB.OidcClients().FindById(clientId)
	if err != nil {
		t.Fail()
	}
}
