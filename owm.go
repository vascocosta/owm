/*
Package owm implements a wrapper for the OpenWeatherMap API.

Read more about the OpenWeatherMap API here: http://openweathermap.org/api.
*/
package owm

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
)

type weatherLine struct {
	Main        string
	Description string
}

type forecastLine struct {
	Dt   int
	Main struct {
		Temp      float64
		TempMin   float64 `json:"temp_min"`
		TempMax   float64 `json:"temp_max"`
		Pressure  float64
		SeaLevel  float64 `json:"sea_level"`
		GrndLevel float64 `json:"grnd_level"`
		Humidity  int
	}
	Weather []weatherLine
	Wind    struct {
		Speed float64
		Deg   float64
		Gust  float64
	}
	Clouds struct {
		All int
	}
	Sys struct {
		Pod string
	}
	DtTxt string `json:"dt_txt"`
}

type Weather struct {
	Coord struct {
		Lat float64
		Lon float64
	}
	Sys struct {
		Country string
		Sunrise int
		Sunset  int
	}
	Weather []weatherLine
	Main    struct {
		Temp      float64
		Humidity  int
		TempMin   float64 `json:"temp_min"`
		TempMax   float64 `json:"temp_max"`
		Pressure  float64
		SeaLevel  float64 `json:"sea_level"`
		GrndLevel float64 `json:"grnd_level"`
	}
	Wind struct {
		Speed float64
		Deg   float64
		Gust  float64
	}
	Clouds struct {
		All int
	}
	Dt   int
	Id   int
	Name string
	Cod  int
}

type Forecast struct {
	Cod  int
	City struct {
		Id    int
		Name  string
		Coord struct {
			Lon float64
			Lat float64
		}
		Country string
	}
	Cnt          int
	ForecastLine []forecastLine
}

type Client struct {
	key     string
	baseURL string
}

// NewClient returns a new client given a key.
func NewClient(key string) *Client {
	return &Client{key, "http://api.openweathermap.org/data/2.5/"}
}

func (c *Client) data(url string) (data []byte, err error) {
	res, err := http.Get(url)
	if err != nil {
		err = errors.New("owm: error while getting weather data")
		return
	}
	data, err = ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		err = errors.New("owm: error while getting weather data")
		return
	}
	return
}

func (c *Client) weather(url string) (w Weather, err error) {
	data, err := c.data(url)
	if err != nil {
		err = errors.New("owm: error while decoding weather")
		return
	}
	err = json.Unmarshal(data, &w)
	if err != nil {
		err = errors.New("owm: error while decoding weather")
		return
	}
	return
}

// WeatherByName decodes the current weather given the city name and the units.
func (c *Client) WeatherByName(name string, units string) (w Weather, err error) {
	w, err = c.weather(c.baseURL +
		"weather" +
		"?q=" + name +
		"&units=" + units)
	if err != nil {
		err = errors.New("owm: error while decoding weather")
	}
	return
}

// WeatherById decodes the current weather given the city id and the units.
func (c *Client) WeatherById(id int, units string) (w Weather, err error) {
	w, err = c.weather(c.baseURL +
		"weather" +
		"?id=" + strconv.Itoa(id) +
		"&units=" + units)
	if err != nil {
		err = errors.New("owm: error while decoding weather")
	}
	return
}

// WeatherByCoord decodes the current weather given the city coordinates and the units.
func (c *Client) WeatherByCoord(lat float64, lon float64, units string) (w Weather, err error) {
	w, err = c.weather(c.baseURL +
		"weather" +
		"?lat=" + strconv.FormatFloat(lat, 'f', 2, 64) +
		"&lon=" + strconv.FormatFloat(lon, 'f', 2, 64) +
		"&units=" + units)
	if err != nil {
		err = errors.New("owm: error while decoding weather")
	}
	return
}
