package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/bayugyug/benjerry-icecream/models"
	"github.com/bayugyug/benjerry-icecream/utils"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

func (api *ApiHandler) CreateIcecream(w http.ResponseWriter, r *http.Request) {
	token := api.GetAuthToken(r)
	data := models.NewIcecream()
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

	if token == "" {
		utils.Dumper("INVALID_TOKEN:", token)
		//403
		api.ReplyErrContent(w, r, http.StatusForbidden, "Invalid token")
		return
	}
	utils.Dumper("TOKEN:", token)
	//get author
	udata := models.NewUser()
	urow, err := udata.Get(ApiInstance.Context, ApiInstance.DB, token)
	//sanity
	if err != nil {
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

	data.CreatedBy = urow.ID
	//create
	prodID := data.Create(ApiInstance.Context, ApiInstance.DB, data)
	switch prodID {
	case 0, -1, -2:
		utils.Dumper("RECORD_CREATE_FAILED", data.Name, prodID)
		//400
		api.ReplyErrContent(w, r, http.StatusBadRequest, "Record create failed")
		return
	case -3:
		pdata, perr := data.GetByName(ApiInstance.Context, ApiInstance.DB, data.Name)
		if pdata == nil || perr != nil {
			utils.Dumper("RECORD_CREATE_FAILED::OLD_GET", data.Name, prodID, perr)
			//400
			api.ReplyErrContent(w, r, http.StatusBadRequest, "Record create failed")
			return
		}
		prodID, _ = strconv.ParseInt(pdata.ID, 10, 64)
	}
	if prodID <= 0 {
		utils.Dumper("RECORD_CREATE_FAILED::PROD_ID", data.Name, prodID)
		//400
		api.ReplyErrContent(w, r, http.StatusBadRequest, "Record create failed")
		return
	}
	//add all sourcing-values
	for _, source := range data.SourcingValues {
		_ = data.CreateSourcingValue(ApiInstance.Context, ApiInstance.DB, source, prodID)
	}
	//add all ingredients
	for _, cup := range data.Ingredients {
		_ = data.CreateIngredient(ApiInstance.Context, ApiInstance.DB, cup, prodID)
	}
	//response send
	render.JSON(w, r, IcereamResponse{
		Code:      http.StatusOK,
		Status:    "Create successful",
		ProductID: fmt.Sprintf("%d", prodID),
	})

}

func (api *ApiHandler) UpdateIcecream(w http.ResponseWriter, r *http.Request) {

	token := api.GetAuthToken(r)
	data := models.NewIcecream()
	//get from url
	data.ID = strings.TrimSpace(chi.URLParam(r, "id"))
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

	if token == "" {
		utils.Dumper("INVALID_TOKEN:", token)
		//403
		api.ReplyErrContent(w, r, http.StatusForbidden, "Invalid token")
		return
	}
	//get author
	udata := models.NewUser()
	urow, err := udata.Get(ApiInstance.Context, ApiInstance.DB, token)
	//sanity
	if err != nil {
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
	data.ModifiedBy = urow.ID
	//create
	if exists := data.Exists(ApiInstance.Context, ApiInstance.DB, data); exists <= 0 {
		utils.Dumper("RECORD_UPDATE_FAILED::PROD_ID_MISSING", data.Name, data.ID)
		//404
		api.ReplyErrContent(w, r, http.StatusNotFound, "Icecream not found")
		return
	}
	if oks, err := data.Update(ApiInstance.Context, ApiInstance.DB, data); !oks || err != nil {
		utils.Dumper("RECORD_UPDATE_FAILED::PROD_ID", data.Name)
		//400
		api.ReplyErrContent(w, r, http.StatusBadRequest, "Record create failed")
		return
	}
	//get the product_id in number notation
	prodID, _ := strconv.ParseInt(data.ID, 10, 64)
	//free-up
	_ = data.DeleteSourcingValue(ApiInstance.Context, ApiInstance.DB, prodID)
	//add all sourcing-values
	for _, source := range data.SourcingValues {
		_ = data.CreateSourcingValue(ApiInstance.Context, ApiInstance.DB, source, prodID)
	}
	//free-up
	_ = data.DeleteIngredient(ApiInstance.Context, ApiInstance.DB, prodID)
	//add all ingredients
	for _, cup := range data.Ingredients {
		_ = data.CreateIngredient(ApiInstance.Context, ApiInstance.DB, cup, prodID)
	}
	//response send
	render.JSON(w, r, APIResponse{
		Code:   http.StatusOK,
		Status: "Update successful",
	})

}

func (api *ApiHandler) GetIcecream(w http.ResponseWriter, r *http.Request) {
	token := api.GetAuthToken(r)
	data := models.NewIcecream()
	//get from url
	data.ID = strings.TrimSpace(chi.URLParam(r, "id"))
	//chk
	if !data.SanityCheck(data, "GET") {
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
	if gdata == nil || gerr != nil {
		utils.Dumper("RECORD_NOT_FOUND::PROD_ID", data.Name, prodID)
		//400
		api.ReplyErrContent(w, r, http.StatusNotFound, "Icecream not found")
		return
	}

	//response send
	render.JSON(w, r, IcereamGetResponse{
		Code:   http.StatusOK,
		Status: "Record found",
		Result: gdata,
	})

}

func (api *ApiHandler) DeleteIcecream(w http.ResponseWriter, r *http.Request) {

	token := api.GetAuthToken(r)
	data := models.NewIcecream()
	//get from url
	data.ID = strings.TrimSpace(chi.URLParam(r, "id"))
	//chk
	if !data.SanityCheck(data, "DELETE") {
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
	utils.Dumper("TOKEN:", token)
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
	utils.Dumper("PRODUCT_ID:", prodID)
	gdata, gerr := data.Get(ApiInstance.Context, ApiInstance.DB, prodID)
	if gerr != nil || gdata.ID == "" {
		utils.Dumper("RECORD_NOT_FOUND::PROD_ID", data.Name, prodID)
		//400
		api.ReplyErrContent(w, r, http.StatusNotFound, "Icecream not found")
		return
	}
	if gdata.Status == "deleted" {
		utils.Dumper("RECORD_NOT_FOUND::STATUS_DELETED", data.Name, prodID, gdata.Status)
		//404
		api.ReplyErrContent(w, r, http.StatusNotFound, "Icecream already deleted")
		return
	}

	data.ID = gdata.ID
	data.ModifiedBy = urow.ID
	if oks, err := data.Delete(ApiInstance.Context, ApiInstance.DB, data); !oks || err != nil {
		utils.Dumper("RECORD_DELETE_FAILED::PROD_ID", data.Name)
		//400
		api.ReplyErrContent(w, r, http.StatusBadRequest, "Record delete failed")
		return
	}

	//response send
	render.JSON(w, r, APIResponse{
		Code:   http.StatusOK,
		Status: "Delete successful",
	})
}
