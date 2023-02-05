package controllers

import (
	"encoding/json"
	"fmt"
	"go-stac-api-postgres/database"
	"go-stac-api-postgres/models"
	"net/http"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// GetSearch godoc
// @Summary GET Search request
// @Description Search for STAC items via the Search endpoint
// @Tags Search
// @ID get-search
// @Accept  json
// @Produce  json
// @Param bbox1, bbox2, bbox3, bbox4 path float true "Bbox"
// @Router /search [get]
func GetSearch(c *fiber.Ctx) error {
	var items []models.Item
	var search models.Search
	var searchMap models.SearchMap
	geoString := ""

	bboxString := c.Query("bbox")
	collectionsString := c.Query("collections")
	limitString := c.Query("limit")
	geometryString := c.Query("geometry")

	geomType := ""
	line := [][2]float64{}
	point := models.GeoJSONPoint{}.Coordinates
	if geometryString != "" {
		geomSlice := strings.FieldsFunc(geometryString, split)
		fmt.Println(geomSlice)
		geomType = geomSlice[1]
		if geomType == "Point" {
			point[0], _ = strconv.ParseFloat(geomSlice[3], 32)
			point[1], _ = strconv.ParseFloat(geomSlice[4], 32)
			searchMap.Geometry = 1
		}
		if geomType == "LineString" {
			val1, _ := strconv.ParseFloat(geomSlice[3], 32)
			val2, _ := strconv.ParseFloat(geomSlice[4], 32)
			val3, _ := strconv.ParseFloat(geomSlice[5], 32)
			val4, _ := strconv.ParseFloat(geomSlice[6], 32)

			line = [][2]float64{
				{val1, val2},
				{val3, val4},
			}
			searchMap.Geometry = 1
		}
	}

	limit := 100
	if limitString != "" {
		limit, _ = strconv.Atoi(limitString)
	}
	fmt.Println("limit: ", limit)

	if bboxString != "" {
		searchMap.Geometry = 1
	}

	if collectionsString != "" {
		searchMap.Collections = 1
		collections := strings.Split(collectionsString, ",")
		for i := 0; i < len(collections); i++ {
			search.Collections = append(search.Collections, collections[i])
		}
	}

	searchString := sQLString(searchMap)

	searchString += fmt.Sprintf("LIMIT %d", limit)

	if searchMap.Geometry == 1 {
		if bboxString != "" {
			bbox := strings.Split(bboxString, ",")

			b1, _ := strconv.ParseFloat(bbox[0], 32)
			b2, _ := strconv.ParseFloat(bbox[1], 32)
			b3, _ := strconv.ParseFloat(bbox[2], 32)
			b4, _ := strconv.ParseFloat(bbox[3], 32)

			search.Bbox = append(search.Bbox, b1, b2, b3, b4)
			geoString = bbox2polygon(search.Bbox)
		} else if geometryString != "" {
			if geomType == "Point" {
				geoString = pointString(point)
			} else if geomType == "LineString" {
				geoString = lineString(line)
			}
		}

		encodedString := toWKT(geoString)

		if len(search.Collections) > 0 {
			database.DB.Db.Raw(searchString, encodedString, search.Collections, limit).Scan(&items)
		} else {
			database.DB.Db.Raw(searchString, encodedString, limit).Scan(&items)
		}
	} else if len(search.Collections) > 0 {
		database.DB.Db.Raw(searchString, search.Collections, limit).Scan(&items)
	}

	context := models.Context{
		Returned: len(items),
		Limit:    limit,
	}

	var stac_items []interface{}
	for _, a_item := range items {
		var itemMap map[string]interface{}
		json.Unmarshal([]byte(a_item.Data), &itemMap)
		stac_items = append(stac_items, itemMap)
	}

	c.Status(http.StatusOK).JSON(&fiber.Map{
		"context":  context,
		"type":     "FeatureCollection",
		"features": stac_items,
	})

	return nil
}

// PostSearch godoc
// @Summary POST Search request
// @Description Search for STAC items via the Search endpoint
// @Tags Search
// @ID post-search
// @Accept  json
// @Produce  json
// @Param search body models.Search true "Search body json"
// @Router /search [post]
func PostSearch(c *fiber.Ctx) error {
	var search models.Search
	var items []models.Item
	var searchMap models.SearchMap

	if err := c.BodyParser(&search); err != nil {
		return c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"Status":  http.StatusBadRequest,
			"Message": "error",
			"Data":    err.Error(),
		})
	}
	limit := 100
	if search.Limit > 0 {
		limit = search.Limit
	}

	bbox := fix3dBbox(search)

	if len(bbox) == 4 || search.Geometry.Type == "Point" ||
		search.Geometry.Type == "Polygon" || search.Geometry.Type == "LineString" {
		searchMap.Geometry = 1
	}
	if len(search.Collections) > 0 {
		searchMap.Collections = 1
	}
	if len(search.Ids) > 0 {
		searchMap.Ids = 1
	}
	searchString := sQLString(searchMap)

	if len(search.Sortby) > 0 {
		searchString = BuildSortString(searchString, search)
	}

	searchString += fmt.Sprintf(" LIMIT %d", limit)

	if searchMap.Geometry == 1 {
		geoString := ""
		if len(bbox) == 4 {
			geoString = bbox2polygon(bbox)
		} else if search.Geometry.Type == "Point" {
			geom := models.GeoJSONPoint{}.Coordinates
			json.Unmarshal(search.Geometry.Coordinates, &geom)
			geoString = pointString(geom)
		} else if search.Geometry.Type == "Polygon" {
			geom := models.GeoJSONPolygon{}.Coordinates
			json.Unmarshal(search.Geometry.Coordinates, &geom)
			geoString = polygonString(geom)
		} else if search.Geometry.Type == "LineString" {
			geom := models.GeoJSONLine{}.Coordinates
			json.Unmarshal(search.Geometry.Coordinates, &geom)
			geoString = lineString(geom)
		}

		encodedString := toWKT(geoString)

		if searchMap.Collections == 1 && searchMap.Ids == 1 {
			database.DB.Db.Raw(searchString, encodedString, search.Collections, search.Ids).Scan(&items)
		} else if searchMap.Ids == 1 {
			database.DB.Db.Raw(searchString, encodedString, search.Ids).Scan(&items)
		} else if searchMap.Collections == 1 {
			database.DB.Db.Raw(searchString, encodedString, search.Collections).Scan(&items)
		} else {
			database.DB.Db.Raw(searchString, encodedString, limit).Scan(&items)
		}
	} else if searchMap.Collections == 1 || searchMap.Ids == 1 {
		tx1 := database.DB.Db.Limit(limit)
		tx2 := database.DB.Db.Limit(limit)
		if len(search.Collections) > 0 {
			tx1 = database.DB.Db.Limit(limit).Where("collection IN ?", search.Collections)
			tx2 = tx1.Limit(limit)
			fmt.Println(tx1)
		}

		if len(search.Ids) > 0 {
			tx2 = tx1.Limit(limit).Where("id IN ?", search.Ids)
			fmt.Println(tx1)
		}

		err := tx2.Find(&items).Error

		if err != nil {
			c.Status(http.StatusBadRequest).JSON(
				&fiber.Map{"message": "could not get items"})
			return err
		}
	}

	context := models.Context{
		Returned: len(items),
		Limit:    limit,
	}

	var stac_items []interface{}
	for _, a_item := range items {
		var itemMap map[string]interface{}
		json.Unmarshal([]byte(a_item.Data), &itemMap)
		stac_items = append(stac_items, itemMap)
	}

	c.Status(http.StatusOK).JSON(&fiber.Map{
		"context":  context,
		"type":     "FeatureCollection",
		"features": stac_items,
	})

	return nil
}

// buildSearchString takes a pointer to a Search struct and returns a string that is the ORDER BY
// clause for a SQL query.
func BuildSortString(searchString string, search models.Search) string {
	var fieldStrings []string
	for _, sort := range search.Sortby {
		var field_string string
		if strings.ContainsRune(sort.Field, '.') {
			substrings := strings.Split(sort.Field, ".")
			field_string = fmt.Sprintf("'%s' -> '%s'", substrings[0], substrings[1])
		} else {
			field_string = fmt.Sprintf("'%s'", sort.Field)
		}
		fieldStrings = append(fieldStrings, fmt.Sprintf("data -> %s %s", field_string, sort.Direction))
	}
	searchString += " ORDER BY " + strings.Join(fieldStrings, ", ")
	return searchString
}
