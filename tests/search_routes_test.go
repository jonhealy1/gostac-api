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

func TestSearch(t *testing.T) {
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
}
