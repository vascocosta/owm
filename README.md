# owm

Package owm implements a wrapper for the OpenWeatherMap API.

Read more about the OpenWeatherMap API here: http://openweathermap.org/api.

## Installation

`go get github.com/vascocosta/owm`

## Documentation

`godoc [github.com/vascocosta/owm](http://godoc.org/github.com/vascocosta/owm)`

## Examples

```go
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
```

```go
package main

import (
	"fmt"
	"github.com/vascocosta/owm"
	"log"
)

func main() {
	// Create a new Client given an API key.
	c := owm.NewClient("YOUR_OPEN_WEATHER_MAP_API_KEY")
	// Decode the current weather of a location given the city id and
	// units. WeatherById returns a Weather.
	w, err := c.WeatherById(2267057, "metric")
	// If there is an error, print it and terminate the program.
	if err != nil {
		log.Fatal(err)
	}
	// Print a string representation of w using the Stringer interface.
	fmt.Println(&w)
}
```

```go
package main

import (
	"fmt"
	"github.com/vascocosta/owm"
	"log"
)

func main() {
	// Create a new Client given an API key.
	c := owm.NewClient("YOUR_OPEN_WEATHER_MAP_API_KEY")
	// Decode the current weather of a location given the city coordinates
	// and units. WeatherByCoord returns a Weather.
	w, err := c.WeatherByCoord(38.72, -9.13, "metric")
	// If there is an error, print it and terminate the program.
	if err != nil {
		log.Fatal(err)
	}
	// Print a string representation of w using the Stringer interface.
	fmt.Println(&w)
}
```

```go
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
```

```go
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
```

```go
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
```

```go
package main

import (
	"fmt"
	"github.com/vascocosta/owm"
	"log"
)

func main() {
	// Create a new Client given an API key.
	c := owm.NewClient("YOUR_OPEN_WEATHER_MAP_API_KEY")
	// Decode the current forecast of a location given the city name, days
	// and units. ForecastByName returns a Forecast.
	f, err := c.ForecastByName("Lisbon", 4, "metric")
	// If there is an error, print it and terminate the program.
	if err != nil {
		log.Fatal(err)
	}
	// Print a string representation of f using the Stringer interface.
	fmt.Println(&f)
}
```

```go
package main

import (
	"fmt"
	"github.com/vascocosta/owm"
	"log"
)

func main() {
	// Create a new Client given an API key.
	c := owm.NewClient("YOUR_OPEN_WEATHER_MAP_API_KEY")
	// Decode the current forecast of a location given the city id, days
	// and units. ForecastById returns a Forecast.
	f, err := c.ForecastById(2267057, 8, "metric")
	// If there is an error, print it and terminate the program.
	if err != nil {
		log.Fatal(err)
	}
	// Print a string representation of f using the Stringer interface.
	fmt.Println(&f)
}
```

```go
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
```
