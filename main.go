package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	"strings"
)

type Weather struct {
	Location string `json:"timezone"`
	Timezone string `json:"timezone_abbreviation"`
	Latitude float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`

	Units struct {
		Time string `json:"time"`
		Interval string `json:"interval"`
		Temperature string `json:"temperature_2m"`
		Precipitation string `json:"precipitation"`
	} `json:"current_units"`

	Current struct {
		Time string `json:"time"`
		Temperature float64 `json:"temperature_2m"`
	} `json:"current"`

	Hourly struct {
		Time []string `json:"time"`
		Temperature []float64 `json:"temperature_2m"`
	} `json:"hourly"`

}

type Pair[T, U any] struct {
    First  T
    Second U
}

func Zip[T, U any](ts []T, us []U) []Pair[T,U] {
    if len(ts) != len(us) {
        fmt.Errorf("zip: slices have different length.")
    }
    pairs := make([]Pair[T,U], len(ts))
    for i := 0; i < len(ts); i++ {
        pairs[i] = Pair[T,U]{ts[i], us[i]}
    }
    return pairs
}

func main() {
	var get_url string = "https://api.open-meteo.com/v1/forecast?latitude=43.7001&longitude=-79.4163&current=temperature_2m,relative_humidity_2m,is_day,precipitation,rain,showers,snowfall,wind_speed_10m&hourly=temperature_2m,precipitation_probability,precipitation,rain,wind_speed_10m&daily=temperature_2m_max,temperature_2m_min,sunrise,sunset,daylight_duration,uv_index_max&temperature_unit=fahrenheit&timezone=America%2FNew_York"
	get_url = strings.Replace(get_url, "fahrenheit", "celsius", -1)
	
	res, err := http.Get(get_url)
	if err != nil {
		panic(err)
	}
	fmt.Println("Http status code: " + res.Status)
	defer res.Body.Close()

	if res.StatusCode != 200 {
		fmt.Println("Http status code: " + res.Status)
		panic("Weather API not available")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	var weather Weather
	err = json.Unmarshal(body, &weather)
	if err != nil {
		panic(err)
	}

	location, current, hours := weather.Location, weather.Current, weather.Hourly

	fmt.Printf(
		"%s, %s, Current Temperature: %.0f°C\n",
		location,
		weather.Timezone,
		current.Temperature,
	)

	p := Zip(hours.Time, hours.Temperature)

	for _, hourly := range p {

		format := "2006-01-02 15:04"
		t, err := time.Parse(format, strings.Replace(hourly.First, "T", " ", -1))
		if err != nil {
			fmt.Println(err)
		}

		loc, _ := time.LoadLocation("America/New_York")

		if t.Before(time.Now().UTC().In(loc)) {
			continue
		}

		fmt.Printf(
			"%s - %.0f°C\n",
			t.Format(format),
			hourly.Second,
		)
	}
}