package controllers

import "fmt"

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
