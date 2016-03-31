package handler

import (
	"encoding/json"
	"github.com/coralproject/pillar/app/pillar/config"
	"github.com/coralproject/pillar/pkg/web"
	"net/http"
	"github.com/coralproject/pillar/pkg/service"
	"log"
)

func doRespond(c *web.AppContext, object interface{}, appErr *web.AppError) {
	if appErr != nil {
		config.Logger().Printf("Call failed [%v]", appErr)
		payload, err := json.Marshal(appErr)
		if err != nil {
			http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
			return
		}

		c.Writer.WriteHeader(appErr.Code)
		c.Writer.Write(payload)
		return
	}

	payload, err := json.Marshal(object)
	if err != nil {
		config.Logger().Printf("Call failed [%v]", err)
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	//Publish to MQ if there is one
	if c.MQ.IsValid() {
		publish(c, object)
	}

	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Write(payload)
}

func publish(c *web.AppContext, object interface{}) {

	payload := service.GetPayload(c, object)
	if payload == nil {
		log.Printf("MQ - nothing to send\n\n")
		return
	}

	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error sending message: %s\n\n", err)
		return
	}

	c.MQ.Publish(data)
}
