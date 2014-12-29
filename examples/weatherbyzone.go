package main

import (
	"fmt"
	"github.com/vascocosta/owm"
	"log"
)

func main() {
	// Create a new Client given an API key.
	c := owm.NewClient("YOUR_OPEN_WEATHER_MAP_API_KEY")
	// Decode the current weather of multiple locations given the zone
	// coordinates, map zoom and units. WeatherByZone returns a []Weather.
	w, err := c.WeatherByZone(12, 32, 15, 37, 10, "metric")
	// If there is an error, print it and terminate the program.
	if err != nil {
		log.Fatal(err)
	}
	// Print a string representation of w using the Stringer interface.
	for i := range w {
		fmt.Println(&w[i])
	}
}
