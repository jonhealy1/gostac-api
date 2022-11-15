package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-stac-api-postgres/database"
	"go-stac-api-postgres/models"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/spatial-go/geoos/geoencoding"
)

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

	geoString := ""
	if search.Geometry.Type == "Point" || search.Geometry.Type == "Polygon" {
		if search.Geometry.Type == "Point" {
			geom := models.GeoJSONPoint{}.Coordinates
			json.Unmarshal(search.Geometry.Coordinates, &geom)
			geoString = fmt.Sprintf(`{"type":"Point", "Coordinates":[%f,%f]}`, geom[0], geom[1])
			fmt.Println(geoString)
		} else if search.Geometry.Type == "Polygon" {
			geom := models.GeoJSONPolygon{}.Coordinates
			json.Unmarshal(search.Geometry.Coordinates, &geom)
			geoString = fmt.Sprintf(`{"type":"Polygon", "Coordinates":[[`)
			for i := 0; i < len(geom[0])-1; i++ {
				geoString += fmt.Sprintf("[%f,", geom[0][i][0])
				geoString += fmt.Sprintf("%f],", geom[0][i][1])
			}
			geoString += fmt.Sprintf("[%f,", geom[0][len(geom[0])-1][0])
			geoString += fmt.Sprintf("%f]", geom[0][len(geom[0])-1][1])
			geoString += fmt.Sprintf("]]}")
		}

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
		if len(search.Collections) > 0 && len(search.Ids) > 0 {
			database.DB.Db.Raw(`
				SELECT * FROM items WHERE ST_Intersects(items.geometry, ST_GeomFromText(?, 4326)) 
				AND items.collection in ? AND items.id in ?`,
				buf.String(), search.Collections, search.Ids).Scan(&items)
		} else if len(search.Ids) > 0 {
			database.DB.Db.Raw(`
				SELECT * FROM items WHERE ST_Intersects(items.geometry, ST_GeomFromText(?, 4326)) 
				AND items.id in ?`,
				buf.String(), search.Ids).Scan(&items)
		} else if len(search.Collections) > 0 {
			database.DB.Db.Raw(`
				SELECT * FROM items WHERE ST_Intersects(items.geometry, ST_GeomFromText(?, 4326)) 
				AND items.collection in ?`,
				buf.String(), search.Collections).Scan(&items)
		} else {
			database.DB.Db.Raw(`
				SELECT * FROM items WHERE ST_Intersects(items.geometry, ST_GeomFromText(?, 4326))`,
				buf.String()).Scan(&items)
		}
	} else if len(search.Collections) > 0 || len(search.Ids) > 0 {
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
		stac_items = append(stac_items, a_item.Data)
	}

	c.Status(http.StatusOK).JSON(&fiber.Map{
		"message":  "item collection retrieved successfully",
		"context":  context,
		"type":     "FeatureCollection",
		"features": stac_items,
	})

	return nil
}
