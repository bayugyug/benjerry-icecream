package controllers

import (
	"net/http"

	"github.com/go-chi/render"
)

func (api *ApiHandler) CreateUser(w http.ResponseWriter, r *http.Request) {

	//response send
	render.JSON(w, r, APIResponse{
		Code:   http.StatusOK,
		Status: "OK::" + r.URL.Path,
	})

}

func (api *ApiHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {

	//response send
	render.JSON(w, r, APIResponse{
		Code:   http.StatusOK,
		Status: "OK::" + r.URL.Path,
	})

}

func (api *ApiHandler) GetUser(w http.ResponseWriter, r *http.Request) {

	//response send
	render.JSON(w, r, APIResponse{
		Code:   http.StatusOK,
		Status: "OK::" + r.URL.Path,
	})

}

func (api *ApiHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {

	//response send
	render.JSON(w, r, APIResponse{
		Code:   http.StatusOK,
		Status: "OK::" + r.URL.Path,
	})

}
