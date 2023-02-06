package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-stac-api-postgres/models"
	"go-stac-api-postgres/responses"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateItem(t *testing.T) {
	var expected_item models.Item
	jsonFile, err := os.Open("setup_data/S2B_1CCV_20181004_0_L2A-test.json")

	if err != nil {
		fmt.Println(err)
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &expected_item)
	responseBody := bytes.NewBuffer(byteValue)

	// Setup the app as it is done in the main function
	app := Setup()

	// Create a new HTTP request
	req, _ := http.NewRequest("POST", "/collections/sentinel-s2-l2a-cogs-test/items", bytes.NewBuffer(responseBody.Bytes()))

	// Set the Content-Type header to indicate the type of data in the request body
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	assert.Equalf(t, 201, resp.StatusCode, "create item")

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var item_response responses.CollectionResponse
	json.Unmarshal(body, &item_response)

	assert.Equalf(t, "success", item_response.Message, "create item")
}

func TestCreateItemNoCollection(t *testing.T) {
	var expected_item models.Item
	jsonFile, err := os.Open("setup_data/S2B_1CCV_20181004_0_L2A-test")

	if err != nil {
		fmt.Println(err)
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &expected_item)
	responseBody := bytes.NewBuffer(byteValue)

	// Setup the app as it is done in the main function
	app := Setup()

	// Create a new HTTP request
	req, _ := http.NewRequest("POST", "/collections/sentinel-s2-l2a-cogs-test-x/items", bytes.NewBuffer(responseBody.Bytes()))

	// Set the Content-Type header to indicate the type of data in the request body
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	assert.Equalf(t, 400, resp.StatusCode, "create item")

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var collection_response responses.CollectionResponse
	json.Unmarshal(body, &collection_response)

	assert.Equalf(t, "Collection does not exist", collection_response.Message, "create item")
}

func TestGetItem(t *testing.T) {
	var expected_item models.Item
	jsonFile, _ := os.Open("setup_data/S2B_1CCV_20181004_0_L2A-test.json")

	byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal(byteValue, &expected_item)

	tests := []struct {
		description   string
		route         string
		expectedError bool
		expectedCode  int
		expectedBody  models.Item
	}{
		{
			description:   "GET item route",
			route:         "/collections/sentinel-s2-l2a-cogs-test/items/S2B_1CCV_20181004_0_L2A-test",
			expectedError: false,
			expectedCode:  200,
			expectedBody:  expected_item,
		},
	}

	// Setup the app as it is done in the main function
	app := Setup()

	// Iterate through test single test cases
	for _, test := range tests {
		// Create a new http request with the route
		// from the test case
		req, _ := http.NewRequest(
			"GET",
			test.route,
			nil,
		)

		// Perform the request plain with the app.
		// The -1 disables request latency.
		res, err := app.Test(req, -1)

		// // verify that no error occured, that is not expected
		assert.Equalf(t, test.expectedError, err != nil, test.description)

		// As expected errors lead to broken responses, the next
		// test case needs to be processed
		if test.expectedError {
			continue
		}

		// Verify if the status code is as expected
		assert.Equalf(t, test.expectedCode, res.StatusCode, test.description)

		// Read the response body
		body, err := ioutil.ReadAll(res.Body)
		assert.Nilf(t, err, "get item")

		var stac_item models.Item

		json.Unmarshal(body, &stac_item)

		// stac_item.Collection = "sentinel-s2-l2a-cogs-test"

		// Reading the response body should work everytime, such that
		// the err variable should be nil
		assert.Nilf(t, err, test.description)

		// Verify, that the reponse body equals the expected body
		assert.Equalf(t, test.expectedBody, stac_item, test.description)
	}
}

func TestGetItemCollection(t *testing.T) {
	tests := []struct {
		description   string
		route         string
		expectedError bool
		expectedCode  int
	}{
		{
			description:   "GET item collection route",
			route:         "/collections/sentinel-s2-l2a-cogs-test/items",
			expectedError: false,
			expectedCode:  200,
		},
	}

	// Setup the app as it is done in the main function
	app := Setup()

	// Iterate through test single test cases
	for _, test := range tests {
		req, _ := http.NewRequest(
			"GET",
			test.route,
			nil,
		)

		// Perform the request plain with the app.
		// The -1 disables request latency.
		res, err := app.Test(req, -1)

		// // verify that no error occured, that is not expected
		assert.Equalf(t, test.expectedError, err != nil, test.description)

		// As expected errors lead to broken responses, the next
		// test case needs to be processed
		if test.expectedError {
			continue
		}

		// Verify if the status code is as expected
		assert.Equalf(t, test.expectedCode, res.StatusCode, test.description)

		// Read the response body
		body, err := ioutil.ReadAll(res.Body)
		assert.Nilf(t, err, test.description)

		var item_collection models.ItemCollection

		json.Unmarshal(body, &item_collection)

		assert.GreaterOrEqual(t, item_collection.Context.Returned, 1, test.description)
	}
}

func TestEditItem(t *testing.T) {
	var expected_item models.Item
	jsonFile, err := os.Open("setup_data/S2B_1CCV_20181004_0_L2A-test-updated.json")

	if err != nil {
		fmt.Println(err)
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &expected_item)
	responseBody := bytes.NewBuffer(byteValue)

	// Setup the app as it is done in the main function
	app := Setup()

	// Create a new HTTP request
	req, _ := http.NewRequest(http.MethodPut, "/collections/sentinel-s2-l2a-cogs-test/items/S2B_1CCV_20181004_0_L2A-test", bytes.NewBuffer(responseBody.Bytes()))

	// Set the Content-Type header to indicate the type of data in the request body
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	assert.Equalf(t, "200 OK", resp.Status, "update item")

	// Read Response Body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	var item_response responses.CollectionResponse
	json.Unmarshal(body, &item_response)

	assert.Equalf(t, "success", item_response.Message, "update item")
}

func TestDeleteItem(t *testing.T) {
	app := Setup()
	resp, err := http.NewRequest(
		"DELETE",
		"/collections/sentinel-s2-l2a-cogs-test/items/S2B_1CCV_20181004_0_L2A-test",
		nil,
	)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}

	res, err := app.Test(resp, -1)

	assert.Equalf(t, 200, res.StatusCode, "delete item")
}
