package main

func Rank(countyData map[string]countyData_t, config Config) map[string]float64 {
	rankings := make(map[string]float64)

	c := make(chan map[string]float64)

	keys := make([]string, len(countyData))
	i := 0
	for k := range countyData {
		keys[i] = k
		i++
	}

	for i := 0; i < len(keys); i += 100 {
		keysLeft := len(keys) - i
		n := 0
		if keysLeft > 100 {
			n = 100
		} else {
			n = keysLeft
		}
		var keySlice = keys[i : i+n]
		go func() {
			rankingsSlice := make(map[string]float64)
			for _, key := range keySlice {
				rankingsSlice[key] = countyData[key].housingPrice*float64(config.Weights.MedianHousingPrice)/100 +
					float64(countyData[key].violentCrime)*float64(config.Weights.ViolentCrimeIncidentsPerYear)/100
			}
			c <- rankingsSlice
		}()
	}

	for len(rankings) < len(countyData) {
		rankingsSlice := <-c
		for k, v := range rankingsSlice {
			rankings[k] = v
		}
	}

	return rankings
}
