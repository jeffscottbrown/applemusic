package auth

import (
	"github.com/gorilla/pat"
	"github.com/gorilla/sessions"
	"github.com/jeffscottbrown/applemusic/secrets"
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

func Configure() {
	slog.Debug("Configuring authentication providers")

	googleId, googleSecret, callbackUrl := oauthConfig()

	slog.Debug("Credential Info", "gid", googleId, "gse", googleSecret, "cbu", callbackUrl)

	goth.UseProviders(
		google.New(googleId, googleSecret, callbackUrl, "profile"),
	)
}

func oauthConfig() (string, string, string) {
	googleSecret := retrieveSecreatValue("GOOGLE_SECRET")
	callbackUrl := retrieveSecreatValue("GOOGLE_CALLBACK_URL")
	googleId := retrieveSecreatValue("GOOGLE_ID")
	return googleId, googleSecret, callbackUrl
}

func retrieveSecreatValue(secretName string) string {
	clientSecret, err := secrets.RetrieveSecret(secretName)
	if err != nil {
		slog.Error("Falling back to OS environment variable", "variable", secretName)
		clientSecret = os.Getenv(secretName)
	}
	return clientSecret
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
