package auth

import (
	"context"
	"github.com/gorilla/sessions"
	"github.com/jeffscottbrown/applemusic/secrets"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"log/slog"
	"net/http"
	"os"
)

func login(res http.ResponseWriter, req *http.Request) {
	if _, err := gothic.CompleteUserAuth(res, req); err == nil {
		res.Header().Set("Location", "/")
		res.WriteHeader(http.StatusTemporaryRedirect)
	} else {
		gothic.BeginAuthHandler(res, req)
	}
}

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

func ConfigureAuthorizationHandlers(router *http.ServeMux) {
	router.HandleFunc("/auth/{provider}/callback", providerAwareHandler(authCallback))
	router.HandleFunc("/logout/{provider}", providerAwareHandler(logout))
	router.HandleFunc("/auth/{provider}", providerAwareHandler(login))
}

// gothic tries a number of techniques to retrieve the provider
// from the URL but it does not use PathValue, which is how
// the standard library provides access to the value
// see https://github.com/markbates/goth/blob/260588e82ba14930ae070a80acadcf0f75348c05/gothic/gothic.go#L263
// this wrapper will add the provider to the context in a way that gothic can use
func providerAwareHandler(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		provider := r.PathValue("provider")
		r = r.WithContext(context.WithValue(context.Background(), "provider", provider))

		h(w, r)
	}
}
