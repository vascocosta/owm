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
)

type weatherLine struct {
	Id          int
	Main        string
	Description string
	Icon        string
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

// Weather represents the current weather at a specific location.
//
// It is returned by the WeatherBy* methods of Client.
type Weather struct {
	Coord struct {
		Lat float64
		Lon float64
	}
	Sys struct {
		Type    int
		Id      int
		Country string
		Sunrise int
		Sunset  int
	}
	Weather []weatherLine
	Base    string
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

// WeatherSet represents the current weather at a set of specific locations.
//
// It is returned by the WeatherSetBy* methods of Client.
type WeatherSet struct {
	Cnt     int
	Weather []Weather `json:"list"`
}

// Forecast represents the weather forecast for a specific location.
//
// It is returned by the ForecastBy* methods of Client.
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

// Client represents an OpenWeatherMap API client.
type Client struct {
	key     string
	baseURL string
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
	if c.key != "" {
		url += "&APPID=" + c.key
	}
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

func (c *Client) weatherSet(url string) (ws WeatherSet, err error) {
	if c.key != "" {
		url += "&APPID=" + c.key
	}
	data, err := c.data(url)
	if err != nil {
		err = errors.New("owm: error while decoding weather")
		return
	}
	err = json.Unmarshal(data, &ws)
	if err != nil {
		err = errors.New("owm: error while decoding weather")
		return
	}
	return
}

// WeatherSetByIds decodes the current weather given a slice of city ids and the units.
func (c *Client) WeatherSetByIds(id []int, units string) (ws WeatherSet, err error) {
	var ids string
	for i := range id {
		ids += strconv.Itoa(id[i]) + ","
	}
	ws, err = c.weatherSet(c.baseURL +
		"group" +
		"?id=" + ids[:len(ids)-1] +
		"&units=" + units)
	if err != nil {
		err = errors.New("owm: error while decoding weather")
	}
	return
}
