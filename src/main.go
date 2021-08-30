package main

import (
	"fmt"
	"math/rand"
	"time"

	"gonum.org/v1/gonum/stat"
)

func insertionSort(array []int) {
	for i := 1; i < len(array); i++ {
		j := i
		for j > 0 && array[j-1] > array[j] {
			array[j], array[j-1] = array[j-1], array[j]
			j--
		}
	}
}

func selectionSort(array []int) {
	for i := 0; i < len(array); i++ {
		var minIndex = i
		for j := i; j < len(array); j++ {
			if array[j] < array[minIndex] {
				minIndex = j
			}
		}
		array[i], array[minIndex] = array[minIndex], array[i]
	}
}

// Returns the average amount of time needed to sort a list
// over 10 tries.
func singleInstance(s string, array []int) float64 {
	var runTimes []float64
	// Copies the array and starts the timer for a sorting algorithm
	for j := 0; j < 10; j++ {
		arrayToSort := make([]int, len(array))
		copy(arrayToSort, array)
		start := time.Now()
		// Determines which sorting algorithm to use
		if s == "insertion" {
			insertionSort(arrayToSort)
		} else {
			selectionSort(arrayToSort)
		}
		elapsed := time.Since(start).Seconds()
		runTimes = append(runTimes, elapsed)
	}
	avg := stat.Mean(runTimes, nil)
	return avg
}

// Creates x amount of arrays with size n and uses two sorting
// algorithms on each instance. Sends the data back in a 2D array,
// where the first slice consists of runtimes for insertion sort
// and the second slice consists of the runtimes for selection sort.
func multipleInstance(x int, n int, ch chan [][]float64) {
	var alg1Times []float64
	var alg2Times []float64
	// Runs sorting algorithms x times.
	for i := 0; i < x; i++ {
		// Generates array
		var array []int
		for a := 0; a < n; a++ {
			array = append(array, rand.Intn(n*10))
		}
		alg1Times = append(alg1Times, singleInstance("insertion", array))
		alg2Times = append(alg2Times, singleInstance("selection", array))
	}
	// Formatting data to fit in channel
	var data = [][]float64{}
	data = append(data, alg1Times)
	data = append(data, alg2Times)
	ch <- data
}

func compareEfficiency(n int) {
	// Compares the effiency of two algorithms over 1000 instances.
	// This will compute the average time to run each algorithm along
	// with the standard deviation of the times.
	var alg1Combined = []float64{}
	var alg2Combined = []float64{}
	// creates 10 go routines that each run 100 instances
	ch := make(chan [][]float64)
	for i := 0; i < 10; i++ {
		go multipleInstance(100, n, ch)
	}
	// receives data from go routines
	for j := 0; j < 10; j++ {
		data := <-ch
		alg1Data := data[0]
		alg2Data := data[1]
		alg1Combined = append(alg1Combined, alg1Data...)
		alg2Combined = append(alg2Combined, alg2Data...)
	}
	// Computing averages and standard deviations
	alg1Avg := stat.Mean(alg1Combined, nil)
	alg1StdDev := stat.StdDev(alg1Combined, nil)
	alg2Avg := stat.Mean(alg2Combined, nil)
	alg2StdDev := stat.StdDev(alg2Combined, nil)

	fmt.Printf("Insertion Sort (n = %d): %v, %v\n", n, alg1Avg, alg1StdDev)
	fmt.Printf("Selection Sort (n = %d): %v, %v\n", n, alg2Avg, alg2StdDev)
}

func main() {
	rand.Seed(int64(time.Now().Nanosecond()))
	var size = [4]int{100, 1000, 10000, 100000}
	for n := 0; n < len(size); n++ {
		compareEfficiency(size[n])
	}

}
