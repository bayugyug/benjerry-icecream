package controllers

import (
	"net/http"

	"github.com/go-chi/render"
)

func (api *ApiHandler) CreateIcecream(w http.ResponseWriter, r *http.Request) {

	//response send
	render.JSON(w, r, APIResponse{
		Code:   http.StatusOK,
		Status: "OK::" + r.URL.Path,
	})

}

func (api *ApiHandler) UpdateIcecream(w http.ResponseWriter, r *http.Request) {

	//response send
	render.JSON(w, r, APIResponse{
		Code:   http.StatusOK,
		Status: "OK::" + r.URL.Path,
	})

}

func (api *ApiHandler) GetIcecream(w http.ResponseWriter, r *http.Request) {

	//response send
	render.JSON(w, r, APIResponse{
		Code:   http.StatusOK,
		Status: "OK::" + r.URL.Path,
	})

}

func (api *ApiHandler) DeleteIcecream(w http.ResponseWriter, r *http.Request) {

	//response send
	render.JSON(w, r, APIResponse{
		Code:   http.StatusOK,
		Status: "OK::" + r.URL.Path,
	})

}
