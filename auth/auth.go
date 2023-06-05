package auth

import (
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

var UsersPasswords map[string][]byte

func VerifyUserPass(name, password string) bool {
	wantPass, hasUser := UsersPasswords[name]
	if !hasUser {
		return false
	}
	if cmperr := bcrypt.CompareHashAndPassword(wantPass, []byte(password)); cmperr == nil {
		return true
	}
	return false
}
func BasicAuth(w *http.ResponseWriter, req *http.Request) bool {
	user, pass, ok := req.BasicAuth()
	if !ok || !VerifyUserPass(user, pass) {
		(*w).Header().Set("WWW-Authenticate", `Basic realm="api"`)
		http.Error(*w, "Unauthorized", http.StatusUnauthorized)
		return false
	}
	return true
}
