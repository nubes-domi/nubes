package router

import (
	"nubes/sum/db"
	"nubes/sum/oidc"
	"nubes/sum/sessions"
	"nubes/sum/users"

	"github.com/gin-gonic/gin"
)

func New() *gin.Engine {
	router := gin.Default()

	// Set some shared context
	router.Use(func(c *gin.Context) {
		c.Set("db", &db.DB)
	})

	router.Use(sessions.Middleware)

	router.LoadHTMLGlob("views/**/*")

	router.GET("/.well-known/openid-configuration", oidc.Discovery)
	router.GET("/openid/jwks", oidc.Jwks)
	router.POST("/openid/registration", oidc.Registration)
	router.GET("/openid/registration/:id", oidc.GetClient)
	router.DELETE("/openid/registration/:id", oidc.DeleteClient)

	router.GET("/openid/authorization", oidc.CreateAuthorizationRequest)
	router.GET("/openid/authorization/:id", sessions.EnsureSignedIn, oidc.ShowAuthorizationRequest)
	router.POST("/openid/authorization/:id", oidc.ConfirmAuthorizationRequest)

	router.POST("/openid/token", oidc.Token)

	router.GET("/openid/userinfo", oidc.Userinfo)
	router.POST("/openid/userinfo", oidc.Userinfo)

	router.GET("/signin", sessions.New)
	router.POST("/signin", sessions.Create)

	usersNamespace := router.Group("/users", sessions.EnsureSignedIn)

	usersNamespace.GET("/", users.Index)
	usersNamespace.POST("/", users.Create)
	usersNamespace.GET("/:id", users.Show)
	usersNamespace.POST("/:id", users.Update)
	usersNamespace.DELETE("/:id", users.Delete)

	router.Static("/assets", "./assets")

	return router
}
