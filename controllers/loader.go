package controllers

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/bayugyug/benjerry-icecream/models"
)

func (api *ApiHandler) Preload(fn string) (int, error) {

	body, err := ioutil.ReadFile(fn)
	if err != nil {
		log.Println("ERROR_READ_FILE::", err)
		return -1, err
	}

	var icecreams []models.Icecream
	if err := json.Unmarshal(body, &icecreams); err != nil {
		log.Println("ERROR_JSON_DECODE::", err)
		return -2, err
	}

	var all int
	log.Println("ALL:", len(icecreams))
	ice := models.NewIcecream()
	for _, v := range icecreams {
		pid := ice.Create(ApiInstance.Context, ApiInstance.DB, &v)

		//log.Println("new-product-id:", pid, v.Name)

		//add all sourcing-values
		for _, source := range v.SourcingValues {
			_ = ice.CreateSourcingValue(ApiInstance.Context, ApiInstance.DB, source, pid)
		}
		//add all ingredients
		for _, cup := range v.Ingredients {
			_ = ice.CreateIngredient(ApiInstance.Context, ApiInstance.DB, cup, pid)
		}
		all++
	}
	//good
	return all, nil
}
