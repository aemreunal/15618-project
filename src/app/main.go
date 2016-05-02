package main

import (
	"fmt"
	"lockmap"
	"math/rand"
	"runtime"
	"sync"
)

const (
	NumWritesInWriteOnlyTestSmall = 1024 * 1024 * 16 // 16 M
)

type IMap interface {
	Get(k interface{}) (interface{}, bool)
	Put(k, v interface{}) interface{}
	Remove(k interface{}) (interface{}, bool)
}

type testFunc func(IMap, int, int)

func main() {
	// Set num. threads to be equal to num. of logical processors
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Test concurrentWrites
	fmt.Println("Testing concurrentWrites")
	runTest(lockmap.NewLockMap(), NumWritesInWriteOnlyTestSmall, NumWritesInWriteOnlyTestSmall, concurrentWrites)
}

/*
 * Normally distributed random number generator
 */
func getNextNormalRandom(upperLimit int) int {
	mean := upperLimit / 2
	stdDev := upperLimit / 6
	var next int
	for {
		next = rand.NormFloat64()*stdDev + mean
		if next >= 0 && next < upperLimit {
			return int(next)
		}
	}
}

/*
 * Generic wrapper to run test functions
 */
func runTest(m IMap, numWrites int, numKeys int, testToRun testFunc) {
	// Create waitgroup to wait for goroutines to finish
	var wg sync.WaitGroup
	wg.Add(runtime.NumCPU())

	// Create goroutines
	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			defer wg.Done()
			// Run the test function
			testToRun(m, numWrites, numKeys)
		}()
	}

	// Wait for goroutines to finish
	wg.Wait()
}

/*
 * 1.1
 */
func concurrentWrites(m IMap, numWrites int, numKeys int) {
	for i := 0; i < numWrites; i++ {
		k := rand.Int()
		v := k
		m.Put(k, v)
	}
}

/*
 * 1.2
 */
func concurrentWritesNormalDist(m IMap, numWrites int, numKeys int) {
	for i := 0; i < numWrites; i++ {
		k := getNextNormalRandom(numWrites)
		v := k
		m.Put(k, v)
	}
}

// /*
//  * 1.2
//  */
// func func_name(m IMap, numWrites int, numKeys int) {

// }
