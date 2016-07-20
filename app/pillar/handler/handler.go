package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/coralproject/pillar/pkg/service"
	"github.com/coralproject/pillar/pkg/web"
)

func doRespond(c *web.AppContext, object interface{}, appErr *web.AppError) {
	if appErr != nil {
		log.Printf("Error doing the web request. Error: %v", appErr)

		payload, err := json.Marshal(appErr)
		if err != nil {
			log.Printf("Error marshalling the error message. Error: %v", err)
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
		log.Printf("Marshalling the object. Error %v", err)
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		c.SD.Client.Inc("JSON_Error", 1, 1.0)
		return
	}

	//publish event before sending response
	service.PublishEvent(c, object, nil)

	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Write(payload)
}
