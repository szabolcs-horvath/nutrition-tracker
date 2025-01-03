package auth

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/szabolcs-horvath/nutrition-tracker/util"
	"net/http"
	"net/url"
)

const Prefix = "/auth"

func Routes() map[string]http.HandlerFunc {
	return map[string]http.HandlerFunc{
		"GET /login":    loginHandler,
		"GET /callback": callbackHandler,
		"GET /logout":   logoutHandler,
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	state, err := generateRandomState()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session, err := util.CookieStoreInstance.New(r, "auth-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	session.Values["state"] = state
	if err = session.Save(r, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, util.AuthenticatorInstance.AuthCodeURL(state), http.StatusTemporaryRedirect)
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	token, err := util.AuthenticatorInstance.Exchange(r.Context(), r.FormValue("code"))
	if err != nil {
		http.Error(w, "Failed to exchange an authorization code for a token.", http.StatusUnauthorized)
		return
	}
	idToken, err := util.AuthenticatorInstance.VerifyIDToken(r.Context(), token)
	if err != nil {
		http.Error(w, "Failed to verify ID Token.", http.StatusInternalServerError)
		return
	}
	var profile map[string]interface{}
	if err = idToken.Claims(&profile); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	session, err := util.CookieStoreInstance.Get(r, "auth-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	session.Values["access_token"] = token.AccessToken
	session.Values["profile"] = profile
	if err = session.Save(r, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/htmx", http.StatusTemporaryRedirect)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	logoutUrl, err := url.Parse("https://" + util.GetEnvSafe("AUTH0_DOMAIN") + "/v2/logout")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	returnTo, err := url.Parse(scheme + "://" + r.Host)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	parameters := url.Values{}
	parameters.Add("returnTo", returnTo.String())
	parameters.Add("client_id", util.GetEnvSafe("AUTH0_CLIENT_ID"))
	logoutUrl.RawQuery = parameters.Encode()

	session, err := util.CookieStoreInstance.Get(r, "auth-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	clear(session.Values)
	if err = session.Save(r, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, logoutUrl.String(), http.StatusTemporaryRedirect)
}

func generateRandomState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	state := base64.StdEncoding.EncodeToString(b)

	return state, nil
}
