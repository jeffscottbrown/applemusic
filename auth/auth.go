package auth

import (
	"github.com/gorilla/pat"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"log/slog"
	"net/http"
	"os"
)

func logout(res http.ResponseWriter, req *http.Request) {
	gothic.Logout(res, req)
	slog.Info("User logged out")
	res.Header().Set("Location", "/")
	res.WriteHeader(http.StatusTemporaryRedirect)
}

var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
var callbackUrl = getCallbackUrl()
var googleSecret = os.Getenv("GOOGLE_SECRET")
var googleId = os.Getenv("GOOGLE_ID")

func authCallback(res http.ResponseWriter, req *http.Request) {
	user, err := gothic.CompleteUserAuth(res, req)
	if err != nil {
		slog.Error("Error authenticating user", "error", err)
		return
	}

	slog.Info("User authenticated", "name", user.Name)

	gothic.StoreInSession("authenticatedUser", user.Name, req, res)

	res.Header().Set("Location", "/")
	res.WriteHeader(http.StatusTemporaryRedirect)
}

func getCallbackUrl() string {
	value, exists := os.LookupEnv("GOOGLE_CALLBACK_URL")
	if !exists {
		value = "http://localhost:8080/auth/google/callback"
		slog.Warn("GOOGLE_CALLBACK_URL not set, using default", "value", value)
	}
	return value
}

func Configure() {
	slog.Debug("Configuring authentication providers")
	goth.UseProviders(
		google.New(googleId, googleSecret, callbackUrl, "profile"),
	)
}

func ConfigureAuthorizationHandlers(router *pat.Router) {
	router.Get("/auth/{provider}/callback", authCallback)
	router.Get("/logout/{provider}", logout)

	router.Get("/auth/{provider}", func(res http.ResponseWriter, req *http.Request) {
		if _, err := gothic.CompleteUserAuth(res, req); err == nil {
			res.Header().Set("Location", "/")
			res.WriteHeader(http.StatusTemporaryRedirect)
		} else {
			gothic.BeginAuthHandler(res, req)
		}
	})
	router.Get("/logout/{provider}", func(res http.ResponseWriter, req *http.Request) {
		gothic.Logout(res, req)
		res.Header().Set("Location", "/")
		res.WriteHeader(http.StatusTemporaryRedirect)
	})
}
