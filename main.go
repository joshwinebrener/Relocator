package main

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

type (
	userInput struct {
		// weights     configParameters
		// min_filters configParameters
		// max_filters configParameters
		median_housing_price int //`toml:"median_housing_price"`
	}
	// configParameters struct {
	// 	median_housing_price int `toml:median_housing_price`
	// }
	Config struct {
		Weights    ConfigParameters `toml:"weights"`
		MinFilters ConfigParameters `toml:"min_filters"`
		MaxFilters ConfigParameters `toml:"max_filters"`
	}
	ConfigParameters struct {
		MedianHousingPrice           int `toml:"median_housing_price"`
		ViolentCrimeIncidentsPerYear int `toml:"violent_crime_incidents_per_year"`
	}
)

func main() {
	// var input userInput
	var conf Config
	// _, err := toml.DecodeFile("input.toml", &input)
	_, err := toml.DecodeFile("input.toml", &conf)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", conf)
}
