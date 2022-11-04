package main

import (
	"fmt"

	"os"

	"github.com/BurntSushi/toml"
)

type UserInput struct {
	weights struct {
		median_housing_price float32 `toml:"weights.median_housing_price"`
	}
	min_filters struct {
		median_housing_price int `toml:"min_filters.median_housing_price"`
	}
	max_filters struct {
		median_housing_price int `toml:"max_filters.median_housing_price"`
	}
}

func main() {
	tomlbytes, err := os.ReadFile("input.toml")
	if err != nil {
		panic(err)
	}
	tomlstr := string(tomlbytes)
	var input UserInput
	_, err2 := toml.Decode(tomlstr, &input)
	if err2 != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", input)
}
