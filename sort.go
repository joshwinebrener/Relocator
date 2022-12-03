package main

import (
	"sync"
)

// Parallel Merge Sort function Implemented using the Go concurrency model
func mergeSort(items []ranked) []ranked {
	if len(items) < 2 {
		return items
	}

	first := items
	second := items

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		first = mergeSort(items[:len(items)/2])
	}()

	go func() {
		defer wg.Done()
		second = mergeSort(items[len(items)/2:])
	}()

	wg.Wait()

	outdata := merge(first, second)
	return outdata
}

func merge(a []ranked, b []ranked) []ranked {
	final := []ranked{}
	i := 0
	j := 0
	for i < len(a) && j < len(b) {
		if a[i].rank < b[j].rank {
			final = append(final, a[i])
			i++
		} else {
			final = append(final, b[j])
			j++
		}
	}
	for ; i < len(a); i++ {
		final = append(final, a[i])
	}
	for ; j < len(b); j++ {
		final = append(final, b[j])
	}
	return final
}
