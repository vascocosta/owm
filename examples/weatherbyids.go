package main

import (
	"fmt"
	"github.com/vascocosta/owm"
	"log"
)

func main() {
	// Create a new Client given an API key.
	c := owm.NewClient("YOUR_OPEN_WEATHER_MAP_API_KEY")
	// Decode the current weather of multiple locations given a slice of
	// city ids and units. WeatherByIds returns a []Weather.
	w, err := c.WeatherByIds([]int{2267057, 2735943, 2268339}, "metric")
	// If there is an error, print it and terminate the program.
	if err != nil {
		log.Fatal(err)
	}
	// Print a string representation of w using the Stringer interface.
	for i := range w {
		fmt.Println(&w[i])
	}
}
