package main

import (
	"fmt"
	"lockmap"
	"math/rand"
	"runtime"
	"sync"
)

/*
 * -------------------------------------------------
 * Types & Constants
 * -------------------------------------------------
 */

const (
	mapTypeConcurrentMap              = "chinese"
	mapTypeGotomicMap                 = "gotomic"
	mapTypeLockMap                    = "lock"
	mapTypeParallelMap                = "parallel"
	mapTypeRWLockMap                  = "rwlock"
	numIterationInConcurrentReadWrite = 10 * 1024 * 16
	numKeysInBigMap                   = 1024 * 1024 * 16       // 16 M
	numKeysInLargeMap                 = 1024 * 1024 * 1024 * 2 // 2 G
	numKeysInSmallMap                 = 1024 * 16
	numReadsInReadOnlyTestLarge       = 1024 * 1024 * 1024 * 2 // 2 G
	numReadsInReadOnlyTestSmall       = 1024 * 1024 * 16       // 16 M
	numWriteDeleteIter                = 5
	numWritesInRWTestLarge            = 1024 * 1024 * 1024 * 2 // 2 G
	numWritesInRWTestSmall            = 1024 * 1024 * 16       // 16 M
	numWritesInWriteOnlyTestLarge     = 1024 * 1024 * 1024 * 2 // 2 G
	numWritesInWriteOnlyTestSmall     = 1024 * 1024 * 16       // 16 M
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
 * -------------------------------------------------
 * Helper Goroutines
 * -------------------------------------------------
 */

func writer(do chan bool, done chan bool, m iMap, numKeys int, numWrites int) {
	<-do
	for i := 0; i < numWrites; i++ {
		k := rand.Intn(numKeys)
		v := fmt.Sprintf("%12d", k)
		m.Put(k, v)
	}
	done <- true
}

func reader(do chan bool, done chan bool, m iMap, numKeys int, numReads int) {
	<-do
	for i := 0; i < numReads; i++ {
		k := rand.Intn(numKeys)
		v, ok := m.Get(k)
		if ok {
			expectedV := fmt.Sprintf("%12d", k)
			if v != expectedV {
				fmt.Errorf("Wrong value for key", k, ". Expect ", expectedV, ". Got ", v)
			}
		} else {
			fmt.Errorf("Failed to get key ", k)
		}
	}
	done <- true
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

/*
 * 4.1
 */
func lotsReads(m iMap, numReads int, numKeys int) {
	initializeMap(numKeys, m)
	for i := 0; i < numReads; i++ {
		k := rand.Intn(numKeys)
		v, ok := m.Get(k)
		if ok {
			expectedV := fmt.Sprintf("%12d", k)
			if v != expectedV {
				fmt.Errorf("Wrong value for key", k, ". Expect ", expectedV, ". Got ", v)
			}
		} else {
			fmt.Errorf("Failed to get key ", k)
		}
	}
}

/*
 * 4.2
 */
func lotsReadsNormalDist(m iMap, numReads int, numKeys int) {
	initializeMap(numKeys, m)
	for i := 0; i < numReads; i++ {
		k := getNextNormalRandom(numKeys)
		v, ok := m.Get(k)
		if ok {
			expectedV := fmt.Sprintf("%12d", k)
			if v != expectedV {
				fmt.Errorf("Wrong value for key", k, ". Expect ", expectedV, ". Got ", v)
			}
		} else {
			fmt.Errorf("Failed to get key ", k)
		}
	}
}

/*
 * 4.3
 */
func lotsReadsSequential(m iMap, numReads int, numKeys int) {
	currentKey := 0
	initializeMap(numKeys, m)
	for i := 0; i < numReads; i++ {
		k := currentKey
		currentKey = (currentKey + 1) % numKeys
		v, ok := m.Get(k)
		if ok {
			expectedV := fmt.Sprintf("%12d", k)
			if v != expectedV {
				fmt.Errorf("Wrong value for key", k, ". Expect ", expectedV, ". Got ", v)
			}
		} else {
			fmt.Errorf("Failed to get key ", k)
		}
	}
}

/*
 * 8/9/10
 */
func concurrentWriterReaders(m iMap, numWriters int, numReaders int) {
	do := make(chan bool)
	done := make(chan bool)
	initializeMap(numKeysInSmallMap, m)

	/* Start writers */
	for i := 0; i < numWriters; i++ {
		go writer(do, done, m, numKeysInSmallMap, numIterationInConcurrentReadWrite)
	}
	/* Start readers */
	for i := 0; i < numReaders; i++ {
		go reader(do, done, m, numKeysInSmallMap, numIterationInConcurrentReadWrite)
	}

	close(do)
	/* Readers/Writers finish */
	for i := 0; i < numWriters+numReaders; i++ {
		<-done
	}
}

/*
 * 11
 */
func concurrentWriteDeleteWrite(m iMap, numWrites int, numKeys int) {
	for i := 0; i < numWriteDeleteIter; i++ {
		for i := 0; i < numKeysInSmallMap; i++ {
			k := i
			v := fmt.Sprintf("%12d", k)
			m.Put(k, v)
		}
		for i := 0; i < numKeysInSmallMap; i++ {
			m.Remove(i)
		}
	}
}
