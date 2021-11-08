package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type endpoint struct {
	SvcName string
	Ips     []string
}

var (
	svcList   = []string{""}              // names of all services
	endpoints = make(map[string][]string) // all endpoints for all services
)

func getEndpoints(svcName string) {
	req, err := http.NewRequest("GET", "http://epwatcher:62000/"+svcName, nil)
	if err != nil {
		log.Println("Error reading request. ", err)
	}

	req.Header.Set("Cache-Control", "no-cache")

	client := &http.Client{Timeout: time.Second * 10}

	resp, err := client.Do(req)
	if err != nil {
		log.Println("error getting response: ", err.Error())
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("error reading response: ", err.Error())
		return
	}

	var ep endpoint
	err = json.Unmarshal(body, &ep)
	if err != nil {
		log.Println("error json unmarshalling: ", err.Error())
		return
	}
	endpoints[ep.SvcName] = ep.Ips
}

func getAllEndpoints() {
	if len(svcList) > 0 {
		for _, svc := range svcList {
			getEndpoints(svc)
		}
	}
}

func runComm(done chan bool) {
	go func() {
		for {
			select {
			case <-time.Tick(time.Microsecond * 10):
				getAllEndpoints()
			case <-done:
				return
			}
		}
	}()
}
