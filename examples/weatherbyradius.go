package main

import (
	"fmt"
	"github.com/vascocosta/owm"
	"log"
)

func main() {
	// Create a new Client given an API key.
	c := owm.NewClient("YOUR_OPEN_WEATHER_MAP_API_KEY")
	// Decode the current weather of multiple locations given the center
	// coordinates, radius and units. WeatherByRadius returns a []Weather.
	w, err := c.WeatherByRadius(38.72, -9.13, 10, "metric")
	// If there is an error, print it and terminate the program.
	if err != nil {
		log.Fatal(err)
	}
	// Print a string representation of w using the Stringer interface.
	for i := range w {
		fmt.Println(&w[i])
	}
}
