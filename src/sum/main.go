package main

import (
	"crypto/rand"
	"crypto/rsa"
	"log"
	"nubes/sum/db"

	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/jwk"
)

var RSAKey *rsa.PrivateKey

func PrepareKeys() {
	var err error

	JWKSet = jwk.NewSet()
	JWKPublicSet = jwk.NewSet()

	RSAKey, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Panicf("failed to generate new RSA key: %s\n", err)
	}

	key, err := jwk.New(RSAKey)
	if err != nil {
		log.Panicf("failed to JWK from RSA: %s\n", err)
	}

	key.Set("use", "sig")
	key.Set("kid", "deadbeef")

	JWKSet.Add(key)

	publicKey, err := key.PublicKey()
	if err != nil {
		log.Panicf("expected jwk.SymmetricKey, got %T\n", key)
		return
	}

	JWKSet.Add(key)
	JWKPublicSet.Add(publicKey)
}

func main() {
	PrepareKeys()
	db.Init()

	router := gin.Default()
	router.LoadHTMLFiles("new_session.html", "error.html")

	router.GET("/.well-known/openid-configuration", openidConfiguration)
	router.GET("/openid/jwks", jwks)
	router.POST("/openid/registration", registration)

	router.GET("/openid/authorization", authorizationStart)
	router.POST("/openid/authorization", authorizationSubmit)

	router.POST("/openid/token", token)

	router.GET("/openid/userinfo", userinfo)

	router.Run("localhost:8080")
}
