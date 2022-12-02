package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
)

func main() {
	// Read config file
	var conf Config
	_, err := toml.DecodeFile("config.toml", &conf)
	if err != nil {
		panic(err)
	}

	// Read housing pricing CSV
	file, err := os.Open("zillow_home_prices.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	lines, err := csv.NewReader(file).ReadAll()
	if err != nil {
		panic(err)
	}

	housingPrices := make(map[string]float64)

	for lineNo, values := range lines {
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

	// Basic assertions.  TODO: remove before submission.
	if housingPrices["san augustine"] != 180910.0 {
		panic("housingPrices[\"san augustine\"] != 180910.0")
	}
	if housingPrices["teton"] != 346561.0 {
		panic("housingPrices[\"teton\"] != 346561.0")
	}
	if housingPrices["rock"] != 202578.0 {
		panic("housingPrices[\"rock\"] != 202578.0")
	}

	fmt.Println(housingPrices)
}
