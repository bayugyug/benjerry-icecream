package controllers

import (
	"net/http"

	"github.com/go-chi/render"
)

func (api *ApiHandler) CreateSourcing(w http.ResponseWriter, r *http.Request) {

	//response send
	render.JSON(w, r, APIResponse{
		Code:   http.StatusOK,
		Status: "OK::" + r.URL.Path,
	})

}

func (api *ApiHandler) UpdateSourcing(w http.ResponseWriter, r *http.Request) {

	//response send
	render.JSON(w, r, APIResponse{
		Code:   http.StatusOK,
		Status: "OK::" + r.URL.Path,
	})

}

func (api *ApiHandler) GetSourcing(w http.ResponseWriter, r *http.Request) {

	//response send
	render.JSON(w, r, APIResponse{
		Code:   http.StatusOK,
		Status: "OK::" + r.URL.Path,
	})

}

func (api *ApiHandler) DeleteSourcing(w http.ResponseWriter, r *http.Request) {

	//response send
	render.JSON(w, r, APIResponse{
		Code:   http.StatusOK,
		Status: "OK::" + r.URL.Path,
	})

}
