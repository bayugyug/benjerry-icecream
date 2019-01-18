package controllers

import (
	"log"
	"net/http"

	"github.com/bayugyug/benjerry-icecream/models"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
)

type APIResponse struct {
	Code   int
	Status string
	Result interface{} `json:",omitempty"`
}

type TokenResponse struct {
	Code   int
	Status string
	Token  string
}

type OtpResponse struct {
	Code      int
	Status    string
	Otp       string
	OtpExpiry string
}

type IcereamResponse struct {
	Code      int
	Status    string
	ProductID string
}

type IcereamGetResponse struct {
	Code   int
	Status string
	Result *models.Icecream
}

type ApiHandler struct {
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
	tok, claims, _ := jwtauth.FromContext(r.Context())

	//sanity check
	if tok == nil || !tok.Valid {
		log.Println("INVALID_TOKEN")
		return ""
	}
	//try checking it
	if token, ok := claims["user"].(string); ok {
		//validate maybe fr db?
		data := models.NewUser()
		row, err := data.GetByToken(ApiInstance.Context, ApiInstance.DB, tok.Raw, token)
		if err != nil || row == nil {
			log.Println("ERR_VALIDATE_TOKEN::", err, token)
			return ""
		}
		if row.ExpiredToken > 0 {
			log.Println("ERR_VALIDATE_TOKEN::ALREADY_EXPIRED")
			return ""
		}
		if row.Token != tok.Raw {
			log.Println("ERR_VALIDATE_TOKEN::MISMATCH", tok.Raw)
			return ""
		}
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
