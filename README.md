# owm

Package owm implements a wrapper for the OpenWeatherMap API.

Read more about the OpenWeatherMap API here: http://openweathermap.org/api.

## Example

```go
package main

import (
	"fmt"
	"github.com/vascocosta/owm"
	"log"
	"time"
)

func main() {
	// Creates a new owm.Client given a key.
	c := owm.NewClient("YOUR_OPEN_WEATHER_MAP_API_KEY")

	// Get the current weather given the city name and the units.
	w, err := c.WeatherByName("Lisbon", "metric")

	// Exit the program if there is an error.
	if err != nil {
		log.Fatal(err)
	}

	// Get the current weather given the city id and the units.
	w, err = c.WeatherById(2643743, "imperial")
	
	// Exit the program if there is an error.
	if err != nil {
		log.Fatal(err)
	}

	// Get the current weather given the city coordinates and the units.
	w, err := c.WeatherByCoord(40.71, -74.01, "kelvin")

	// Exit the program if there is an error.
	if err != nil {
		log.Fatal(err)
	}
}
```
