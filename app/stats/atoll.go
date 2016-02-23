package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/coralproject/pillar/server/model"
)

type atollPayloadWrapper struct {
	Data []model.Comment `json:"data"`
}

type atollPayload struct {
	Comments []model.Comment `json:"comments"`
}

func get(cs []model.Comment) {
	url := "https://atoll_stg.coralproject.net/pipelines/comments/score"

	p := atollPayloadWrapper{cs}

	jsonStr, err := json.MarshalIndent(p, "", "    ")

	fmt.Printf("\n\n%s", jsonStr)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}
