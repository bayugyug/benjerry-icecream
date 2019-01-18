package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/bayugyug/benjerry-icecream/models"
	"github.com/bayugyug/benjerry-icecream/utils"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

func (api *ApiHandler) CreateSourcing(w http.ResponseWriter, r *http.Request) {

	token := api.GetAuthToken(r)
	data := models.NewIcecream()
	if err := render.Bind(r, data); err != nil {
		utils.Dumper("BIND_FAILED", err)
		//206
		api.ReplyErrContent(w, r, http.StatusPartialContent, http.StatusText(http.StatusPartialContent))
		return
	}
	//get from url
	data.ID = strings.TrimSpace(chi.URLParam(r, "id"))
	//chk
	if !data.SanityCheck(data, "ADD-SOURCING") {
		utils.Dumper("MISSING_REQUIRED_PARAMS", data)
		//206
		api.ReplyErrContent(w, r, http.StatusPartialContent, http.StatusText(http.StatusPartialContent))
		return
	}

	if token == "" {
		utils.Dumper("INVALID_TOKEN:", token)
		//403
		api.ReplyErrContent(w, r, http.StatusForbidden, "Invalid token")
		return
	}
	//get author
	udata := models.NewUser()
	//sanity
	urow, err := udata.Get(ApiInstance.Context, ApiInstance.DB, token)
	if urow == nil || err != nil {
		utils.Dumper("RECORD_NOT_FOUND::USER", err)
		//404
		api.ReplyErrContent(w, r, http.StatusNotFound, "User not found")
		return
	}
	if urow.Status != "active" {
		utils.Dumper("RECORD_NOT_FOUND::USER_NOT_ACTIVE", urow.Status)
		//404
		api.ReplyErrContent(w, r, http.StatusNotFound, "User not found")
		return
	}
	//get the product_id in number notation
	prodID, _ := strconv.ParseInt(data.ID, 10, 64)

	gdata, gerr := data.Get(ApiInstance.Context, ApiInstance.DB, prodID)
	if gdata == nil || gerr != nil || gdata.ID == "" {
		utils.Dumper("RECORD_NOT_FOUND::PROD_ID", gdata.Name, prodID)
		//400
		api.ReplyErrContent(w, r, http.StatusNotFound, "Icecream not found")
		return
	}

	//free-up
	_ = data.DeleteSourcingValue(ApiInstance.Context, ApiInstance.DB, prodID)
	//add all
	for _, cup := range data.SourcingValues {
		_ = data.CreateSourcingValue(ApiInstance.Context, ApiInstance.DB, cup, prodID)
	}
	//response send
	render.JSON(w, r, APIResponse{
		Code:   http.StatusOK,
		Status: "Create Sourcing Values successful",
	})

}

func (api *ApiHandler) DeleteSourcing(w http.ResponseWriter, r *http.Request) {

	token := api.GetAuthToken(r)
	data := models.NewIcecream()
	//get from url
	data.ID = strings.TrimSpace(chi.URLParam(r, "id"))
	//chk
	if !data.SanityCheck(data, "DELETE-SOURCING") {
		utils.Dumper("MISSING_REQUIRED_PARAMS", data)
		//206
		api.ReplyErrContent(w, r, http.StatusPartialContent, http.StatusText(http.StatusPartialContent))
		return
	}

	if token == "" {
		utils.Dumper("INVALID_TOKEN:", token)
		//403
		api.ReplyErrContent(w, r, http.StatusForbidden, "Invalid token")
		return
	}
	//get author
	udata := models.NewUser()
	//sanity
	urow, err := udata.Get(ApiInstance.Context, ApiInstance.DB, token)
	if urow == nil || err != nil {
		utils.Dumper("RECORD_NOT_FOUND::USER", err)
		//404
		api.ReplyErrContent(w, r, http.StatusNotFound, "User not found")
		return
	}
	if urow.Status != "active" {
		utils.Dumper("RECORD_NOT_FOUND::USER_NOT_ACTIVE", urow.Status)
		//404
		api.ReplyErrContent(w, r, http.StatusNotFound, "User not found")
		return
	}
	//get the product_id in number notation
	prodID, _ := strconv.ParseInt(data.ID, 10, 64)

	gdata, gerr := data.Get(ApiInstance.Context, ApiInstance.DB, prodID)
	if gdata == nil || gerr != nil || gdata.ID == "" {
		utils.Dumper("RECORD_NOT_FOUND::PROD_ID", gdata.Name, prodID)
		//400
		api.ReplyErrContent(w, r, http.StatusNotFound, "Icecream not found")
		return
	}

	if len(gdata.SourcingValues) <= 0 {
		utils.Dumper("RECORD_NOT_FOUND::ALREADY_EMPTY", gdata.Name, prodID)
		//400
		api.ReplyErrContent(w, r, http.StatusNotFound, "Sourcing Values already deleted")
		return
	}

	//free-up
	_ = data.DeleteSourcingValue(ApiInstance.Context, ApiInstance.DB, prodID)
	//response send
	render.JSON(w, r, APIResponse{
		Code:   http.StatusOK,
		Status: "Delete Sourcing Values successful",
	})
}
