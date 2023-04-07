package tests

import (
	"testing"

	"github.com/jonhealy1/goapi-stac/controllers"
	"github.com/jonhealy1/goapi-stac/models"
)

func TestBuildSearchString(t *testing.T) {
	tests := []struct {
		search   *models.Search
		expected string
	}{
		{
			&models.Search{
				Sortby: []models.Sort{
					{Field: "field1", Direction: "ASC"},
					{Field: "field2", Direction: "DESC"},
				},
			},
			" ORDER BY data -> 'field1' ASC, data -> 'field2' DESC",
		},
		{
			&models.Search{
				Sortby: []models.Sort{
					{Field: "field1.subfield1", Direction: "ASC"},
					{Field: "field2", Direction: "DESC"},
				},
			},
			" ORDER BY data -> 'field1' -> 'subfield1' ASC, data -> 'field2' DESC",
		},
		{
			&models.Search{
				Sortby: []models.Sort{
					{Field: "field1", Direction: "ASC"},
				},
			},
			" ORDER BY data -> 'field1' ASC",
		},
	}
	for _, test := range tests {
		result := controllers.BuildSortString("", *test.search)
		if result != test.expected {
			t.Errorf("Expected %q but got %q", test.expected, result)
		}
	}
}
