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
