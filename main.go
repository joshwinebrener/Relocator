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
	violentCrime int
}

func main() {
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
	if conf.MinFilters.ViolentCrimeIncidentsPerYear <= 0 {
		conf.MinFilters.ViolentCrimeIncidentsPerYear = 0
	}
	if conf.MaxFilters.ViolentCrimeIncidentsPerYear <= 0 {
		conf.MaxFilters.ViolentCrimeIncidentsPerYear = 80000000000
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
	if (countyData["Lincoln, Georgia"].housingPrice != 231617) || (countyData["Lincoln, Georgia"].violentCrime != 24) {
		panic(fmt.Sprintf("Housing price or violent crime sanity check failed: %f, %d",
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
		} else if data.violentCrime < conf.MinFilters.ViolentCrimeIncidentsPerYear {
			delete(countyData, county)
		} else if data.violentCrime > conf.MaxFilters.ViolentCrimeIncidentsPerYear {
			delete(countyData, county)
		}
	}

	ranked := rank(countyData, conf)
	sorted := mergeSort(ranked)
	fmt.Println(sorted)

	outputBuffer = ""
	for county, data := range countyData {
		outputBuffer += fmt.Sprintf("%s,%.1f,%d\n", county, data.housingPrice, data.violentCrime)
	}
	os.WriteFile("outputFiltered.csv", []byte(outputBuffer), 0644)

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
