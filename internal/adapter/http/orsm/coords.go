package osrm

import (
	"fmt"
	"strings"
)

type Coord struct {
	Latitude  float64
	Longitude float64
}

func pointsBuilder(coords []Coord) string {
	var pointsBuilder strings.Builder
	userLen := len(coords)
	for i, coord := range coords {
		//fmt.Printf("lat: %f | long: %f | address: %s\n", student.Latitude, student.Longitude, student.Address)
		fmt.Fprintf(&pointsBuilder, "%.5f,%.5f", coord.Longitude, coord.Latitude)
		if i < userLen-1 {
			pointsBuilder.WriteString(";")
		}

	}

	return pointsBuilder.String()
}
