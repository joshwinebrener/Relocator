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
	housingPrice float64
	violentCrime int
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

	countyData := make(map[string]countyData_t)
	for lineNo, values := range housePricinglines[1:] {
		county := values[2]

		// Normalize county name
		county = strings.TrimSpace(county)
		county = strings.ToLower(county)
		county = strings.Replace(county, " county", "", 1)

		// Get the most recent house cd of the last 6 months
		var price float64
		for i := 1; i <= 6; i++ {
			price, err = strconv.ParseFloat(values[len(values)-i], 32)
			if err != nil {
				if i == 6 {
					fmt.Printf("Error parsing housing price for county \"%s\" at line %d: ", county, lineNo+1)
					fmt.Println(err)
				}
				continue
			}
			break
		}

		// Set violentCrime to -1 so we know if it is never set again
		countyData[county] = countyData_t{housingPrice: price, violentCrime: -1}
	}

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

	for lineNo, values := range violentCrimeLines {
		if lineNo == 0 {
			// Skip header
			continue
		}
		county := values[1]

		// Normalize county name
		county = strings.TrimSpace(county)
		county = strings.ToLower(county)
		county = strings.Replace(county, " police department", "", 1)

		var violentCrimeElement int
		values[2] = strings.Replace(values[2], ",", "", 1)
		violentCrimeElement, err = strconv.Atoi(values[2])
		if err != nil {
			// fmt.Printf("Error parsing crime data for county \"%s\" at line %d: ", county, lineNo+1)
			// fmt.Println(err)
			continue
		}

		cd, ok := countyData[county]
		if ok {
			cd.violentCrime = violentCrimeElement
			countyData[county] = cd
		}
	}

	// Remove elements that have only housing price and not violent crime
	for county, data := range countyData {
		if data.violentCrime < 0 {
			delete(countyData, county)
		}
	}

	// Basic assertions.  TODO: remove before submission.
	if (2328182.9 < countyData["weston"].housingPrice &&
		countyData["weston"].housingPrice < 2328183.1) ||
		(countyData["weston"].violentCrime != 2) {
		panic(fmt.Sprintf("Housing price or violent crime sanity check failed: %f, %d",
			countyData["weston"].housingPrice,
			countyData["weston"].violentCrime))
	}

	outputBuffer := ""
	for county, data := range countyData {
		outputBuffer += fmt.Sprintf("%s,%.1f,%d\n", county, data.housingPrice, data.violentCrime)
	}
	os.WriteFile("output.csv", []byte(outputBuffer), 0644)
}
