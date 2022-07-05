/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

// geolocCmd represents the geoloc command
var geolocCmd = &cobra.Command{
	Use:   "geoloc",
	Short: "This command gives you geolocation info",
	Long: `This command gives you the opportunity to access
	geolocation information for a specific city or zip code`,
	Run: func(cmd *cobra.Command, args []string) {
		cityName, _ := cmd.Flags().GetString("city-name")
		usCode, _ := cmd.Flags().GetString("us-code")
		countryCode, _ := cmd.Flags().GetString("country-code")
		getGeolocation(cityName, usCode, countryCode)
	},
}

func init() {
	rootCmd.AddCommand(geolocCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// geolocCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// geolocCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	geolocCmd.Flags().String("city-name", "mieres", "Name of the city")
	geolocCmd.Flags().String("us-code", "", "State code for US only")
	geolocCmd.Flags().String("country-code", "esp", "ISO 3166 country code")
}



type Geoloc struct {
	Lat float32 `json:"lat"`
	Lon float32 `json:"lon"`
	Name string `json:"name"`
}


func getGeolocation(cityName, usCode, countryCode string) {
	url := "http://api.openweathermap.org/geo/1.0/direct"
	responseBytes := getLocationData(url, cityName, usCode, countryCode)
	geoloc := []Geoloc{}

	if err := json.Unmarshal(responseBytes, &geoloc); err != nil {
		log.Printf("Could not unmarshall response - %v", err)
	}

	for _, item := range geoloc {
		fmt.Printf("City with name %s have longitud %f and latitud %f\n", item.Name, item.Lon, item.Lat)
	}
}


func getLocationData(url, cityName, usCode, countryCode string) []byte {
	request, err := http.NewRequest(
		http.MethodGet,
		url,
		nil,
	)

	if err != nil {
		log.Printf("Could not request a geoloc data - %v", err)
	}

	// Configure the query parameters
	q := request.URL.Query()
	qParam := []string{cityName, usCode, countryCode}
	q.Add("q", strings.Join(qParam, ","))
	err = godotenv.Load(".env")
	if err != nil {
		log.Printf("Could not read the .env file - %v", err)
	}
	q.Add("appid", os.Getenv("OPENWEATHER_API_KEY"))
	request.URL.RawQuery = q.Encode()

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Printf("Could not make a request - %v", err)
	}

	responseByte, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("Could not read response body - %v", err)
	}
	return responseByte
}


