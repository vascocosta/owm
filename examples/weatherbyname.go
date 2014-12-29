package main

import (
	"fmt"
	"github.com/vascocosta/owm"
	"log"
)

func main() {
	// Create a new Client given an API key.
	c := owm.NewClient("YOUR_OPEN_WEATHER_MAP_API_KEY")
	// Decode the current weather of a location given the city name and
	// units. WeatherByName returns a Weather.
	w, err := c.WeatherByName("Lisbon", "metric")
	// If there is an error, print it and terminate the program.
	if err != nil {
		log.Fatal(err)
	}
	// Print a string representation of w using the Stringer interface.
	fmt.Println(&w)
}
