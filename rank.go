package main

type ranked struct {
	county string
	rank   float64
}

func rank(countyData map[string]countyData_t, config Config) []ranked {
	rankings := make([]ranked, len(countyData))

	c := make(chan []ranked)

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
			rankingsSlice := make([]ranked, n)
			j := 0
			for _, key := range keySlice {
				rankingsSlice[j] = ranked{
					key,
					countyData[key].housingPrice*float64(config.Weights.MedianHousingPrice)/100 +
						float64(countyData[key].violentCrime)*float64(config.Weights.ViolentCrimeIncidentsPerYear)/100,
				}
			}
			c <- rankingsSlice
		}()
	}

	i = 0
	for i < len(rankings) {
		rankingsSlice := <-c
		for _, v := range rankingsSlice {
			rankings[i] = v
			i++
		}
	}

	return rankings
}
