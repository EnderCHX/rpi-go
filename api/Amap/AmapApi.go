package Amap

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

var (
	apiUrl = "https://restapi.amap.com/v3/weather/weatherInfo"
)

type Weather struct {
	Key      string
	CityCode string
	Result   WeatherResponse
}

type WeatherResponse struct {
	Status   string `json:"status"`
	Count    string `json:"count"`
	Info     string `json:"info"`
	InfoCode string `json:"infocode"`
	Lives    []struct {
		Province         string `json:"province"`
		City             string `json:"city"`
		Adcode           string `json:"adcode"`
		Weather          string `json:"weather"`
		Temperature      string `json:"temperature"`
		WindDirection    string `json:"windDirection"`
		WindPower        string `json:"windPower"`
		Humidity         string `json:"humidity"`
		ReportTime       string `json:"reportTime"`
		TemperatureFloat string `json:"temperature_float"`
		HumidityFloat    string `json:"humidity_float"`
	}
}

func (w *Weather) GetWeather(cityCode string) WeatherResponse {
	var url string
	if len(cityCode) == 0 {
		url = fmt.Sprintf("%s?&city=%skey=%s", apiUrl, w.CityCode, w.Key)
	} else {
		url = fmt.Sprintf("%s?&city=%s&key=%s", apiUrl, cityCode, w.Key)
	}

	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	bodyJson := WeatherResponse{}
	json.Unmarshal(body, &bodyJson)

	w.Result = bodyJson

	return bodyJson
}
