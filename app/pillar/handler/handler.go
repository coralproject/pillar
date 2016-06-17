package handler

import (
	"encoding/json"
	"github.com/coralproject/pillar/app/pillar/config"
	"github.com/coralproject/pillar/pkg/service"
	"github.com/coralproject/pillar/pkg/web"
	"net/http"
)

func doRespond(c *web.AppContext, object interface{}, appErr *web.AppError) {
	if appErr != nil {
		config.Logger().Printf("Call failed [%v]", appErr)
		payload, err := json.Marshal(appErr)
		if err != nil {
			http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
			c.SD.Client.Inc("Internal_Server_Error", 1, 1.0)
			return
		}

		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.WriteHeader(appErr.Code)
		c.Writer.Write(payload)
		c.SD.Client.Inc("App_Error", 1, 1.0)
		return
	}

	payload, err := json.Marshal(object)
	if err != nil {
		config.Logger().Printf("Call failed [%v]", err)
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		c.SD.Client.Inc("JSON_Error", 1, 1.0)
		return
	}

	//publish event before sending response
	service.PublishEvent(c, object, nil)

	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Write(payload)
	c.SD.Client.Inc("HTTP_Success", 1, 1.0)
}
