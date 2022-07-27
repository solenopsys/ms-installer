package main

import (
	"io"
	"k8s.io/apimachinery/pkg/util/json"
	"log"
	"net/http"
	"time"
)

var cache = make(map[string]*map[string][]HelmChart)

type HelmChart struct {
	Name        string    `json:"name"`
	Version     string    `json:"version"`
	Description string    `json:"description"`
	Icon        string    `json:"icon"`
	APIVersion  string    `json:"apiVersion"`
	Urls        []string  `json:"urls"`
	Created     time.Time `json:"created"`
	Digest      string    `json:"digest"`
}

var test = map[string][]HelmChart{}

func getHelmCharts(url string) *map[string][]HelmChart {
	reqUrl := url + "/api/charts"
	println(reqUrl)
	req, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	err = json.Unmarshal(bodyBytes, &test)

	if err != nil {
		log.Fatalln(err)
	}

	return &test
}

func cachedReq(refresh bool, repositories map[string]string) map[string]*map[string][]HelmChart {
	if refresh || len(cache) == 0 {
		for name, url := range repositories {
			cache[name] = getHelmCharts(url)
		}
	}

	return cache
}
