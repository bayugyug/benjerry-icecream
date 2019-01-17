package controllers

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/bayugyug/benjerry-icecream/models"
	"github.com/bayugyug/benjerry-icecream/utils"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-chi/jwtauth"
	"github.com/go-chi/render"
)

func (api *ApiHandler) CreateUser(w http.ResponseWriter, r *http.Request) {

	data := models.NewUser()
	if err := render.Bind(r, data); err != nil {
		utils.Dumper("BIND_FAILED", err)
		//206
		api.ReplyErrContent(w, r, http.StatusPartialContent, http.StatusText(http.StatusPartialContent))
		return
	}
	defer r.Body.Close()
	//chk
	if !data.SanityCheck(data, "ADD") {
		utils.Dumper("MISSING_REQUIRED_PARAMS", data)
		//206
		api.ReplyErrContent(w, r, http.StatusPartialContent, http.StatusText(http.StatusPartialContent))
		return
	}
	data.Pass = data.Hash
	//exists
	if old := data.Exists(ApiInstance.Context, ApiInstance.DB, data.User); old > 0 {
		utils.Dumper("RECORD_EXISTS", data.User)
		//409
		api.ReplyErrContent(w, r, http.StatusConflict, "Record already exists")
		return
	}
	//add pin
	data.Otp = fmt.Sprintf("%05d", rand.Intn(9999))
	data.OtpExp = time.Now().Local().Add(time.Minute * time.Duration(5)).Format("2006-01-02 15:04:05")

	//create
	if oks := data.Create(ApiInstance.Context, ApiInstance.DB, data); oks <= 0 {
		utils.Dumper("RECORD_CREATE_FAILED", data.User, oks)
		//400
		api.ReplyErrContent(w, r, http.StatusBadRequest, "Record create failed")
		return
	}

	//response send
	render.JSON(w, r, OtpResponse{
		Code:      http.StatusOK,
		Status:    "Create successful",
		Otp:       data.Otp,
		OtpExpiry: data.OtpExp,
	})
}

func (api *ApiHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {

	token := api.GetAuthToken(r)
	data := models.NewUser()
	if err := render.Bind(r, data); err != nil {
		utils.Dumper("BIND_FAILED", err)
		//206
		api.ReplyErrContent(w, r, http.StatusPartialContent, http.StatusText(http.StatusPartialContent))
		return
	}
	defer r.Body.Close()
	//chk
	if !data.SanityCheck(data, "UPDATE") {
		utils.Dumper("MISSING_REQUIRED_PARAMS", data)
		//206
		api.ReplyErrContent(w, r, http.StatusPartialContent, http.StatusText(http.StatusPartialContent))
		return
	}
	//token mismatched
	if data.User != token {
		utils.Dumper("INVALID_TOKEN:", token, data.User)
		//403
		api.ReplyErrContent(w, r, http.StatusForbidden, "Invalid token")
		return
	}
	row, err := data.Get(ApiInstance.Context, ApiInstance.DB, data.User)
	//sanity
	if err != nil {
		utils.Dumper("RECORD_NOT_FOUND", err)
		//404
		api.ReplyErrContent(w, r, http.StatusNotFound, "Record not found")
		return
	}

	if row.Status == "deleted" {
		utils.Dumper("INVALID_STATUS", row.Status)
		//206
		api.ReplyErrContent(w, r, http.StatusPartialContent, "Status is not active")
		return

	}
	//md5
	data.Pass = data.Hash
	//create
	if oks, err := data.Update(ApiInstance.Context, ApiInstance.DB, data); !oks || err != nil {
		utils.Dumper("RECORD_UPDATE_FAILED", data.User, err)
		//400
		api.ReplyErrContent(w, r, http.StatusBadRequest, "Record update failed")
		return
	}

	//reply
	render.JSON(w, r, APIResponse{
		Code:   http.StatusOK,
		Status: "Update successful",
	})

}

func (api *ApiHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	//response send
	//404
	api.ReplyErrContent(w, r, http.StatusNotFound, "Not yet supported")
	return
}

func (api *ApiHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {

	token := api.GetAuthToken(r)
	data := models.NewUser()
	if err := render.Bind(r, data); err != nil {
		utils.Dumper("BIND_FAILED", err)
		//206
		api.ReplyErrContent(w, r, http.StatusPartialContent, http.StatusText(http.StatusPartialContent))
		return
	}
	defer r.Body.Close()
	//chk
	if !data.SanityCheck(data, "DELETE") {
		utils.Dumper("MISSING_REQUIRED_PARAMS", data)
		//206
		api.ReplyErrContent(w, r, http.StatusPartialContent, http.StatusText(http.StatusPartialContent))
		return
	}
	//token mismatched
	if data.User != token {
		utils.Dumper("INVALID_TOKEN:", token, data.User)
		//403
		api.ReplyErrContent(w, r, http.StatusForbidden, "Invalid token")
		return
	}
	row, err := data.Get(ApiInstance.Context, ApiInstance.DB, data.User)
	//sanity
	if err != nil {
		utils.Dumper("RECORD_NOT_FOUND", err)
		//404
		api.ReplyErrContent(w, r, http.StatusNotFound, "Record not found")
		return
	}

	if row.Status != "active" {
		utils.Dumper("INVALID_STATUS", row.Status)
		//206
		api.ReplyErrContent(w, r, http.StatusPartialContent, "Status is not active")
		return

	}

	//delete
	if oks, err := data.Delete(ApiInstance.Context, ApiInstance.DB, data.User); !oks || err != nil {
		utils.Dumper("RECORD_DELETE_FAILED", data.User, err)
		//400
		api.ReplyErrContent(w, r, http.StatusBadRequest, "Record delete failed")
		return
	}

	//reply
	render.JSON(w, r, APIResponse{
		Code:   http.StatusOK,
		Status: "Delete successful",
	})

}

func (api *ApiHandler) Login(w http.ResponseWriter, r *http.Request) {

	data := models.NewUser()
	if err := render.Bind(r, data); err != nil {
		utils.Dumper("BIND_FAILED", err)
		//206
		api.ReplyErrContent(w, r, http.StatusPartialContent, http.StatusText(http.StatusPartialContent))
		return
	}
	defer r.Body.Close()
	//chk
	if !data.SanityCheck(data, "LOG") {
		utils.Dumper("MISSING_REQUIRED_PARAMS", data)
		//206
		api.ReplyErrContent(w, r, http.StatusPartialContent, http.StatusText(http.StatusPartialContent))
		return
	}
	row, err := data.Get(ApiInstance.Context, ApiInstance.DB, data.User)
	//sanity
	if err != nil {
		utils.Dumper("RECORD_NOT_FOUND", err)
		//404
		api.ReplyErrContent(w, r, http.StatusNotFound, "Record not found")
		return
	}

	if row.Status != "active" {
		utils.Dumper("LOGIN_ACCOUNT_NOT_ACTIVE", row.Status)
		//403
		api.ReplyErrContent(w, r, http.StatusForbidden, "Account is not active")
		return
	}
	//good then check password match
	if data.Hash != row.Pass {
		utils.Dumper("LOGIN_PASSWORD_MISMATCH")
		//403
		api.ReplyErrContent(w, r, http.StatusForbidden, "Password mismatch or invalid")
		return
	}
	var token string
	//generate new token
	token, err = ApiInstance.AppJwtToken.GenToken(
		jwt.MapClaims{
			"user": data.User,
			"uuid": data.ID,
			"exp":  jwtauth.ExpireIn(30 * 24 * time.Hour),
		},
	)
	if err != nil {
		utils.Dumper("ERROR_TOKEN", err)
		//500
		api.ReplyErrContent(w, r, http.StatusInternalServerError, "Failed to generate token")
		return
	}
	//add pin
	data.Token = token
	data.TokenExp = time.Now().Local().Add(time.Hour * time.Duration(24*30)).Format("2006-01-02 15:04:05")

	//set flag
	_, _ = data.SetUserLogStatus(ApiInstance.Context, ApiInstance.DB, data)

	//token send
	render.JSON(w, r, TokenResponse{
		Code:   http.StatusOK,
		Status: "Login Successfull",
		Token:  token,
	})

}
