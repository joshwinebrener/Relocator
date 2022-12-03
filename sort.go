package main

import (
	"sync"
)

// Parallel Merge Sort function Implemented using the Go concurrency model
func mergeSort(items []ranked) []ranked {
	if len(items) < 2 {
		return items
	}

	var first []ranked
	var second []ranked

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		first = mergeSort(items[:len(items)/2])
	}()

	// Here I allow the concurrency of go call the second mergesort call.
	// This way the number of goroutines triggered becomes linear instead of being exponential.
	// https://teivah.medium.com/parallel-merge-sort-in-go-fe14c1bc006
	second = mergeSort(items[len(items)/2:])

	wg.Wait()

	return merge(first, second)
}

// Because sometimes you want to go slow.
func mergeSortSerial(items []ranked) []ranked {
	if len(items) < 2 {
		return items
	}

	first := mergeSortSerial(items[:len(items)/2])
	second := mergeSortSerial(items[len(items)/2:])

	return merge(first, second)
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
