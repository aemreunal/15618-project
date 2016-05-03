package main

import (
	"concurrent"
	"fmt"
	"gotomic"
	"lockmap"
	"math/rand"
	"os"
	"pmap"
	"runtime"
	"rwlockmap"
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

/*
 * -------------------------------------------------
 * Main
 * -------------------------------------------------
 */

func main() {
	// Set num. threads to be equal to num. of logical processors
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Parse arguments
	if len(os.Args) == 3 {
		// Passed args:
		// 1) Program name
		// 2) Test number
		// 3) Map type
		testNum := os.Args[1]
		mapType := os.Args[2]
		createAndRunTest(testNum, mapType)
	} else {
		printHelpText()
		os.Exit(0)
	}
}

func createAndRunTest(testNum string, mapType string) {
	// Create test map object
	var testMap iMap
	switch mapType {
	case mapTypeConcurrentMap:
		testMap = concurrent.NewConcurrentMap()
	case mapTypeGotomicMap:
		testMap = gotomic.NewGotomicMap()
	case mapTypeLockMap:
		testMap = lockmap.NewLockMap()
	case mapTypeParallelMap:
		testMap = pmap.NewParallelMap()
	case mapTypeRWLockMap:
		testMap = rwlockmap.NewRWLockMap()
	default:
		fmt.Errorf("Invalid map type entered")
		os.Exit(-1)
	}

	switch testNum {
	case "1":
		fmt.Println("Testing lots of concurrent writes")
		runTestSingleRW(testMap, numWritesInWriteOnlyTestSmall, numWritesInWriteOnlyTestSmall, concurrentWrites)
	case "2":
		fmt.Println("Testing lots of concurrent writes, few concurrent reads")
		runTestSingleRW(testMap, numWritesInRWTestSmall, numWritesInRWTestSmall, lotsWritesFewReads)
	case "3":
		fmt.Println("Testing lots of concurrent writes, lots of concurrent reads")
		runTestSingleRW(testMap, numWritesInRWTestSmall, numWritesInRWTestSmall, lotsWritesLotsReads)
	case "4":
		fmt.Println("Testing lots of concurrent reads")
		runTestSingleRW(testMap, numReadsInReadOnlyTestSmall, numKeysInBigMap, lotsReads)
	case "5.1":
		fmt.Println("Testing 1-lots of concurrent writes in large data set.")
		runTestSingleRW(testMap, numWritesInWriteOnlyTestLarge, numWritesInWriteOnlyTestLarge, concurrentWrites)
	case "5.2":
		fmt.Println("Testing 2-concurrent writes, few concurrent reads in large data set.")
		runTestSingleRW(testMap, numWritesInRWTestLarge, numWritesInRWTestLarge, lotsWritesFewReads)
	case "5.3":
		fmt.Println("Testing 3-concurrent writes, lots of concurrent reads in large data set.")
		runTestSingleRW(testMap, numWritesInRWTestLarge, numWritesInRWTestLarge, lotsWritesLotsReads)
	case "5.4":
		fmt.Println("Testing 4-lots of concurrent reads in large data set.")
		runTestSingleRW(testMap, numReadsInReadOnlyTestLarge, numKeysInLargeMap, lotsReads)
	case "6.1":
		fmt.Println("Testing 1-lots of concurrent writes in normal distribution keys.")
		runTestSingleRW(testMap, numWritesInWriteOnlyTestSmall, numWritesInWriteOnlyTestSmall, concurrentWritesNormalDist)
	case "6.2":
		fmt.Println("Testing 2-concurrent writes, few concurrent reads in normal distribution keys.")
		runTestSingleRW(testMap, numWritesInRWTestSmall, numWritesInRWTestSmall, lotsWritesFewReadsNormalDist)
	case "6.3":
		fmt.Println("Testing 3-concurrent writes, lots of concurrent reads in normal distribution keys.")
		runTestSingleRW(testMap, numWritesInRWTestSmall, numWritesInRWTestSmall, lotsWritesLotsReadsNormalDist)
	case "6.4":
		fmt.Println("Testing 4-lots of concurrent reads in normal distribution keys.")
		runTestSingleRW(testMap, numReadsInReadOnlyTestSmall, numKeysInBigMap, lotsReadsNormalDist)
	case "7.1":
		fmt.Println("Testing 1-lots of concurrent writes in sequential keys.")
		runTestSingleRW(testMap, numWritesInWriteOnlyTestSmall, numWritesInWriteOnlyTestSmall, concurrentWritesNormalDist)
	case "7.2":
		fmt.Println("Testing 2-concurrent writes, few concurrent reads in sequential keys.")
		runTestSingleRW(testMap, numWritesInRWTestSmall, numWritesInRWTestSmall, lotsWritesFewReadsNormalDist)
	case "7.3":
		fmt.Println("Testing 3-concurrent writes, lots of concurrent reads in sequential keys.")
		runTestSingleRW(testMap, numWritesInRWTestSmall, numWritesInRWTestSmall, lotsWritesLotsReadsNormalDist)
	case "7.4":
		fmt.Println("Testing 4-lots of concurrent reads in sequential keys.")
		runTestSingleRW(testMap, numReadsInReadOnlyTestSmall, numKeysInBigMap, lotsReadsNormalDist)
	case "8":
		fmt.Println("Testing 100 concurrent writers, 10 concurrent readers")
		runTestConcurrentRW(testMap, 100, 10, concurrentWriterReaders)
	case "9":
		fmt.Println("Testing 10 concurrent writers, 100 concurrent readers")
		runTestConcurrentRW(testMap, 10, 100, concurrentWriterReaders)
	case "10":
		fmt.Println("Testing 1 concurrent writers, 100 concurrent readers")
		runTestConcurrentRW(testMap, 1, 100, concurrentWriterReaders)
	case "11":
		fmt.Println("Testing resize behavior")
		// Pass 0 because values aren't used
		runTestSingleRW(testMap, 0, 0, concurrentWriteDeleteWrite)

	default:
		fmt.Errorf("Invalid test number entered")
		os.Exit(-1)
	}
}

/*
 * -------------------------------------------------
 * Helper functions
 * -------------------------------------------------
 */

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
func runTest(m iMap, arg1 int, arg2 int, testToRun testFunc) {
	// Create waitgroup to wait for goroutines to finish
	var wg sync.WaitGroup
	wg.Add(runtime.NumCPU())

	// Create goroutines
	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			defer wg.Done()
			// Run the test function
			testToRun(m, arg1, arg2)
		}()
	}

	// Wait for goroutines to finish
	wg.Wait()
}

/*
 * Wrapper to run test functions with a single reader/writer per core
 */
func runTestSingleRW(m iMap, numIterations int, numKeys int, testToRun testFunc) {
	runTest(m, numIterations, numKeys, testToRun)
}

/*
 * Wrapper to run test functions with concurrent readers and writers per core
 */
func runTestConcurrentRW(m iMap, numWriters int, numReaders int, testToRun testFunc) {
	runTest(m, numWriters, numReaders, testToRun)
}

/*
 * Prints console help text
 */
func printHelpText() {
	fmt.Println("Usage: ./app <test_num> <map_type>")
	fmt.Println("Map types:")
	fmt.Println("\tchinese")
	fmt.Println("\tgotomic")
	fmt.Println("\tlock")
	fmt.Println("\tparallel")
	fmt.Println("\trwlock")
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
