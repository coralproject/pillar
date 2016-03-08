package handler

import (
	"encoding/json"
	"github.com/coralproject/pillar/app/pillar/config"
	"github.com/coralproject/pillar/pkg/service"
	"net/http"
)

type AppHandlerFunc func(rw http.ResponseWriter, r *http.Request, c *service.AppContext)

func (h AppHandlerFunc) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	c := service.NewContext(r.Body)
	defer c.Close()
	h(rw, r, c)
}

func doRespond(w http.ResponseWriter, object interface{}, appErr *service.AppError) {
	if appErr != nil {
		config.Logger.Printf("Call failed [%v]", appErr)
		payload, err := json.Marshal(appErr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(appErr.Code)
		w.Write(payload)
		return
	}

	payload, err := json.Marshal(object)
	if err != nil {
		config.Logger.Printf("Call failed [%v]", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(payload)
}

func CreateIndex(w http.ResponseWriter, r *http.Request, c *service.AppContext) {
	err := service.CreateIndex(c)
	doRespond(w, nil, err)
}

func HandleUserAction(w http.ResponseWriter, r *http.Request, c *service.AppContext) {
	err := service.CreateUserAction(c)
	doRespond(w, nil, err)
}

