package responses

import "github.com/gofiber/fiber/v2"

type SearchResponse struct {
	Status   int        `json:"status"`
	Message  string     `json:"message"`
	Type     string     `json:"type"`
	Context  Context    `json:"context"`
	Features *fiber.Map `json:"features"`
}

type Context struct {
	Returned int `json:"returned"`
	Limit    int `json:"limit"`
}
