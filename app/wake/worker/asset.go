package worker

import (
	"encoding/json"
	"github.com/coralproject/pillar/pkg/model"
	"github.com/coralproject/pillar/pkg/web"
	"log"
	"os"
	"strconv"
	"time"
	"bytes"
)

type story struct {
	Creators  []creator `json:"creator"`
	Tracking  tracking  `json:"tracking"`
	Created   string    `json:"created_date"`
	Modified  string    `json:"last_modified"`
	Published string    `json:"published_date"`
}

type tracking struct {
	Section section `json:"section"`
}

type section struct {
	Section    string `json:"section"`
	Subsection string `json:"subsection"`
}
type creator struct {
	ID       string `json:"_id"`
	Name     string `json:"name"`
	Page     string `json:"bio_page"`
	Twitter  string `json:"twitter"`
	Facebook string `json:"facebook"`
}

func UpdateAsset(event model.Event) {

	var asset model.Asset
	if err := MnU(event.Payload, &asset); err != nil {
		log.Printf("Error [%v] unmarshalling Asset from Event [%v]", err, event)
	}

	if asset.URL == "" {
		log.Printf("Invalid Asset - Empty URL.\n")
		return
	}

	//If asset's fields are already populated, then return
	if !isIncomplete(asset) {
		log.Printf("The Asset looks Complete - No more updates needed!\n")
		return
	}

	s := getStory(asset.URL)
	if s == nil {
		return
	}

	//update authors, section, subsection
	doUpdateAsset(&asset, s)

	//Call Pillar to update this Asset
	data, _ := json.Marshal(asset)
	if _, err := web.Request(web.POST, os.Getenv("PILLAR_URL")+EPAsset,
		getHeader(), bytes.NewBuffer(data)); err != nil {
		log.Printf("Error updating Asset [%v]!\n", err)
	}
}

func isIncomplete(a model.Asset) bool {
	return a.DateCreated.IsZero() &&
		a.DateUpdated.IsZero() &&
		a.DatePublished.IsZero() &&
		(a.Authors == nil || len(a.Authors) == 0) &&
		a.Section == "" &&
		a.Subsection == ""
}

func getStory(url string) *story {

	saURL := os.Getenv("SA_URL") + url
	response, err := web.Request(web.GET, saURL, nil, nil)
	if err != nil || response.StatusCode != 200 {
		log.Printf("Bad response from Source!\n")
		return nil
	}

	var s story
	if err := json.Unmarshal([]byte(response.Body), &s); err != nil {
		log.Printf("Error transforming to a Story [%v]\n", err)
		return nil
	}

	return &s
}

func doUpdateAsset(asset *model.Asset, s *story) {
	//find authors
	creators := s.Creators
	authors := make([]model.Author, len(creators))
	for i := 0; i < len(creators); i++ {
		if creators[i].ID == "" {
			continue
		}

		authors[i].ID = creators[i].ID
		authors[i].URL = creators[i].Page
		authors[i].Name = creators[i].Name
		authors[i].Twitter = creators[i].Twitter
		authors[i].Facebook = creators[i].Facebook
	}

	asset.Authors = authors
	asset.Section = s.Tracking.Section.Section
	asset.Subsection = s.Tracking.Section.Subsection
	asset.DateCreated = parseTime(s.Created)
	asset.DateUpdated = parseTime(s.Modified)
	asset.DatePublished = parseTime(s.Published)
}

func getHeader() map[string]string {
	m := make(map[string]string)
	m["Content-Type"] = "application/json"

	return m
}

func parseTime(t string) time.Time {
	i64, err := strconv.ParseInt(t, 10, 64)
	if err != nil {
		var t time.Time
		return t
	}
	return time.Unix(i64, 0)
}
