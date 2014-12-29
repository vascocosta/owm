package main

import (
	"fmt"
	"github.com/vascocosta/owm"
	"log"
)

func main() {
	// Create a new Client given an API key.
	c := owm.NewClient("YOUR_OPEN_WEATHER_MAP_API_KEY")
	// Decode the current forecast of a location given the city coordinates,
	// days and units. ForecastByCoord returns a Forecast.
	f, err := c.ForecastByCoord(38.72, -9.13, 16, "metric")
	// If there is an error, print it and terminate the program.
	if err != nil {
		log.Fatal(err)
	}
	// Print a string representation of f using the Stringer interface.
	fmt.Println(&f)
}
