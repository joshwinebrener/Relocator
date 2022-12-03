package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"

	"github.com/BurntSushi/toml"
)

type countyData_t struct {
	housingPrice float64
	violentCrime float64
	population   float64
}

func main() {
	// // A handy map of US state codes to full names
	var stateCodeNameMap = map[string]string{
		"AL": "alabama",
		"AK": "alaska",
		"AZ": "arizona",
		"AR": "arkansas",
		"CA": "california",
		"CO": "colorado",
		"CT": "connecticut",
		"DE": "delaware",
		"FL": "florida",
		"GA": "georgia",
		"HI": "hawaii",
		"ID": "idaho",
		"IL": "illinois",
		"IN": "indiana",
		"IA": "iowa",
		"KS": "kansas",
		"KY": "kentucky",
		"LA": "louisiana",
		"ME": "maine",
		"MD": "maryland",
		"MA": "massachusetts",
		"MI": "michigan",
		"MN": "minnesota",
		"MS": "mississippi",
		"MO": "missouri",
		"MT": "montana",
		"NE": "nebraska",
		"NV": "nevada",
		"NH": "new hampshire",
		"NJ": "new jersey",
		"NM": "new mexico",
		"NY": "new york",
		"NC": "north carolina",
		"ND": "north dakota",
		"OH": "ohio",
		"OK": "oklahoma",
		"OR": "oregon",
		"PA": "pennsylvania",
		"RI": "rhode island",
		"SC": "south carolina",
		"SD": "south dakota",
		"TN": "tennessee",
		"TX": "texas",
		"UT": "utah",
		"VT": "vermont",
		"VA": "virginia",
		"WA": "washington",
		"WV": "west virginia",
		"WI": "wisconsin",
		"WY": "wyoming",
		// Territories
		"AS": "american samoa",
		"DC": "district of columbia",
		"FM": "federated states of micronesia",
		"GU": "guam",
		"MH": "marshall islands",
		"MP": "northern mariana islands",
		"PW": "palau",
		"PR": "puerto rico",
		"VI": "virgin islands",
		// Armed Forces (AE includes Europe, Africa, Canada, and the Middle East)
		"AA": "armed forces americas",
		"AE": "armed forces europe",
		"AP": "armed forces pacific",
	}

	// Read config file
	var conf Config
	_, err := toml.DecodeFile("config.toml", &conf)
	if err != nil {
		panic(err)
	}

	// If the toml feilds haven't been set set default values
	if conf.MinFilters.MedianHousingPrice <= 0 {
		conf.MinFilters.MedianHousingPrice = 0.0
	}
	if conf.MaxFilters.MedianHousingPrice <= 0 {
		conf.MaxFilters.MedianHousingPrice = 80000000000
	}
	if conf.MinFilters.YearlyViolentCrimeIncidentsPerCapita <= 0 {
		conf.MinFilters.YearlyViolentCrimeIncidentsPerCapita = 0
	}
	if conf.MaxFilters.YearlyViolentCrimeIncidentsPerCapita <= 0 {
		conf.MaxFilters.YearlyViolentCrimeIncidentsPerCapita = 80000000000
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
		county = county + ", " + strings.ToLower(strings.TrimSpace(stateCodeNameMap[values[4]]))
		county = Capitalize(county)

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

	// Read population
	censusData, err := os.Open("census.csv")
	if err != nil {
		panic(err)
	}
	defer censusData.Close()
	censusLines, err := csv.NewReader(censusData).ReadAll()
	if err != nil {
		panic(err)
	}

	for lineNo, values := range censusLines {
		if lineNo < 6 {
			// Skip header
			continue
		}

		county := values[0]

		// Normalize county name
		county = strings.TrimSpace(county)
		county = strings.ToLower(county)
		county = strings.Replace(county, " county", "", 1)
		county = strings.Replace(county, ".", "", 1)
		county = Capitalize(county)

		values[3] = strings.Replace(values[3], ",", "", -1)
		population, err := strconv.ParseFloat(values[3], 32)
		if err != nil {
			fmt.Printf("Error parsing crime data for county \"%s\" at line %d: ", county, lineNo+1)
			fmt.Println(err)
			continue
		}

		cd, ok := countyData[county]
		if ok {
			cd.population = population
			countyData[county] = cd
		}

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
		county = county + ", " + strings.ToLower(strings.TrimSpace(values[0]))
		county = Capitalize(county)

		var violentCrimeElement float64
		values[2] = strings.Replace(values[2], ",", "", 1)
		violentCrimeElement, err = strconv.ParseFloat(values[2], 32)
		if err != nil {
			// fmt.Printf("Error parsing crime data for county \"%s\" at line %d: ", county, lineNo+1)
			// fmt.Println(err)
			continue
		}

		cd, ok := countyData[county]
		if ok {
			cd.violentCrime = violentCrimeElement / cd.population
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
	if (countyData["Lincoln, Georgia"].housingPrice != 231617) ||
		(0.0001534095267-0.0000001 < countyData["Lincoln, Georgia"].violentCrime &&
			countyData["Lincoln, Georgia"].violentCrime < 0.0001534095267+0.0000001) {
		panic(fmt.Sprintf("Housing price or violent crime sanity check failed: %f, %f",
			countyData["Lincoln, Georgia"].housingPrice,
			countyData["Lincoln, Georgia"].violentCrime))
	}

	outputBuffer := ""
	for county, data := range countyData {
		outputBuffer += fmt.Sprintf("%s,%.1f,%d\n", county, data.housingPrice, data.violentCrime)
	}
	os.WriteFile("output.csv", []byte(outputBuffer), 0644)

	//Filter county based on user set min and max values for desired properties
	for county, data := range countyData {
		if data.housingPrice < conf.MinFilters.MedianHousingPrice {
			delete(countyData, county)
		} else if data.housingPrice > conf.MaxFilters.MedianHousingPrice {
			delete(countyData, county)
		} else if data.violentCrime < conf.MinFilters.YearlyViolentCrimeIncidentsPerCapita {
			delete(countyData, county)
		} else if data.violentCrime > conf.MaxFilters.YearlyViolentCrimeIncidentsPerCapita {
			delete(countyData, county)
		}
	}

	// Get the min and max countyData values
	var min countyData_t
	var max countyData_t
	for _, data := range countyData {
		min = data
		max = data
		break
	}
	for _, data := range countyData {
		if data.housingPrice < min.housingPrice {
			min.housingPrice = data.housingPrice
		}
		if data.housingPrice > max.housingPrice {
			max.housingPrice = data.housingPrice
		}
		if data.violentCrime < min.violentCrime {
			min.violentCrime = data.violentCrime
		}
		if data.violentCrime > max.violentCrime {
			max.violentCrime = data.violentCrime
		}
		if data.population < min.population {
			min.population = data.population
		}
		if data.population > max.population {
			max.population = data.population
		}
	}

	diff := countyData_t{
		max.housingPrice - min.housingPrice,
		max.violentCrime - min.violentCrime,
		max.population - min.population,
	}
	// Normalize county data
	for k, v := range countyData {
		countyData[k] = countyData_t{
			housingPrice: (v.housingPrice - min.housingPrice) / diff.housingPrice,
			violentCrime: (v.violentCrime - min.violentCrime) / diff.violentCrime,
		}
	}

	ranked := rank(countyData, conf)
	sorted := mergeSort(ranked)
	// fmt.Println(sorted)

	outputBuffer = ""
	for _, ranked := range sorted {
		outputBuffer += fmt.Sprintf("\"%s\",%f\n", ranked.county, ranked.rank)
	}
	os.WriteFile("output.csv", []byte(outputBuffer), 0644)
}

func Capitalize(s string) string {
	rs := []rune(s)
	inWord := false
	for i, r := range rs {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			if !inWord {
				rs[i] = unicode.ToTitle(r)
			}
			inWord = true
		} else {
			inWord = false
		}
	}
	return string(rs)
}
