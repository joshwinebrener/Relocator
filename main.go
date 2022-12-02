package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
)

// Do we have a decided naming convention?
type countyData_t struct {
	housingPrices map[string]float64
	violentCrime  map[string]int
}

func main() {
	// Read config file
	var conf Config
	_, err := toml.DecodeFile("config.toml", &conf)
	if err != nil {
		panic(err)
	}

	// Read housing pricing CSV
	housePricingfile, err := os.Open("zillow_home_prices.csv")
	if err != nil {
		panic(err)
	}
	defer housePricingfile.Close()

	housePricinglines, err := csv.NewReader(housePricingfile).ReadAll()
	if err != nil {
		panic(err)
	}

	housingPrices := make(map[string]float64)
	for lineNo, values := range housePricinglines {

		if lineNo == 0 {
			// Skip header
			continue
		}
		county := values[2]

		// Normalize county name
		county = strings.TrimSpace(county)
		county = strings.ToLower(county)
		county = strings.Replace(county, " county", "", 1)

		// Get the most recent house price of the last 6 months
		var price float64
		for i := 1; i <= 6; i++ {
			price, err = strconv.ParseFloat(values[len(values)-i], 32)
			if err != nil {
				if i == 6 {
					fmt.Printf("Error parsing county \"%s\" at line %d: ", county, lineNo+1)
					fmt.Println(err)
				}
				continue
			}
			break
		}
		housingPrices[county] = price
	}
	countyData := countyData_t{housingPrices: housingPrices}

	// Read crime data CSV
	violentCrimefile, err := os.Open("cjis_crime_data.csv")
	if err != nil {
		panic(err)
	}
	defer violentCrimefile.Close()

	violentCrimeLines, err := csv.NewReader(violentCrimefile).ReadAll()
	if err != nil {
		panic(err)
	}

	violentCrime := make(map[string]int)
	for lineNo, values := range violentCrimeLines {

		if lineNo == 0 {
			// Skip header
			continue
		}
		county := values[1]

		// Normalize county name
		county = strings.TrimSpace(county)
		county = strings.ToLower(county)

		var violentCrimeElement int
		violentCrimeElement, err = strconv.Atoi(values[2])

		violentCrime[county] = violentCrimeElement
	}

	countyData.violentCrime = violentCrime

	fmt.Println(countyData.violentCrime["weston"])
	fmt.Println(countyData.housingPrices["weston"])

	// Basic assertions.  TODO: remove before submission.
	if (countyData.housingPrices["weston"] != 2328183) || (countyData.violentCrime["weston"] != 2.0) {
		panic("Housing price or violent crime sanity check failed")
	}

}
