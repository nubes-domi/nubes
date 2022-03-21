package sessions

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"nubes/sum/db"
	"nubes/sum/utils"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/jwt"
)

func New(c *gin.Context) {
	c.HTML(http.StatusOK, "sessions/new", gin.H{
		"continue": c.Query("continue"),
	})
}

func Create(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	next := c.PostForm("continue")

	user, err := db.DB.Users().FindByCredentials(username, password)
	if err != nil {
		c.HTML(http.StatusOK, "sessions/new", gin.H{
			"error": "Invalid username or password",
		})
	} else {
		Start(c, user)
		if next != "" {
			c.Redirect(302, next)
		} else {
			c.Redirect(302, fmt.Sprintf("/users/%d", user.ID))
		}
	}
}

func Middleware(c *gin.Context) {
	database := utils.CtxMustGet[*db.Database](c, "db")

	// Expire old sessions in the database
	database.UserSessions().CleanupExpired()

	// Read the sessions cookie and validate JWTs
	activeSessions, invalid := retrieveSessions(c)

	// Get the currently active user
	currentSessionToken, _ := c.Cookie("current_session")

	// Save valid sessions context for later use
	c.Set("activeSessions", activeSessions)

	// If any JWT has expired or has been marked invalid
	if len(invalid) > 0 {
		// Rewrite the cookie with the good sessions
		updateSessionsCookie(c)

		// The session that just expired was the active one
		if utils.Contains(invalid, currentSessionToken) {
			c.SetCookie("current_session", "", -1, "", "", false, true)
			currentSessionToken = ""
		}
	}

	if currentSessionToken != "" {
		currentSession := activeSessions[currentSessionToken]

		// Token is valid, not expired and in DB. We WILL always find a user.
		// (minus db corruption, handled in the function)
		user, _ := utils.CtxMustGet[*db.Database](c, "db").Users().FindById(currentSession.UserID)

		c.Set("currentSession", &currentSession)
		c.Set("currentUser", user)
	}

	c.Next()
}

func EnsureSignedIn(c *gin.Context) {
	_, ok := utils.CtxGet[*db.User](c, "currentUser")
	if ok {
		c.Next()
	} else {
		c.Redirect(302, "/signin?continue="+url.QueryEscape(c.Request.RequestURI))
	}
}

func Start(c *gin.Context, user *db.User) {
	database := utils.CtxMustGet[*db.Database](c, "db")

	// Genearate and save the new session
	session := user.NewSession(c)
	database.UserSessions().Create(&session)

	// Prepare a token to be given as cookie
	token, err := jwt.NewBuilder().
		JwtID(session.ID).
		Subject(fmt.Sprintf("%d", user.ID)).
		Expiration(session.ExpiresAt).
		Audience([]string{"sessions"}).
		Build()
	if err != nil {
		log.Panicf("Could not build the JWT token: %v", err)
	}

	// Sign the token
	signed := utils.JwtSign(token, "RS256")

	// Remember the signed version in the context, for possible future use
	session.SignedToken = signed

	// Put the new session in the context
	sessions := utils.CtxMustGet[map[string]db.UserSession](c, "activeSessions")
	sessions[session.ID] = session
	c.Set("activeSessions", sessions)

	// Rewrite the cookie
	updateSessionsCookie(c)
	c.SetCookie("current_session", session.ID, 60*60*24*365*10, "", "", false, true)

	// Remember the sign in
	c.Set("currentSession", &session)
	c.Set("currentUser", user)
}

func Terminate(c *gin.Context, session *db.UserSession) {
	database := utils.CtxMustGet[*db.Database](c, "db")
	database.UserSessions().Delete(session)
}

func CurrentUser(c *gin.Context) *db.User {
	currentUser, _ := utils.CtxGet[*db.User](c, "currentUser")
	return currentUser
}

func IsAnyoneSignedIn(c *gin.Context) bool {
	sessions := utils.CtxMustGet[map[string]db.UserSession](c, "activeSessions")
	return len(sessions) > 0
}

func IsUserSignedIn(c *gin.Context, user *db.User) bool {
	sessions := utils.CtxMustGet[map[string]db.UserSession](c, "activeSessions")
	return utils.AnyMap(sessions, func(_ string, s db.UserSession) bool {
		return s.UserID == user.ID
	})
}

func retrieveSessions(c *gin.Context) (result map[string]db.UserSession, badSessions []string) {
	result = make(map[string]db.UserSession)
	sessionsCookie, err := c.Cookie("sessions")
	if err != nil {
		return result, []string{}
	}

	sessionTokens := strings.Split(sessionsCookie, "|")
	for _, sessionToken := range sessionTokens {
		token, err := utils.JwtVerify(sessionToken)
		if err == nil {
			sub, _ := strconv.Atoi(token.Subject())
			result[token.JwtID()] = db.UserSession{
				ID:          token.JwtID(),
				UserID:      uint(sub),
				SignedToken: sessionToken,
			}
		} else {
			badSessions = append(badSessions, sessionToken)
		}
	}

	return
}

func updateSessionsCookie(c *gin.Context) {
	sessions := utils.CtxMustGet[map[string]db.UserSession](c, "activeSessions")

	// SignedToken is filled by retrieveSessions
	tokens := utils.CollectMap(sessions, func(s db.UserSession) string {
		return s.SignedToken
	})

	c.SetCookie("sessions", strings.Join(tokens, "|"), 60*60*24*365*10, "", "", false, true)
}
