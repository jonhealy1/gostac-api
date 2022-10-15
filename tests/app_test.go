package tests

// import (
// 	"io/ioutil"
// 	"net/http"
// 	"testing"

// 	"github.com/gofiber/fiber/v2"
// 	"github.com/gofiber/fiber/v2/middleware/cache"
// 	"github.com/gofiber/fiber/v2/middleware/compress"
// 	"github.com/gofiber/fiber/v2/middleware/cors"
// 	"github.com/gofiber/fiber/v2/middleware/etag"
// 	"github.com/gofiber/fiber/v2/middleware/favicon"
// 	"github.com/gofiber/fiber/v2/middleware/recover"
// 	"github.com/stretchr/testify/assert"

// 	database "go-stac-api-postgres/database"
// 	router "go-stac-api-postgres/router"
// )

// func Setup() *fiber.App {
// 	database.ConnectDb()
// 	app := fiber.New()

// 	app.Use(cors.New())
// 	app.Use(compress.New())
// 	app.Use(cache.New())
// 	app.Use(etag.New())
// 	app.Use(favicon.New())
// 	app.Use(recover.New())

// 	router.CollectionRoute(app)
// 	router.ItemRoute(app)

// 	return app
// }

// func TestCollectionsRoute(t *testing.T) {
// 	// Define a structure for specifying input and output
// 	// data of a single test case. This structure is then used
// 	// to create a so called test map, which contains all test
// 	// cases, that should be run for testing this function
// 	LoadCollection()
// 	LoadItems()

// 	tests := []struct {
// 		description string

// 		// Test input
// 		route string

// 		// Expected output
// 		expectedError bool
// 		expectedCode  int
// 		expectedBody  string
// 	}{
// 		{
// 			description:   "root catalog route",
// 			route:         "/",
// 			expectedError: false,
// 			expectedCode:  200,
// 			expectedBody:  "{\"stac_version\":\"1.0.0\",\"id\":\"test-catalog\",\"title\":\"go-stac-api\",\"description\":\"test catalog for go-stac-api, please edit\",\"links\":[{\"rel\":\"self\",\"href\":\"/\",\"type\":\"application/json\",\"title\":\"root catalog\"},{\"rel\":\"children\",\"href\":\"/collections\",\"type\":\"application/json\",\"title\":\"stac child collections\"}]}",
// 		},
// 		{
// 			description:   "all collections route",
// 			route:         "/collections/sentinel-s2-l2a-cogs-testing",
// 			expectedError: false,
// 			expectedCode:  200,
// 			expectedBody:  "{\"id\":\"sentinel-s2-l2a-cogs-testing\",\"stac_version\":\"1.0.0\"}",
// 		},
// 	}

// 	// Setup the app as it is done in the main function
// 	app := Setup()

// 	// Iterate through test single test cases
// 	for _, test := range tests {
// 		// Create a new http request with the route
// 		// from the test case
// 		req, _ := http.NewRequest(
// 			"GET",
// 			test.route,
// 			nil,
// 		)

// 		// Perform the request plain with the app.
// 		// The -1 disables request latency.
// 		res, err := app.Test(req, -1)

// 		// // verify that no error occured, that is not expected
// 		// assert.Equalf(t, test.expectedError, err != nil, test.description)

// 		// As expected errors lead to broken responses, the next
// 		// test case needs to be processed
// 		if test.expectedError {
// 			continue
// 		}

// 		// Verify if the status code is as expected
// 		assert.Equalf(t, test.expectedCode, res.StatusCode, test.description)

// 		// Read the response body
// 		body, err := ioutil.ReadAll(res.Body)

// 		// Reading the response body should work everytime, such that
// 		// the err variable should be nil
// 		assert.Nilf(t, err, test.description)

// 		// Verify, that the reponse body equals the expected body
// 		assert.Equalf(t, test.expectedBody, string(body), test.description)
// 	}
// }
