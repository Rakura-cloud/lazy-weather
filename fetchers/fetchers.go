//Bullshit name please change to someting less terrible

package fetchers

import (
	"io"
	"log"
	"net/http"

	"example.com/mod/classes"
)

const weatherUrl string = "https://api.openweathermap.org/data/2.5/weather?lat=48.1486&lon=17.1077&appid=9f097f8dd06d8a3f002f5b3bfc8c3012"

func GetWeatherFromLatLig(lat, lon float64) *classes.Weather {

	resp, err := http.Get(weatherUrl)

	if err != nil {
		log.Fatal(err)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	tmpWeather, error := classes.ParseWeather(bodyBytes)

	if error != nil {
		log.Fatal(error)
	}

	return tmpWeather
}
