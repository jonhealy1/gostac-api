package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/jonhealy1/goapi-stac/database"
	"github.com/jonhealy1/goapi-stac/models"
	"github.com/jonhealy1/goapi-stac/responses"
	routes "github.com/jonhealy1/goapi-stac/router"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/stretchr/testify/assert"
)

func EsSetup() *fiber.App {
	log.Println("Setting up Elasticsearch for testing")
	database.ConnectES()
	app := fiber.New()

	app.Use(cors.New())
	app.Use(compress.New())
	//app.Use(cache.New())
	app.Use(etag.New())
	app.Use(favicon.New())
	app.Use(recover.New())

	routes.ESCollectionRoute(app)
	routes.ESItemRoute(app)
	//routes.SearchRoute(app)

	return app
}

func TestEsCreateCollection(t *testing.T) {
	var expected_collection models.Collection
	jsonFile, err := os.Open("setup_data/collection.json")

	if err != nil {
		fmt.Println(err)
	}
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &expected_collection)
	responseBody := bytes.NewBuffer(byteValue)

	// Setup the app as it is done in the main function
	app := EsSetup()

	// Create a new HTTP request
	req, _ := http.NewRequest("POST", "/es/collections", bytes.NewBuffer(responseBody.Bytes()))

	// Set the Content-Type header to indicate the type of data in the request body
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req, -1)
	if err != nil {
		log.Fatalln(err)
	}

	assert.Equalf(t, 201, resp.StatusCode, "create collection")

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var collection_response responses.CollectionResponse
	json.Unmarshal(body, &collection_response)

	assert.Equalf(t, "success", collection_response.Message, "create collection")
}

func TestEsGetCollection(t *testing.T) {
	// Setup the app as it is done in the main function
	app := EsSetup()

	LoadEsCollection()
	// LoadEsItems()

	var expected_collection models.Collection
	jsonFile, _ := os.Open("setup_data/collection.json")

	byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal(byteValue, &expected_collection)

	tests := []struct {
		description   string
		route         string
		expectedError bool
		expectedCode  int
		expectedBody  models.Collection
	}{
		{
			description:   "GET collection route",
			route:         "/es/collections/sentinel-s2-l2a-cogs-test-2",
			expectedError: false,
			expectedCode:  200,
			expectedBody:  expected_collection,
		},
	}

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
		assert.Nilf(t, err, "Create collection")

		var stac_collection models.Collection

		json.Unmarshal(body, &stac_collection)

		// Reading the response body should work everytime, such that
		// the err variable should be nil
		assert.Nilf(t, err, test.description)

		// Verify, that the reponse body equals the expected body
		assert.Equalf(t, test.expectedBody, stac_collection, test.description)
	}
}

func TestEsGetAllCollections(t *testing.T) {
	tests := []struct {
		description   string
		route         string
		expectedError bool
		expectedCode  int
	}{
		{
			description:   "GET collections route",
			route:         "/es/collections",
			expectedError: false,
			expectedCode:  200,
		},
	}

	// Setup the app as it is done in the main function
	app := EsSetup()

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

		var stac_collection []models.Collection

		json.Unmarshal(body, &stac_collection)
	}
}

func TestEsEditCollection(t *testing.T) {
	// Setup the app as it is done in the main function
	app := EsSetup()

	// Assuming the collection ID is known, you can directly fetch the existing collection
	collectionId := "sentinel-s2-l2a-cogs-test-2"

	// Create a new HTTP request for getting the existing collection
	getReq, _ := http.NewRequest("GET", "/es/collections/"+collectionId, nil)

	getResp, err := app.Test(getReq, -1)
	if err != nil {
		log.Fatalln(err)
	}

	assert.Equalf(t, 200, getResp.StatusCode, "get existing collection")

	body, err := ioutil.ReadAll(getResp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var existing_collection models.StacCollection
	json.Unmarshal(body, &existing_collection)

	// Change the "stac_version" field of the existing_collection
	existing_collection.StacVersion = "1.0.0"

	// Marshal the updated existing_collection into JSON
	updatedCollectionJSON, _ := json.Marshal(existing_collection)
	updatedCollectionBuffer := bytes.NewBuffer(updatedCollectionJSON)

	// Create a new HTTP request for updating the collection
	updateReq, _ := http.NewRequest("PUT", "/es/collections/"+collectionId, updatedCollectionBuffer)
	updateReq.Header.Set("Content-Type", "application/json")

	updateResp, err := app.Test(updateReq, -1)
	if err != nil {
		log.Fatalln(err)
	}

	assert.Equalf(t, 200, updateResp.StatusCode, "update collection")

	// Create a new HTTP request for getting the updated collection
	getUpdatedReq, _ := http.NewRequest("GET", "/es/collections/"+collectionId, nil)

	getUpdatedResp, err := app.Test(getUpdatedReq, -1)
	if err != nil {
		log.Fatalln(err)
	}

	assert.Equalf(t, 200, getUpdatedResp.StatusCode, "get updated collection")

	updatedBody, err := ioutil.ReadAll(getUpdatedResp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var updated_collection_response models.StacCollection
	json.Unmarshal(updatedBody, &updated_collection_response)

	// Check if the "stac_version" field is updated
	assert.Equalf(t, "1.0.0", updated_collection_response.StacVersion, "check updated stac_version")
}

func TestEsDeleteCollection(t *testing.T) {
	app := EsSetup()

	// Create Request
	req, err := http.NewRequest("DELETE", "/es/collections/sentinel-s2-l2a-cogs-test-2", nil)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}

	// Fetch Request
	resp, err := app.Test(req, -1)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	assert.Equalf(t, "200 OK", resp.Status, "delete collection")

	// Read Response Body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	var collection_response responses.CollectionResponse
	json.Unmarshal(body, &collection_response)

	assert.Equalf(t, "success", collection_response.Message, "delete collection")
}
