package controllers

import (
	"bytes"
	"fmt"
	"go-stac-api-postgres/models"
	"log"

	"github.com/spatial-go/geoos/geoencoding"
)

func bbox2polygon(bbox []float64) string {
	geoString := ""
	geoString += fmt.Sprintf(`{"type":"Polygon", "Coordinates":[[`)
	geoString += fmt.Sprintf("[%f,", bbox[0])
	geoString += fmt.Sprintf("%f],", bbox[1])
	geoString += fmt.Sprintf("[%f,", bbox[2])
	geoString += fmt.Sprintf("%f],", bbox[1])
	geoString += fmt.Sprintf("[%f,", bbox[2])
	geoString += fmt.Sprintf("%f],", bbox[3])
	geoString += fmt.Sprintf("[%f,", bbox[0])
	geoString += fmt.Sprintf("%f],", bbox[3])
	geoString += fmt.Sprintf("[%f,", bbox[0])
	geoString += fmt.Sprintf("%f]", bbox[1])
	geoString += fmt.Sprintf("]]}")
	return geoString
}

func lineString(geom [][2]float64) string {
	geoString := ""
	geoString = fmt.Sprintf(`{"type":"LineString", "Coordinates":[`)
	geoString += fmt.Sprintf("[%f,", geom[0][0])
	geoString += fmt.Sprintf("%f],", geom[0][1])
	geoString += fmt.Sprintf("[%f,", geom[1][0])
	geoString += fmt.Sprintf("%f]", geom[1][1])
	geoString += fmt.Sprintf("]}")
	return geoString
}

func polygonString(geom [][][2]float64) string {
	geoString := ""
	geoString = fmt.Sprintf(`{"type":"Polygon", "Coordinates":[[`)
	for i := 0; i < len(geom[0])-1; i++ {
		geoString += fmt.Sprintf("[%f,", geom[0][i][0])
		geoString += fmt.Sprintf("%f],", geom[0][i][1])
	}
	geoString += fmt.Sprintf("[%f,", geom[0][len(geom[0])-1][0])
	geoString += fmt.Sprintf("%f]", geom[0][len(geom[0])-1][1])
	geoString += fmt.Sprintf("]]}")
	return geoString
}

func pointString(geom [2]float64) string {
	geoString := ""
	geoString = fmt.Sprintf(`{"type":"Point", "Coordinates":[%f,%f]}`, geom[0], geom[1])
	return geoString
}

func sQLString(search_map models.SearchMap) string {
	if search_map.Ids == 0 && search_map.Collections == 1 && search_map.Geometry == 0 {
		return `SELECT * FROM items WHERE items.collection in ? LIMIT ?`
	} else if search_map.Ids == 0 && search_map.Collections == 0 && search_map.Geometry == 1 {
		return `SELECT * FROM items WHERE ST_Intersects(items.geometry, ST_GeomFromText(?, 4326)) LIMIT ?`
	} else if search_map.Ids == 1 && search_map.Collections == 0 && search_map.Geometry == 1 {
		return `SELECT * FROM items WHERE ST_Intersects(items.geometry, ST_GeomFromText(?, 4326)) 
		AND items.id in ? LIMIT ?`
	} else if search_map.Ids == 0 && search_map.Collections == 1 && search_map.Geometry == 1 {
		return `SELECT * FROM items WHERE ST_Intersects(items.geometry, ST_GeomFromText(?, 4326)) 
		AND items.collection in ? LIMIT ?`
	} else if search_map.Ids == 1 && search_map.Collections == 1 && search_map.Geometry == 1 {
		return `SELECT * FROM items WHERE ST_Intersects(items.geometry, ST_GeomFromText(?, 4326)) 
		AND items.collection in ? AND items.id in ? LIMIT ?`
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
