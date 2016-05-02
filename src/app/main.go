package main

import (
	"fmt"
	"lockmap"
	"math/rand"
	"runtime"
	"sync"
)

const (
	numWritesInWriteOnlyTestSmall     = 1024 * 1024 * 16 // 16 M
	numWritesInRWTestSmall            = 1024 * 1024 * 16 // 16 M
	numReadsInReadOnlyTestSmall       = 1024 * 1024 * 16 // 16 M
	numIterationInConcurrentReadWrite = 10 * 1024 * 16
	numWriteDeleteIter                = 5
	numKeysInBigMap                   = 1024 * 1024 * 16 // 16 M
	numKeysInSmallMap                 = 1024 * 16
	writeRatioHigh                    = 1000
	writeRatioLow                     = 2
)

type iMap interface {
	Get(k interface{}) (interface{}, bool)
	Put(k, v interface{}) interface{}
	Remove(k interface{}) (interface{}, bool)
}

type testFunc func(iMap, int, int)

func main() {
	// Set num. threads to be equal to num. of logical processors
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Test concurrentWrites
	fmt.Println("Testing concurrentWrites")
	runTest(lockmap.NewLockMap(), numWritesInWriteOnlyTestSmall, numWritesInWriteOnlyTestSmall, concurrentWrites)
}

/*
 * Initializes a map
 */
func initializeMap(nKeys int, m iMap) {
	for i := 0; i < nKeys; i++ {
		k := i
		v := fmt.Sprintf("%12d", k)
		m.Put(k, v)
	}
}

/*
 * Normally distributed random number generator
 */
func getNextNormalRandom(upperLimit int) int {
	mean := float64(upperLimit / 2)
	stdDev := float64(upperLimit / 6)
	var next float64
	for {
		next = rand.NormFloat64()*stdDev + mean
		if next >= 0 && int(next) < upperLimit {
			return int(next)
		}
	}
}

/*
 * Generic wrapper to run test functions
 */
func runTest(m iMap, numWrites int, numKeys int, testToRun testFunc) {
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
func concurrentWrites(m iMap, numWrites int, numKeys int) {
	for i := 0; i < numWrites; i++ {
		k := rand.Int()
		v := fmt.Sprintf("%12d", k)
		m.Put(k, v)
	}
}

/*
 * 1.2
 */
func concurrentWritesNormalDist(m iMap, numWrites int, numKeys int) {
	for i := 0; i < numWrites; i++ {
		k := getNextNormalRandom(numWrites)
		v := fmt.Sprintf("%12d", k)
		m.Put(k, v)
	}
}

/*
 * 1.3
 */
func concurrentWritesSequential(m iMap, numWrites int, numKeys int) {
	currentKey := 0
	for i := 0; i < numWrites; i++ {
		k := currentKey
		currentKey = (currentKey + 1) % numWrites
		v := fmt.Sprintf("%12d", k)
		m.Put(k, v)
	}
}

/*
 * 2.1
 */
func lotsWritesFewReads(m iMap, numWrites int, numKeys int) {
	for i := 0; i < numWrites; i++ {
		if i > 0 && i%writeRatioHigh == 0 {
			/* Do a read */
			k := rand.Int()
			v, ok := m.Get(k)
			if ok {
				expectedV := fmt.Sprintf("%12d", k)
				if v != expectedV {
					fmt.Errorf("Wrong value for key", k, ". Expect ", expectedV, ". Got ", v)
				}
			}
		} else {
			/* Do write */
			k := rand.Int()
			v := fmt.Sprintf("%12d", k)
			m.Put(k, v)
		}
	}
}

/*
 * 2.2
 */
func lotsWritesFewReadsNormalDist(m iMap, numWrites int, numKeys int) {
	for i := 0; i < numWrites; i++ {
		if i > 0 && i%writeRatioHigh == 0 {
			/* Do a read */
			k := getNextNormalRandom(numWrites)
			v, ok := m.Get(k)
			if ok {
				expectedV := fmt.Sprintf("%12d", k)
				if v != expectedV {
					fmt.Errorf("Wrong value for key", k, ". Expect ", expectedV, ". Got ", v)
				}
			}
		} else {
			/* Do write */
			k := getNextNormalRandom(numWrites)
			v := fmt.Sprintf("%12d", k)
			m.Put(k, v)
		}
	}
}

/*
 * 2.3
 */
func lotsWritesFewReadsSequential(m iMap, numWrites int, numKeys int) {
	currentWriteKey := 0
	currentReadKey := 0
	for i := 0; i < numWrites; i++ {
		if i > 0 && i%writeRatioHigh == 0 {
			/* Do a read */
			k := currentReadKey
			currentReadKey = (currentReadKey + 1) % numWrites
			v, ok := m.Get(k)
			if ok {
				expectedV := fmt.Sprintf("%12d", k)
				if v != expectedV {
					fmt.Errorf("Wrong value for key", k, ". Expect ", expectedV, ". Got ", v)
				}
			}
		} else {
			/* Do write */
			k := currentWriteKey
			currentWriteKey = (currentWriteKey + 1) % numWrites
			v := fmt.Sprintf("%12d", k)
			m.Put(k, v)
		}
	}
}

/*
 * 3.1
 */
func lotsWritesLotsReads(m iMap, numWrites int, numKeys int) {
	for i := 0; i < numWrites; i++ {
		if i > 0 && i%writeRatioLow == 0 {
			/* Do a read */
			k := rand.Int()
			v, ok := m.Get(k)
			if ok {
				expectedV := fmt.Sprintf("%12d", k)
				if v != expectedV {
					fmt.Errorf("Wrong value for key", k, ". Expect ", expectedV, ". Got ", v)
				}
			}
		} else {
			/* Do write */
			k := rand.Int()
			v := fmt.Sprintf("%12d", k)
			m.Put(k, v)
		}
	}
}

/*
 * 3.2
 */
func lotsWritesLotsReadsNormalDist(m iMap, numWrites int, numKeys int) {
	for i := 0; i < numWrites; i++ {
		if i > 0 && i%writeRatioLow == 0 {
			/* Do a read */
			k := getNextNormalRandom(numWrites)
			v, ok := m.Get(k)
			if ok {
				expectedV := fmt.Sprintf("%12d", k)

				if v != expectedV {
					fmt.Errorf("Wrong value for key", k, ". Expect ", expectedV, ". Got ", v)
				}
			}
		} else {
			/* Do write */
			k := getNextNormalRandom(numWrites)
			v := fmt.Sprintf("%12d", k)
			m.Put(k, v)
		}
	}
}

/*
 * 3.3
 */
func lotsWritesLotsReadsSequential(m iMap, numWrites int, numKeys int) {
	currentKey := 0
	for i := 0; i < numWrites; i++ {
		/* Write if i is even, read if i is odd */
		if i%2 == 0 {
			/* Do a read */
			k := currentKey
			v, ok := m.Get(k)
			if ok {
				expectedV := fmt.Sprintf("%12d", k)
				if v != expectedV {
					fmt.Errorf("Wrong value for key", k, ". Expect ", expectedV, ". Got ", v)
				}
			}
		} else {
			/* Do a write */
			k := currentKey
			currentKey = (currentKey + 1) % numWrites
			v := fmt.Sprintf("%12d", k)
			m.Put(k, v)
		}
	}
}

// /*
//  * 1.
//  */
// func func_name(m IMap, numWrites int, numKeys int) {

// }
