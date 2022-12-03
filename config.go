package main

type (
	Config struct {
		Weights    ConfigParameters `toml:"weights"`
		MinFilters ConfigParameters `toml:"min_filters"`
		MaxFilters ConfigParameters `toml:"max_filters"`
	}
	ConfigParameters struct {
		MedianHousingPrice                   float64 `toml:"median_housing_price"`
		YearlyViolentCrimeIncidentsPerCapita float64 `toml:"yearly_violent_crime_incidents_per_capita"`
		Population                           float64 `toml:"population"`
	}
)
