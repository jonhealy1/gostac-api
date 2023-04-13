package controllers

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/jonhealy1/goapi-stac/pg-api/models"

	"github.com/spatial-go/geoos/geoencoding"
)

/**
 * bbox2polygon takes a bounding box represented as a slice of 4 floats and returns a
 * string representation of a polygon geometry in GeoJSON format.
 *
 * @param  bbox  a slice of 4 floats representing the bounding box in the order: minimum longitude, minimum latitude,
 *               maximum longitude, maximum latitude
 * @return a string representation of a polygon geometry in GeoJSON format
 */
func bbox2polygon(bbox []float64) string {
	var buffer bytes.Buffer
	buffer.WriteString(`{"type":"Polygon", "Coordinates":[[`)
	buffer.WriteString(fmt.Sprintf("[%f,%f],", bbox[0], bbox[1]))
	buffer.WriteString(fmt.Sprintf("[%f,%f],", bbox[2], bbox[1]))
	buffer.WriteString(fmt.Sprintf("[%f,%f],", bbox[2], bbox[3]))
	buffer.WriteString(fmt.Sprintf("[%f,%f],", bbox[0], bbox[3]))
	buffer.WriteString(fmt.Sprintf("[%f,%f]]]}", bbox[0], bbox[1]))
	return buffer.String()
}

/**
 * lineString takes an array of 2D points and returns a string representation of a line string geometry in GeoJSON format.
 *
 * @param  geom  an array of 2D points represented as a slice of 2 floats
 * @return a string representation of a line string geometry in GeoJSON format
 */
func lineString(geom [][2]float64) string {
	// Initialize a strings.Builder
	var b strings.Builder

	// Write the start of the LineString JSON string to the strings.Builder
	b.WriteString(`{"type":"LineString", "Coordinates":[`)

	// Loop through the array of 2D points in the input 'geom' parameter
	for i, point := range geom {
		// Write the longitude and latitude of each point to the strings.Builder
		fmt.Fprintf(&b, "[%f,%f]", point[0], point[1])

		// Check if this is the last point in the 'geom' array
		if i != len(geom)-1 {
			// If it's not the last point, add a comma to separate it from the next point
			b.WriteString(",")
		}
	}

	// Write the end of the LineString JSON string to the strings.Builder
	b.WriteString("]}")

	// Return the final geoString
	return b.String()
}

// polygonString returns a string representation of a Polygon geometry
// in GeoJSON format, based on a 3D array of floating-point numbers.
//
// @param geom a 3D array of floating-point numbers representing the Polygon geometry
// @return a string representation of the Polygon geometry in GeoJSON format
func polygonString(geom [][][2]float64) string {
	// polygonString returns a string representation of a Polygon geometry
	// in GeoJSON format, based on a 3D array of floating-point numbers.
	var builder strings.Builder
	builder.WriteString(`{"type":"Polygon", "Coordinates":[[`)
	for _, coord := range geom[0][:len(geom[0])-1] {
		builder.WriteString(fmt.Sprintf("[%f,%f],", coord[0], coord[1]))
	}
	coord := geom[0][len(geom[0])-1]
	builder.WriteString(fmt.Sprintf("[%f,%f]]]}", coord[0], coord[1]))
	return builder.String()
}

/**
 * pointString takes a pair of floating-point numbers and returns a
 * string representation of a Point geometry in GeoJSON format.
 *
 * @param  geom  a pair of floating-point numbers representing the longitude and latitude of the point
 * @return a string representation of a Point geometry in GeoJSON format
 */
func pointString(geom [2]float64) string {
	// pointString returns a string representation of a Point geometry
	// in GeoJSON format, based on a 2D array of floating-point numbers.
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf(`{"type":"Point", "Coordinates":[%f,%f]}`, geom[0], geom[1]))
	return builder.String()
}

// sQLString returns an SQL query string based on the values of the fields in the
// provided `search_map` struct. The returned string will include placeholder
// values for query parameters such as the geometry and item IDs or collections.
//
// If `search_map.Ids` is 1, the query will include a condition to match item IDs.
// If `search_map.Collections` is 1, the query will include a condition to match collections.
// If `search_map.Geometry` is 1, the query will include a condition to match geometries.
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

/**
 * toWKT takes a string in GeoJSON format and returns a string in Well-Known Text (WKT) format.
 * @param geoString a string in GeoJSON format
 * @return a string in Well-Known Text (WKT) format
 */
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

/**
 * fix3dBbox converts a bbox of length 6 or 4 to a bbox of length 4
 * @param  search  a struct containing a bbox field which is a slice of floats representing the bounding box
 * @return a slice of 4 floats representing the bounding box in the order: minimum longitude, minimum latitude, maximum longitude, maximum latitude
 */
func fix3dBbox(search models.Search) []float64 {
	var bbox []float64
	if len(search.Bbox) == 6 {
		bbox = append(bbox, search.Bbox[0], search.Bbox[1], search.Bbox[3], search.Bbox[4])
	} else if len(search.Bbox) == 4 {
		bbox = search.Bbox
	} else {
		log.Println("Error: bbox must be of length 4 or 6")
	}
	return bbox
}

// split returns a bool indicating whether the given rune is a delimiter in a string representation of a JSON object.
// It is used to split a string into tokens, based on the delimiter characters ':', ',', '{', '}', '[', ']', '"', and ' '.
func split(r rune) bool {
	return r == ':' || r == ',' || r == '{' || r == '}' || r == '[' || r == ']' || r == '"' || r == ' '
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

/**
 * sendKafkaMessage sends a message to a specified Kafka topic using the given producer
 * @param  producer  a pointer to an initialized Kafka producer
 * @param  topic  a string representing the Kafka topic to send the message to
 * @param  message  a string containing the message to send to the Kafka topic
 */
func sendKafkaMessage(producer *kafka.Producer, topic string, message string) {
	deliveryChan := make(chan kafka.Event)
	err := producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(message),
	}, deliveryChan)

	if err != nil {
		fmt.Printf("Failed to send message: %v\n", err)
		return
	}

	e := <-deliveryChan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		fmt.Printf("Failed to deliver message: %v\n", m.TopicPartition.Error)
	} else {
		fmt.Printf("Delivered message to %v\n", m.TopicPartition)
	}

	close(deliveryChan)
}

func initKafkaProducer() (*kafka.Producer, error) {
	kafkaConfig := &kafka.ConfigMap{"bootstrap.servers": "kafka:9092"}
	producer, err := kafka.NewProducer(kafkaConfig)
	if err != nil {
		return nil, err
	}
	return producer, nil
}

func isKafkaAvailable() bool {
	timeout := 5 * time.Second
	_, err := net.DialTimeout("tcp", "kafka:9092", timeout)
	if err != nil {
		return false
	}
	return true
}
