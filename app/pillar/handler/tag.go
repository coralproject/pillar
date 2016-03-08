package handler

import (
	"github.com/coralproject/pillar/pkg/service"
	"net/http"
)

func GetTags(w http.ResponseWriter, r *http.Request, c *service.AppContext) {
	dbObject, err := service.GetTags(c)
	doRespond(w, dbObject, err)
}

func CreateUpdateTag(w http.ResponseWriter, r *http.Request, c *service.AppContext) {
	dbObject, err := service.CreateUpdateTag(c)
	doRespond(w, dbObject, err)
}

func DeleteTag(w http.ResponseWriter, r *http.Request, c *service.AppContext) {
	err := service.DeleteTag(c)
	doRespond(w, nil, err)
}

