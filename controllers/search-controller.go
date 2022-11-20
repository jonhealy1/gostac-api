package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-stac-api-postgres/database"
	"go-stac-api-postgres/models"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/spatial-go/geoos/geoencoding"
)

func sQLString(search_map models.SearchMap) string {
	if search_map.Ids == 0 && search_map.Collections == 1 && search_map.Geometry == 0 {
		return `SELECT * FROM items WHERE items.collection in ?`
	} else if search_map.Ids == 0 && search_map.Collections == 0 && search_map.Geometry == 1 {
		return `SELECT * FROM items WHERE ST_Intersects(items.geometry, ST_GeomFromText(?, 4326))`
	} else if search_map.Ids == 1 && search_map.Collections == 0 && search_map.Geometry == 1 {
		return `SELECT * FROM items WHERE ST_Intersects(items.geometry, ST_GeomFromText(?, 4326)) 
		AND items.id in ?`
	} else if search_map.Ids == 0 && search_map.Collections == 1 && search_map.Geometry == 1 {
		return `SELECT * FROM items WHERE ST_Intersects(items.geometry, ST_GeomFromText(?, 4326)) 
		AND items.collection in ?`
	} else if search_map.Ids == 1 && search_map.Collections == 1 && search_map.Geometry == 1 {
		return `SELECT * FROM items WHERE ST_Intersects(items.geometry, ST_GeomFromText(?, 4326)) 
		AND items.collection in ? AND items.id in ?`
	}
	return ""
}

func toWKT(geoString string) string {
	buf := new(bytes.Buffer)
	buf.Write([]byte(geoString))
	got, err := geoencoding.Read(buf, geoencoding.GeoJSON)
	if err != nil {
		log.Println(err)
	}
	err = geoencoding.Write(buf, got, geoencoding.WKT)
	if err != nil {
		log.Println(err)
	}
	return buf.String()
}

func fix3dBbox(search models.Search) []float64 {
	var bbox []float64
	if len(search.Bbox) == 6 {
		bbox = append(bbox, search.Bbox[0])
		bbox = append(bbox, search.Bbox[1])
		bbox = append(bbox, search.Bbox[3])
		bbox = append(bbox, search.Bbox[4])
	} else if len(search.Bbox) == 4 {
		bbox = append(bbox, search.Bbox[0])
		bbox = append(bbox, search.Bbox[1])
		bbox = append(bbox, search.Bbox[2])
		bbox = append(bbox, search.Bbox[3])
	}
	return bbox
}

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

	bboxString := c.Query("bbox")
	collectionsString := c.Query("collections")

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

	if bboxString != "" {
		bbox := strings.Split(bboxString, ",")

		b1, _ := strconv.ParseFloat(bbox[0], 32)
		b2, _ := strconv.ParseFloat(bbox[1], 32)
		b3, _ := strconv.ParseFloat(bbox[2], 32)
		b4, _ := strconv.ParseFloat(bbox[3], 32)

		search.Bbox = append(search.Bbox, b1, b2, b3, b4)

		geoString := bbox2polygon(search.Bbox)

		encodedString := toWKT(geoString)

		if len(search.Collections) > 0 {
			database.DB.Db.Raw(searchString, encodedString, search.Collections).Scan(&items)
		} else {
			database.DB.Db.Raw(searchString, encodedString).Scan(&items)
		}
	} else if len(search.Collections) > 0 {
		database.DB.Db.Raw(searchString, search.Collections).Scan(&items)
	}

	limit := 0
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
		"message":  "item collection retrieved successfully",
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
			database.DB.Db.Raw(searchString, encodedString).Scan(&items)
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
		"message":  "item collection retrieved successfully",
		"context":  context,
		"type":     "FeatureCollection",
		"features": stac_items,
	})

	return nil
}
