package main

type (
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
