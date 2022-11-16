package tests

import (
	"bytes"
	"encoding/json"
	"go-stac-api-postgres/responses"
	"io/ioutil"
	"log"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearchItems(t *testing.T) {
	jsonBody := []byte(`{"ids": ["S2B_1CCV_20181004_0_L2A"]}`)
	bodyReader := bytes.NewReader(jsonBody)

	resp, err := http.Post(
		"http://localhost:6002/search",
		"application/json",
		bodyReader,
	)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()

	assert.Equalf(t, 200, resp.StatusCode, "create item")

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var search_response responses.SearchResponse
	json.Unmarshal(body, &search_response)

	assert.Equalf(t, "item collection retrieved successfully", search_response.Message, "search ids")
	assert.Equalf(t, "FeatureCollection", search_response.Type, "search ids")
	assert.Equalf(t, 100, search_response.Context.Limit, "search ids")
	assert.Equalf(t, 1, search_response.Context.Returned, "search ids")
	assert.Equalf(t, "S2B_1CCV_20181004_0_L2A", search_response.Features[0].Id, "search ids")
}

func TestSearchCollections(t *testing.T) {
	jsonBody := []byte(`{"collections": ["sentinel-s2-l2a-cogs-test"]}`)
	bodyReader := bytes.NewReader(jsonBody)

	resp, err := http.Post(
		"http://localhost:6002/search",
		"application/json",
		bodyReader,
	)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()

	assert.Equalf(t, 200, resp.StatusCode, "create item")

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var search_response responses.SearchResponse
	json.Unmarshal(body, &search_response)

	assert.Equalf(t, "item collection retrieved successfully", search_response.Message, "search collections")
	assert.Equalf(t, "FeatureCollection", search_response.Type, "search collections")
	assert.Equalf(t, 100, search_response.Context.Limit, "search collections")
	assert.Equalf(t, 50, search_response.Context.Returned, "search collections")
	assert.Equalf(t, "sentinel-s2-l2a-cogs-test", search_response.Features[0].Collection, "search collections")
}

func TestSearchNoCollections(t *testing.T) {
	jsonBody := []byte(`{"collections": ["sentinel-s2-l2a-cogs-test-test"]}`)
	bodyReader := bytes.NewReader(jsonBody)

	resp, err := http.Post(
		"http://localhost:6002/search",
		"application/json",
		bodyReader,
	)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()

	assert.Equalf(t, 200, resp.StatusCode, "create item")

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var search_response responses.SearchResponse
	json.Unmarshal(body, &search_response)

	assert.Equalf(t, "item collection retrieved successfully", search_response.Message, "search collections")
	assert.Equalf(t, "FeatureCollection", search_response.Type, "search collections")
	assert.Equalf(t, 100, search_response.Context.Limit, "search collections")
	assert.Equalf(t, 0, search_response.Context.Returned, "search collections")
}

func TestSearchGeometry(t *testing.T) {
	jsonBody := []byte(`{
		"collections": ["sentinel-s2-l2a-cogs-test"],
		"geometry": {
			"type": "Polygon",
        	"coordinates": [[
				[170.8515625, -74.14512718337613],
				[178.35937499999999, -74.14512718337613],
				[178.35937499999999, -70.15296965617042],
				[170.8515625, -70.15296965617042],
				[170.8515625, -74.14512718337613]
          	]]
      	}
	}`)
	bodyReader := bytes.NewReader(jsonBody)

	resp, err := http.Post(
		"http://localhost:6002/search",
		"application/json",
		bodyReader,
	)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()

	assert.Equalf(t, 200, resp.StatusCode, "create item")

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var search_response responses.SearchResponse
	json.Unmarshal(body, &search_response)

	assert.Equalf(t, "item collection retrieved successfully", search_response.Message, "search geometry")
	assert.Equalf(t, "FeatureCollection", search_response.Type, "search geometry")
	assert.Equalf(t, 100, search_response.Context.Limit, "search geometry")
	assert.Equalf(t, 50, search_response.Context.Returned, "search geometry")
	assert.Equalf(t, "sentinel-s2-l2a-cogs-test", search_response.Features[0].Collection, "search collections")
}

func TestSearchNoGeometry(t *testing.T) {
	jsonBody := []byte(`{
		"collections": ["sentinel-s2-l2a-cogs-test"],
		"geometry": {
			"type": "Polygon",
        	"coordinates": [[
				[70.8515625, -74.14512718337613],
				[78.35937499999999, -74.14512718337613],
				[78.35937499999999, -70.15296965617042],
				[70.8515625, -70.15296965617042],
				[70.8515625, -74.14512718337613]
          	]]
      	}
	}`)
	bodyReader := bytes.NewReader(jsonBody)

	resp, err := http.Post(
		"http://localhost:6002/search",
		"application/json",
		bodyReader,
	)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()

	assert.Equalf(t, 200, resp.StatusCode, "create item")

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var search_response responses.SearchResponse
	json.Unmarshal(body, &search_response)

	assert.Equalf(t, "item collection retrieved successfully", search_response.Message, "search geometry")
	assert.Equalf(t, "FeatureCollection", search_response.Type, "search geometry")
	assert.Equalf(t, 100, search_response.Context.Limit, "search geometry")
	assert.Equalf(t, 0, search_response.Context.Returned, "search geometry")
}

func TestSearchPoint(t *testing.T) {
	jsonBody := []byte(`{
		"collections": ["sentinel-s2-l2a-cogs-test"],
		"geometry": {
			"type": "Point",
			"coordinates": [177.064544, -72.690647]
      	}
	}`)
	bodyReader := bytes.NewReader(jsonBody)

	resp, err := http.Post(
		"http://localhost:6002/search",
		"application/json",
		bodyReader,
	)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()

	assert.Equalf(t, 200, resp.StatusCode, "create item")

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var search_response responses.SearchResponse
	json.Unmarshal(body, &search_response)

	assert.Equalf(t, "item collection retrieved successfully", search_response.Message, "search geometry")
	assert.Equalf(t, "FeatureCollection", search_response.Type, "search geometry")
	assert.Equalf(t, 100, search_response.Context.Limit, "search geometry")
	assert.Equalf(t, 49, search_response.Context.Returned, "search geometry")
	assert.Equalf(t, "sentinel-s2-l2a-cogs-test", search_response.Features[0].Collection, "search collections")
}

func TestSearchNoPoint(t *testing.T) {
	jsonBody := []byte(`{
		"collections": ["sentinel-s2-l2a-cogs-test"],
		"geometry": {
			"type": "Point",
			"coordinates": [77.064544, -72.690647]
      	}
	}`)
	bodyReader := bytes.NewReader(jsonBody)

	resp, err := http.Post(
		"http://localhost:6002/search",
		"application/json",
		bodyReader,
	)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()

	assert.Equalf(t, 200, resp.StatusCode, "create item")

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var search_response responses.SearchResponse
	json.Unmarshal(body, &search_response)

	assert.Equalf(t, "item collection retrieved successfully", search_response.Message, "search geometry")
	assert.Equalf(t, "FeatureCollection", search_response.Type, "search geometry")
	assert.Equalf(t, 100, search_response.Context.Limit, "search geometry")
	assert.Equalf(t, 0, search_response.Context.Returned, "search geometry")
}

func TestSearchLine(t *testing.T) {
	jsonBody := []byte(`{
		"collections": ["sentinel-s2-l2a-cogs-test"],
		"geometry": {
			"type": "LineString",
                "coordinates": [
                    [
                        179.85156249999997,
                        -70.554563528593656
                    ],
                    [
                        171.101642, 
                        -75.690647
                    ]
                ]
      	}
	}`)
	bodyReader := bytes.NewReader(jsonBody)

	resp, err := http.Post(
		"http://localhost:6002/search",
		"application/json",
		bodyReader,
	)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()

	assert.Equalf(t, 200, resp.StatusCode, "create item")

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var search_response responses.SearchResponse
	json.Unmarshal(body, &search_response)

	assert.Equalf(t, "item collection retrieved successfully", search_response.Message, "search geometry")
	assert.Equalf(t, "FeatureCollection", search_response.Type, "search geometry")
	assert.Equalf(t, 100, search_response.Context.Limit, "search geometry")
	assert.Equalf(t, 49, search_response.Context.Returned, "search geometry")
	assert.Equalf(t, "sentinel-s2-l2a-cogs-test", search_response.Features[0].Collection, "search collections")
}

func TestSearchNoLine(t *testing.T) {
	jsonBody := []byte(`{
		"collections": ["sentinel-s2-l2a-cogs-test"],
		"geometry": {
			"type": "LineString",
                "coordinates": [
                    [
                        79.85156249999997,
                        -70.554563528593656
                    ],
                    [
                        71.101642, 
                        -75.690647
                    ]
                ]
      	}
	}`)
	bodyReader := bytes.NewReader(jsonBody)

	resp, err := http.Post(
		"http://localhost:6002/search",
		"application/json",
		bodyReader,
	)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}
	defer resp.Body.Close()

	assert.Equalf(t, 200, resp.StatusCode, "create item")

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var search_response responses.SearchResponse
	json.Unmarshal(body, &search_response)

	assert.Equalf(t, "item collection retrieved successfully", search_response.Message, "search geometry")
	assert.Equalf(t, "FeatureCollection", search_response.Type, "search geometry")
	assert.Equalf(t, 100, search_response.Context.Limit, "search geometry")
	assert.Equalf(t, 0, search_response.Context.Returned, "search geometry")
}
