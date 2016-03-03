package handler

import (
	"github.com/coralproject/pillar/pkg/service"
	"net/http"
	"github.com/coralproject/pillar/pkg/model"
	"encoding/json"
)

func GetTags(w http.ResponseWriter, r *http.Request) {
	dbObject, err := service.GetTags(GetAppContext(r, nil))
	doRespond(w, dbObject, err)
}

func CreateUpdateTag(w http.ResponseWriter, r *http.Request) {
	var input model.Tag
	json.NewDecoder(r.Body).Decode(&input)
	dbObject, err := service.CreateUpdateTag(GetAppContext(r, input))
	doRespond(w, dbObject, err)
}

func DeleteTag(w http.ResponseWriter, r *http.Request) {
	var input model.Tag
	json.NewDecoder(r.Body).Decode(&input)
	err := service.DeleteTag(GetAppContext(r, input))
	doRespond(w, nil, err)
}

