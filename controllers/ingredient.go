package controllers

import (
	"net/http"

	"github.com/go-chi/render"
)

func (api *ApiHandler) CreateIngredient(w http.ResponseWriter, r *http.Request) {

	//response send
	render.JSON(w, r, APIResponse{
		Code:   http.StatusOK,
		Status: "OK::" + r.URL.Path,
	})

}

func (api *ApiHandler) UpdateIngredient(w http.ResponseWriter, r *http.Request) {

	//response send
	render.JSON(w, r, APIResponse{
		Code:   http.StatusOK,
		Status: "OK::" + r.URL.Path,
	})

}

func (api *ApiHandler) GetIngredient(w http.ResponseWriter, r *http.Request) {

	//response send
	render.JSON(w, r, APIResponse{
		Code:   http.StatusOK,
		Status: "OK::" + r.URL.Path,
	})

}

func (api *ApiHandler) DeleteIngredient(w http.ResponseWriter, r *http.Request) {

	//response send
	render.JSON(w, r, APIResponse{
		Code:   http.StatusOK,
		Status: "OK::" + r.URL.Path,
	})

}
