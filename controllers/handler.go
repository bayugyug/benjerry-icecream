package controllers

import (
	"net/http"

	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
)

type APIResponse struct {
	Code   int
	Status string
	Result interface{} `json:",omitempty"`
}

type TokenResponse struct {
	Token string
}

type ApiHandler struct {
}

func (api *ApiHandler) Login(w http.ResponseWriter, r *http.Request) {

	var token string
	//token send
	render.JSON(w, r, TokenResponse{Token: token})
}

func (api *ApiHandler) Otp(w http.ResponseWriter, r *http.Request) {

	//response send
	render.JSON(w, r, APIResponse{
		Code:   http.StatusOK,
		Status: "OK::" + r.URL.Path,
	})

}

func (api *ApiHandler) Logout(w http.ResponseWriter, r *http.Request) {
	//NOTE: not yet implemented for token invalidation ;-)
	//reply
	render.JSON(w, r,
		map[string]string{
			"Greeting": "Bye!",
		})
}

func (api *ApiHandler) IndexPage(w http.ResponseWriter, r *http.Request) {

	//reply
	render.JSON(w, r,
		map[string]string{
			"Greeting": "Welcome!",
		})
}

func (api ApiHandler) GetAuthToken(r *http.Request) string {
	_, claims, _ := jwtauth.FromContext(r.Context())

	//try checking it
	if token, ok := claims["mobile"].(string); ok {
		return token
	}

	return ""
}

//ReplyErrContent send 204 msg
//
//  http.StatusNoContent
//  http.StatusText(http.StatusNoContent)
func (api ApiHandler) ReplyErrContent(w http.ResponseWriter, r *http.Request, code int, msg string) {
	render.JSON(w, r, APIResponse{
		Code:   code,
		Status: msg,
	})
}
