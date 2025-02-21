package auth

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/jeffscottbrown/gogoogle/secrets"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

func login(c *gin.Context) {
	req := c.Request
	res := c.Writer
	if _, err := gothic.CompleteUserAuth(res, req); err == nil {
		http.Redirect(res, req, "/", http.StatusTemporaryRedirect)
	} else {
		gothic.BeginAuthHandler(res, req)
	}
}

func logout(c *gin.Context) {
	req := c.Request
	res := c.Writer
	gothic.Logout(res, req)
	slog.Info("User logged out")
	http.Redirect(res, req, "/", http.StatusTemporaryRedirect)
}

func authCallback(c *gin.Context) {
	req := c.Request
	res := c.Writer
	user, err := gothic.CompleteUserAuth(res, req)
	if err != nil {
		slog.Error("Error authenticating user", "error", err)
		return
	}

	slog.Info("User authenticated", "name", user.Name)

	gothic.StoreInSession("authenticatedUser", user.Name, req, res)

	http.Redirect(res, req, "/", http.StatusTemporaryRedirect)
}

func IsAuthenticated(req *http.Request) bool {
	_, err := gothic.GetFromSession("authenticatedUser", req)
	return err == nil
}

func init() {
	gothic.Store = sessions.NewCookieStore([]byte(uuid.NewString()))
	slog.Debug("Configuring authentication providers")

	googleId, googleSecret, callbackUrl := oauthConfig()

	goth.UseProviders(
		google.New(googleId, googleSecret, callbackUrl, "profile"),
	)
}

func oauthConfig() (string, string, string) {
	googleSecret := retrieveSecretValue("GOOGLE_SECRET")
	callbackUrl := retrieveSecretValue("GOOGLE_CALLBACK_URL")
	if callbackUrl == "" {
		callbackUrl = "http://localhost:8080/auth/google/callback"
	}
	googleId := retrieveSecretValue("GOOGLE_ID")
	return googleId, googleSecret, callbackUrl
}

func retrieveSecretValue(secretName string) string {
	clientSecret, err := secrets.RetrieveSecret(secretName)
	if err != nil {
		slog.Error("Falling back to OS environment variable", "variable", secretName)
		clientSecret = os.Getenv(secretName)
	}
	return clientSecret
}

func ConfigureAuthorizationHandlers(router *gin.Engine) {
	router.GET("/auth/:provider/callback", authCallback)
	router.GET("/logout/:provider", logout)
	router.GET("/login/:provider", login)
}
