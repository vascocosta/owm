// Copyright 2014 Vasco Costa. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package owm implements a wrapper for the OpenWeatherMap API.
//
// Read more about the OpenWeatherMap API here: http://openweathermap.org/api.
package owm

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
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

type weatherSet struct {
	Cnt     int       `json:"cnt"`  // Weather line count.
	Weather []Weather `json:"list"` // Weather line.
}

// Forecast represents the weather forecast for a specific location.
//
// It is returned by the ForecastBy* methods of Client.
type Forecast struct {
	Cod  string `json:"cod"` // Code.
	City struct {
		Id    int    `json:"id"`   // City indentification.
		Name  string `json:"name"` // City name.
		Coord struct {
			Lat float64 `json:"lat"` // City latitude.
			Lon float64 `json:"lon"` // City longitude.
		}
		Country string `json:"country"` // Country.
	}
	Cnt     int       `json:"cnt"`  // Weather line count.
	Weather []Weather `json:"list"` // Weather line
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

// WeatherByIds decodes the current weather of multiple locations given the a
// slice of city ids and units. It uses the corresponding web API URL to fetch
// JSON encoded data and returns a []Weather with as much fields decoded fields
// as those available.
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
			"forecasti/daily" +
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
