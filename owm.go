// Copyright 2014 Vasco Costa. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package owm implements a wrapper for the Open Weather Map API.
//
// Read more about the original API here: http://openweathermap.org/api.
package owm

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Weather represents the current weather at a specific location.
//
// It is returned by the WeatherBy* methods of Client.
type Weather struct {
	Coord struct {
		Lat float64 `json:"lat"` // City latitude.
		Lon float64 `json:"lon"` // City longitude.
	}
	Sys struct {
		Type    int    `json:"type"`    // Unused field.
		Id      int    `json:"id"`      // Unused field.
		Country string `json:"country"` // Country.
		Sunrise int    `json:"sunrise"` // Sunrise unix timestamp.
		Sunset  int    `json:"sunset"`  // Sunset unix timestamp.
	}
	Weather []struct {
		Id          int    `json:"id"`          // Weather condition id.
		Main        string `json:"main"`        // Group of weather parameters.
		Description string `json:"description"` // Weather condition within the group.
		Icon        string `json:"icon"`        // Weather icon id.
	}
	Base string `json:"base"` // Unused field.
	Main struct {
		Temp      float64 `json:"temp"`       // Temperature.
		Humidity  int     `json:"humidity"`   // Humidity.
		TempMin   float64 `json:"temp_min"`   // Minimum temperature.
		TempMax   float64 `json:"temp_max"`   // Maximum temperature.
		Pressure  float64 `json:"pressure"`   // Atmospheric pressure.
		SeaLevel  float64 `json:"sea_level"`  // Sea level atmospheric pressure.
		GrndLevel float64 `json:"grnd_level"` // Ground level atmospheric pressure.
	}
	Wind struct {
		Speed float64 `json:"speed"` // Wind speed.
		Deg   float64 `json:"deg"`   // Wind direction.
		Gust  float64 `json:"gust"`  // Wind gust.
	}
	Clouds struct {
		All int `json:"all"` // Cloudiness.
	}
	Dt   int    `json:"dt"`   // Data unix timestamp.
	Id   int    `json:"id"`   // City identification.
	Name string `json:"name"` // City name.
	Cod  int    `json:"cod"`  // Unused field.
}

// String returns a string representation of Weather by implementing Stringer.
func (w *Weather) String() string {
	return fmt.Sprintf("Date: %v\n"+
		"Id: %v\n"+
		"Name: %v\n"+
		"Country: %v\n"+
		"Sunrise: %v\n"+
		"Sunset: %v\n"+
		"Lat: %v\n"+
		"Lon: %v\n"+
		"Weather Id: %v\n"+
		"Weather Main: %v\n"+
		"Weather Description: %v\n"+
		"Weather Icon: %v\n"+
		"Temp: %v\n"+
		"Humidity: %v\n"+
		"Temp Min: %v\n"+
		"Temp Max: %v\n"+
		"Pressure: %v\n"+
		"Sea Level: %v\n"+
		"Ground Level: %v\n"+
		"Wind Speed: %v\n"+
		"Wind Direction: %v\n"+
		"Wind Gust: %v\n"+
		"Cloudiness: %v",
		time.Unix(int64(w.Dt), 0),
		w.Id,
		w.Name,
		w.Sys.Country,
		time.Unix(int64(w.Sys.Sunrise), 0),
		time.Unix(int64(w.Sys.Sunset), 0),
		w.Coord.Lat,
		w.Coord.Lon,
		w.Weather[0].Id,
		w.Weather[0].Main,
		w.Weather[0].Description,
		w.Weather[0].Icon,
		w.Main.Temp,
		w.Main.Humidity,
		w.Main.TempMin,
		w.Main.TempMax,
		w.Main.Pressure,
		w.Main.SeaLevel,
		w.Main.GrndLevel,
		w.Wind.Speed,
		w.Wind.Deg,
		w.Wind.Gust,
		w.Clouds.All)
}

type weatherSet struct {
	Cnt     int       `json:"cnt"`  // Weather line count.
	Weather []Weather `json:"list"` // Weather line.
}

// Forecast represents the weather forecast for a specific location.
//
// It is returned by the ForecastBy* methods of Client.
type Forecast struct {
	Cod     string  `json:"cod"`
	Message float64 `json:"message"`
	City    struct {
		Id      int    `json:"id"`
		Name    string `json:"name"`
		Country string `json:"country"`
		Coord   struct {
			Lat float64 `json:"lat"`
			Lon float64 `json:"lon"`
		} `json:"coord"`
		Population int `json:"population"`
		Sys        struct {
			Population int `json:"population"`
		} `json:"sys"`
	} `json:"city"`
	Cnt      int `json:"cnt"`
	Forecast []struct {
		Dt   int `json:"dt"`
		Main struct {
			Temp      float64 `json:"temp"`
			Humidity  int     `json:"humidity"`
			TempMin   float64 `json:"temp_min"`
			TempMax   float64 `json:"temp_max"`
			Pressure  float64 `json:"pressure"`
			SeaLevel  float64 `json:"sea_level"`
			GrndLevel float64 `json:"grnd_level"`
		} `json:"main"`
		Temp struct {
			Day   float64 `json:"day"`
			Min   float64 `json:"min"`
			Max   float64 `json:"max"`
			Morn  float64 `json:"morn"`
			Eve   float64 `json:"eve"`
			Night float64 `json:"night"`
		} `json:"temp"`
		Pressure float64 `json:"pressure"`
		Humidity int     `json:"humidity"`
		Weather  []struct {
			Id          int    `json:"id"`
			Main        string `json:"main"`
			Description string `json:"description"`
			Icon        string `json:"icon"`
		} `json:"weather"`
		Wind struct {
			Speed float64 `json:"speed"`
			Deg   float64 `json:"deg"`
			Gust  float64 `json:"gust"`
		} `json:"wind"`
		Speed  float64     `json:"speed"`
		Deg    int         `json:"deg"`
		Gust   float64     `json:"gust"`
		Clouds interface{} `json:"clouds"`
	} `json:"list"`
}

func (f *Forecast) String() string {
	return fmt.Sprintf("Id: %v\n"+
		"Name: %v\n"+
		"Country: %v\n"+
		"Lat: %v\n"+
		"Lon: %v\n"+
		"Population: %v\n"+
		"Forecast: %v\n",
		f.City.Id,
		f.City.Name,
		f.City.Country,
		f.City.Coord.Lat,
		f.City.Coord.Lon,
		f.City.Population,
		f.Forecast)
}

// Client represents an OpenWeatherMap API client.
type Client struct {
	key     string // API key.
	baseURL string // API base URL.
}

// NewClient returns a new Client given an API key.
//
// Pass an empty string as the key argument to use the client without an API key.
func NewClient(key string) *Client {
	return &Client{key, "http://api.openweathermap.org/data/2.5/"}
}

func (c *Client) data(url string) (data []byte, err error) {
	res, err := http.Get(url)
	if err != nil {
		err = errors.New("owm: error while fetching weather data")
		return
	}
	data, err = ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		err = errors.New("owm: error while fetching weather data")
		return
	}
	return
}

func (c *Client) weather(url string) (w Weather, err error) {
	if c.key != "" {
		url += "&APPID=" + c.key
	}
	data, err := c.data(url)
	if err != nil || strings.Contains(string(data), `"cod":"404"`) {
		err = errors.New("owm: error while fetching weather data")
		return
	}
	err = json.Unmarshal(data, &w)
	if err != nil {
		err = errors.New("owm: error while decoding weather data")
		return
	}
	return
}

func (c *Client) weatherSet(url string) (ws weatherSet, err error) {
	if c.key != "" {
		url += "&APPID=" + c.key
	}
	data, err := c.data(url)
	if err != nil || strings.Contains(string(data), `"cod":"404"`) {
		err = errors.New("owm: error while fetching weather data")
		return
	}
	err = json.Unmarshal(data, &ws)
	if err != nil {
		err = errors.New("owm: error while decoding weather data")
		return
	}
	return
}

func (c *Client) forecast(url string) (f Forecast, err error) {
	if c.key != "" {
		url += "&APPID=" + c.key
	}
	data, err := c.data(url)
	if err != nil || strings.Contains(string(data), `"cod":"404"`) {
		err = errors.New("owm: error while fetching forecast data")
		return
	}
	err = json.Unmarshal(data, &f)
	if err != nil {
		err = errors.New("owm: error while decoding forecast data")
		return
	}
	return
}

// WeatherByName decodes the current weather of a location given the city name
// and units. It uses the corresponding web API URL to fetch JSON encoded data
// and returns a Weather with as much fields decoded as those available.
//
// An error is returned if there is a problem while fetching weather data from
// the web API or decoding the weather data.
func (c *Client) WeatherByName(name string, units string) (w Weather, err error) {
	w, err = c.weather(c.baseURL +
		"weather" +
		"?q=" + name +
		"&units=" + units)
	return
}

// WeatherById decodes the current weather of a location given the city id and
// units. It uses the corresponding web API URL to fetch JSON encoded data and
// returns a Weather with as much fields decoded as those available.
//
// An error is returned if there is a problem while fetching weather data from
// the web API or decoding the weather data.
func (c *Client) WeatherById(id int, units string) (w Weather, err error) {
	w, err = c.weather(c.baseURL +
		"weather" +
		"?id=" + strconv.Itoa(id) +
		"&units=" + units)
	return
}

// WeatherByCoord decodes the current weather of a location given the city
// coordinates and units. It uses the corresponding web API URL to fetch JSON
// encoded data and returns a Weather with as much fields decoded as those
// available.
//
// An error is returned if there is a problem while fetching weather data from
// the web API or decoding the weather data.
func (c *Client) WeatherByCoord(lat, lon float64, units string) (w Weather, err error) {
	w, err = c.weather(c.baseURL +
		"weather" +
		"?lat=" + strconv.FormatFloat(lat, 'f', 2, 64) +
		"&lon=" + strconv.FormatFloat(lon, 'f', 2, 64) +
		"&units=" + units)
	return
}

// WeatherByZone decodes the current weather of multiple locations given the
// zone coordinates, map zoom and units. It uses the corresponding web API URL
// to fetch JSON encoded data and returns a []Weather with as much fields
// decoded as those available.
//
// An error is returned if there is a problem while fetching weather data from
// the web API or decoding the weather data.
func (c *Client) WeatherByZone(lat1, lon1, lat2, lon2 float64, zoom int, units string) (w []Weather, err error) {
	ws, err := c.weatherSet(c.baseURL +
		"box/city" +
		"?bbox=" +
		strconv.FormatFloat(lat1, 'f', 2, 64) + "," +
		strconv.FormatFloat(lon1, 'f', 2, 64) + "," +
		strconv.FormatFloat(lat2, 'f', 2, 64) + "," +
		strconv.FormatFloat(lon2, 'f', 2, 64) + "," +
		strconv.Itoa(zoom) +
		"&units=" + units)
	w = ws.Weather
	return
}

// WeatherByRadius decodes the current weather of multiple locations given the
// center coordinates, radius and units. It uses the corresponding web API URL
// to fetch JSON encoded data and returns a []Weather with as much fields
// decoded as those available.
//
// An error is returned if there is a problem while fetching weather data from
// the web API or decoding the weather data.
func (c *Client) WeatherByRadius(lat, lon, radius float64, units string) (w []Weather, err error) {
	ws, err := c.weatherSet(c.baseURL +
		"find" +
		"?lat=" + strconv.FormatFloat(lat, 'f', 2, 64) + "," +
		"&lon=" + strconv.FormatFloat(lon, 'f', 2, 64) + "," +
		"&cnt=" + strconv.FormatFloat(radius, 'f', 2, 64) +
		"&units=" + units)
	w = ws.Weather
	return
}

// WeatherByIds decodes the current weather of multiple locations given a slice
// of city ids and units. It uses the corresponding web API URL to fetch JSON
// encoded data and returns a []Weather with as much fields decoded fields as
// those available.
//
// An error is returned if there is a problem while fetching weather data from
// the web API or decoding the weather data.
func (c *Client) WeatherByIds(id []int, units string) (w []Weather, err error) {
	var ids string
	for i := range id {
		ids += strconv.Itoa(id[i]) + ","
	}
	ws, err := c.weatherSet(c.baseURL +
		"group" +
		"?id=" + ids[:len(ids)-1] +
		"&units=" + units)
	w = ws.Weather
	return
}

// ForecastByName decodes the current forecast of a location given the city
// name, days and units. If days is equal to 0, it returns an hourly forecast,
// otherwise if days is greater than 0, it returns a daily forecast for the
// given number of days. It uses the corresponding web API URL to fetch JSON
// encoded data and returns a Forecast with as much fields decoded as those
// available.
//
// An error is returned if there is a problem while fetching forecast data from
// the web API or decoding the forecast data.
func (c *Client) ForecastByName(name string, days int, units string) (f Forecast, err error) {
	if days > 0 {
		f, err = c.forecast(c.baseURL +
			"forecast/daily" +
			"?q=" + name +
			"&cnt=" + strconv.Itoa(days) +
			"&units=" + units)
	} else {
		f, err = c.forecast(c.baseURL +
			"forecast" +
			"?q=" + name +
			"&units=" + units)
	}
	return
}

// ForecastById decodes the current forecast of a location given the city id,
// days and units. If days is equal to 0, it returns an hourly forecast,
// otherwise if days is greater than 0, it returns a daily forecast for the
// given number of days. It uses the corresponding web API URL to fetch JSON
// encoded data and returns a Forecast with as much fields decoded as those
// available.
//
// An error is returned if there is a problem while fetching forecast data from
// the web API or decoding the forecast data.
func (c *Client) ForecastById(id, days int, units string) (f Forecast, err error) {
	if days > 0 {
		f, err = c.forecast(c.baseURL +
			"forecast/daily" +
			"?id=" + strconv.Itoa(id) +
			"&cnt=" + strconv.Itoa(days) +
			"&units=" + units)
	} else {
		f, err = c.forecast(c.baseURL +
			"forecast" +
			"?id=" + strconv.Itoa(id) +
			"&units=" + units)
	}
	return
}

// ForecastByCoord decodes the current forecast of a location given the city
// coordinates, days and units. If days is equal to 0, it returns an hourly
// forecast, otherwise if days is greater than 0, it returns a daily forecast
// for the given number of days. It uses the corresponding web API URL to fetch
// JSON encoded data and returns a Forecast with as much fields decoded as those
// available.
//
// An error is returned if there is a problem while fetching forecast data from
// the web API or decoding the forecast data.
func (c *Client) ForecastByCoord(lat, lon float64, days int, units string) (f Forecast, err error) {
	if days > 0 {
		f, err = c.forecast(c.baseURL +
			"forecast/daily" +
			"?lat=" + strconv.FormatFloat(lat, 'f', 2, 64) +
			"&lon=" + strconv.FormatFloat(lon, 'f', 2, 64) +
			"&cnt=" + strconv.Itoa(days) +
			"&units=" + units)
	} else {
		f, err = c.forecast(c.baseURL +
			"forecast" +
			"?lat=" + strconv.FormatFloat(lat, 'f', 2, 64) +
			"&lon=" + strconv.FormatFloat(lon, 'f', 2, 64) +
			"&units=" + units)
	}
	return
}
