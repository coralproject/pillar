package handler

import (
	"encoding/json"
	"github.com/coralproject/pillar/app/pillar/config"
	"github.com/coralproject/pillar/pkg/web"
	"net/http"
)

func doRespond(w http.ResponseWriter, object interface{}, appErr *web.AppError) {
	if appErr != nil {
		config.Logger().Printf("Call failed [%v]", appErr)
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
		config.Logger().Printf("Call failed [%v]", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(payload)
}
