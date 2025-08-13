//Bullshit name please change to someting less terrible

package fetchers

import (
	"example.com/mod/classes"
	"github.com/joho/godotenv"
	"io"
	"log"
	"net/http"
	"os"
)

const weatherUrl string = "https://api.openweathermap.org/data/2.5/weather?lat=48.1486&lon=17.1077&appid="

func GetWeatherFromLatLig(lat, lon float64) *classes.Weather {

	err := godotenv.Load()

	if err != nil {
		log.Fatal(err)
	}

	openWeatherApi := os.Getenv("OPENWEATHERAPI")

	resp, err := http.Get(weatherUrl + openWeatherApi)

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
