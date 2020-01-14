package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	log.Print("Looking for an OFFERZEN_TOKEN.")
	token := os.Getenv("OFFERZEN_TOKEN")
	if token == "" {
		panic("Please set the OFFERZEN_TOKEN.")
	}

	oc := &offerzenClient{
		http:    &http.Client{},
		baseURL: "https://www.offerzen.com/api",
		token:   token,
	}

	log.Print("Getting offerzen company metadata.")
	meta, err := oc.getPublicProfiles(1)
	if err != nil {
		panic("Error with retrieving company list from offerzen.")
	}

	log.Print("Scraping offerzen company data.")
	var cl []company
	for i := 1; i <= meta.PaginationInfo.TotalPages; i++ {
		log.Println("Retrieving page ", i)
		ocl, err := oc.getPublicProfiles(i)
		if err != nil {
			panic("Error with retrieving company list from offerzen.")
		}

		cl = append(cl, ocl.Result...)
		time.Sleep(1 * time.Second) // Don't request the offerzen api.
	}

	o := struct {
		Companies []company
	}{cl}

	b, err := json.MarshalIndent(&o, "", " ")
	if err != nil {
		panic("Error with marshalling offerzen data as json.")
	}

	err = ioutil.WriteFile("offerzendata.json", b, 0644)
	if err != nil {
		panic("Error with writing offerzendata.json.")
	}
}

// offerzenClient is the http client for accessing the offerzen public api.
type offerzenClient struct {
	http    *http.Client
	baseURL string
	token   string
}

func (oc offerzenClient) getPublicProfiles(page int) (*publicProfileResp, error) {
	url := fmt.Sprintf("%s/company/public_profiles?page=%d", oc.baseURL, page)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("authorization", oc.token)
	res, err := oc.http.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, errors.New("invalid request to offerzen api")
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var ppr publicProfileResp
	err = json.Unmarshal(b, &ppr)
	if err != nil {
		return nil, err
	}

	return &ppr, nil
}

// publicProfileResp is a response when using the api/company/public_profiles endpoint.
// Generated with https://mholt.github.io/json-to-go/
type publicProfileResp struct {
	Success        bool      `json:"success"`
	Result         []company `json:"result"`
	PaginationInfo struct {
		CurrentPage  int `json:"current_page"`
		TotalPages   int `json:"total_pages"`
		ItemsPerPage int `json:"items_per_page"`
		TotalItems   int `json:"total_items"`
	} `json:"pagination_info"`
}

// company is an offerzen representation of a company.
type company struct {
	ID                string        `json:"id"`
	Name              string        `json:"name"`
	LogoURL           string        `json:"logo_url"`
	BrandColor        string        `json:"brand_color"`
	ElevatorPitch     string        `json:"elevator_pitch"`
	Address           string        `json:"address"`
	NumberOfEmployees string        `json:"number_of_employees"`
	URL               string        `json:"url"`
	WebsiteURL        string        `json:"website_url"`
	FullWebsiteURL    string        `json:"full_website_url"`
	Cities            []interface{} `json:"cities"`
	TechStack         []struct {
		ID    string `json:"id"`
		Title string `json:"title"`
	} `json:"tech_stack"`
	Perks []struct {
		ID    string `json:"id"`
		Title string `json:"title"`
		Icon  string `json:"icon"`
	} `json:"perks"`
}
