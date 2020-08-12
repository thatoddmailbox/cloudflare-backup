package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type resultInfo struct {
	TotalPages int `json:"total_pages"`
	Count      int `json:"count"`
	TotalCount int `json:"total_count"`
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
}

type result struct {
	Success    bool       `json:"success"`
	Errors     []string   `json:"errors"`
	Messages   []string   `json:"messages"`
	ResultInfo resultInfo `json:"result_info"`
}

type dnsRecord struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Name      string `json:"name"`
	Content   string `json:"content"`
	Proxiable bool   `json:"proxiable"`
	Proxied   bool   `json:"proxied"`
	TTL       uint64 `json:"ttl"`
	Locked    bool   `json:"locked"`
}

type zone struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	ModifiedOn  string `json:"modified_on"`
	ActivatedOn string `json:"activated_on"`
	CreatedOn   string `json:"created_on"`
}

type dnsRecordsResult struct {
	result
	DNSRecords []dnsRecord `json:"result"`
}

type zonesResult struct {
	result
	Zones []zone `json:"result"`
}

const baseURL = "https://api.cloudflare.com/client/v4/"

var apiKey string

func get(path string, params url.Values, output interface{}) error {
	request, err := http.NewRequest("GET", baseURL+path+"?"+params.Encode(), nil)
	if err != nil {
		return err
	}
	request.Header.Set("Authorization", "Bearer "+apiKey)
	request.Header.Set("Content-Type", "application/json")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, output)
}

func main() {
	log.Println("cloudflare-backup")

	flag.StringVar(&apiKey, "api-key", "", "The CloudFlare API key to use.")
	flag.Parse()

	result := zonesResult{}
	err := get("zones", url.Values{
		"per_page": []string{"50"},
	}, &result)
	if err != nil {
		panic(err)
	}

	if result.ResultInfo.Count != result.ResultInfo.TotalCount {
		// TODO: implement pagination so that this doesn't happen
		log.Fatalln("This program currently does not support accounts with more than 50 zones.")
	}

	for _, zone := range result.Zones {
		log.Printf("Processing %s...", zone.Name)

		// fetch the records for this zone
		dnsResult := dnsRecordsResult{}
		err := get("zones/"+zone.ID+"/dns_records", url.Values{
			"per_page": []string{"100"},
		}, &dnsResult)
		if err != nil {
			panic(err)
		}

		for _, record := range dnsResult.DNSRecords {
			log.Println(record)
		}
	}

	log.Println("Done!")
}
