package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/ardanlabs/kit/log"
	"github.com/coralproject/pillar/server/model"

	//	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type papiResult struct {
	Type               string      `json:"data_type"`
	Id                 int         `json:"data_id"`
	Word_count         int         `json:"word_count"`
	Headline           string      `json:"headline"`
	Summary            string      `json:"summary"`
	Section            interface{} `json:"section"`
	Subsection         interface{} `json:"subsection"`
	Authors            interface{} `json:"authors"`
	PublicationDate    int         `json:"publication_date"`
	PublicationDateISO int         `json:"publication_date_iso"`
}

type papi struct {
	Status string              `json:"status"`
	Result model.AssetMetadata `json:"result"`
}

func getAssetMeta() {

	var as []model.Asset

	err := db.C("asset").Find(bson.M{"metadata": bson.M{"$exists": false}}).Sort("-_id").All(&as)
	if err != nil {
		log.Error("stats", "request", err, "Could not load assets")
	}

	for _, a := range as {

		url := strings.Join([]string{"https://cms-publishapi.prd.nytimes.com/v1/publish/scoop/", strings.Replace(a.URL, "http://", "", 1)}, "")

		//url = "https://cms-publishapi.prd.nytimes.com/v1/publish/scoop/www.nytimes.com/2016/01/21/business/international/us-stock-markets-dow-sp-global-indexes.html"

		fmt.Printf("\n\n\n%+v\n", url)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Error("stats", "request", err, "Could not create request")
			continue
		}

		req.SetBasicAuth("coral", "88swF71pan5s0")

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Error("stats", "request", err, "Could not get request")
			continue
		}

		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)

		var r papi

		err = json.Unmarshal([]byte(body), &r)

		a.Metadata = r.Result

		db.C("asset").Update(bson.M{"_id": a.ID}, bson.M{"$set": bson.M{"metadata": a.Metadata}})

		fmt.Printf("%v\n%#v\n\n", a.ID, a.Metadata)
	}

	//	var meta map[string]interface{}
	/*
	 */
}
