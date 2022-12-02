package main

import (
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
	b, err := os.ReadFile("zillow_home_prices.csv")
	if err != nil {
		panic(err)
	}
	housingPricesCsv := string(b)
	lines := strings.Split(housingPricesCsv, "\n")
	zipCodes := make([]int, len(lines)-1)      // subtract header row
	housingPrices := make([]int, len(lines)-1) // subtract header row
	for i, line := range lines {
		if i != 0 {
			if strings.TrimSpace(line) != "" {
				values := strings.Split(line, ",")
				housingPrices[i-1], err = strconv.Atoi(values[len(values)-1])
				zipCodes[i-1], err = strconv.Atoi(values[2])
			}
		}
	}

	fmt.Println(housingPrices)
	// fmt.Println(zipCodes)
}
