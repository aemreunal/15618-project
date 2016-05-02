package benchmark

import (
	"concurrent"
	"fmt"
	"gotomic"
	"lockmap"
	"math/rand"
	"nativemap"
	"pmap"
	"runtime"
	"rwlockmap"
	"testing"
	"time"
)

/* A key is a int which is 4 bytes. A value is a string which is 12 bytes.
   So each hash entry is 16 bytes. */

const (
	NumWritesInWriteOnlyTestSmall     = 1024 * 1024 * 16 // 16 M
	NumWritesInRWTestSmall            = 1024 * 1024 * 16 // 16 M
	NumReadsInReadOnlyTestSmall       = 1024 * 1024 * 16 // 16 M
	NumIterationInConcurrentReadWrite = 10 * 1024 * 16
	NumWriteDeleteIter                = 5
	NumKeysInBigMap                   = 1024 * 1024 * 16 // 16 M
	NumKeysInSmallMap                 = 1024 * 16
	WriteRatioHigh                    = 1000
	WriteRatioLow                     = 2
)

type IMap interface {
	Get(k interface{}) (interface{}, bool)
	Put(k, v interface{}) interface{}
	Remove(k interface{}) (interface{}, bool)
}

func InitializeMap(nKeys int, m IMap) {
	for i := 0; i < nKeys; i++ {
		k := i
		v := fmt.Sprintf("%12d", k)
		m.Put(k, v)
	}
}

func benchmarkPutGetBasic(m IMap, b *testing.B) {
	for i := 0; i < b.N; i++ {
		k := i
		v := fmt.Sprintf("%12d", k)
		m.Put(k, v)
		v1, ok := m.Get(k)
		if !ok {
			b.Error("Failed to get key ", k)
		}
		if v1 != v {
			b.Error("Should be ", v, ". Got ", v1)
		}
	}
}

/*
 * 1.1. Lots of writes to uniformly random keys, no reads, fits to memory ->
 * helps test cache misses for those keys
 */
func benchmarkConcurrentWrites(m IMap, b *testing.B, numWrites int) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			rand.Seed(time.Now().UTC().UnixNano())
			for i := 0; i < numWrites; i++ {
				k := rand.Int()
				v := fmt.Sprintf("%12d", k)
				m.Put(k, v)
			}
		}
	})
}

/*
 * 1.2. Lots of writes to normally random keys, no reads, fits to memory ->
 * helps test cache misses for those keys
 */
func benchmarkConcurrentWritesNormalDist(m IMap, b *testing.B, numWrites int) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			rand.Seed(time.Now().UTC().UnixNano())
			for i := 0; i < numWrites; i++ {
				k := getNextNormalRandom(numWrites)
				v := fmt.Sprintf("%12d", k)
				m.Put(k, v)
			}
		}
	})
}

/*
 * 1.3. Lots of writes to sequential keys, no reads, fits to memory ->
 * helps test cache misses for those keys
 */
func benchmarkConcurrentWritesSequential(m IMap, b *testing.B, numWrites int) {
	currentKey := 0
	runtime.GOMAXPROCS(runtime.NumCPU())
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			rand.Seed(time.Now().UTC().UnixNano())
			for i := 0; i < numWrites; i++ {
				k := currentKey
				currentKey = (currentKey + 1) % numWrites
				v := fmt.Sprintf("%12d", k)
				m.Put(k, v)
			}
		}
	})
}

/*
 * 2.1. Lots of writes to uniformly random keys, few reads, fits to memory
 */
func benchmarkLotsWritesFewReads(m IMap, b *testing.B, numWrites int) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			rand.Seed(time.Now().UTC().UnixNano())
			for i := 0; i < numWrites; i++ {
				/* Do a read */
				if i > 0 && i%WriteRatioHigh == 0 {
					k := rand.Int()
					v, ok := m.Get(k)
					if ok {
						expectedV := fmt.Sprintf("%12d", k)
						if v != expectedV {
							b.Error("Wrong value for key", k, ". Expect ", expectedV, ". Got ", v)
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
	})
}

/*
 * 2.2. Lots of writes to normally random keys, few reads, fits to memory
 */
func benchmarkLotsWritesFewReadsNormalDist(m IMap, b *testing.B, numWrites int) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			rand.Seed(time.Now().UTC().UnixNano())
			for i := 0; i < numWrites; i++ {
				/* Do a read */
				if i > 0 && i%WriteRatioHigh == 0 {
					k := getNextNormalRandom(numWrites)
					v, ok := m.Get(k)
					if ok {
						expectedV := fmt.Sprintf("%12d", k)
						if v != expectedV {
							b.Error("Wrong value for key", k, ". Expect ", expectedV, ". Got ", v)
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
	})
}

/*
 * 2.3. Lots of writes to sequential keys, few reads, fits to memory
 */
func benchmarkLotsWritesFewReadsSequential(m IMap, b *testing.B, numWrites int) {
	currentWriteKey := 0
	currentReadKey := 0
	runtime.GOMAXPROCS(runtime.NumCPU())
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			rand.Seed(time.Now().UTC().UnixNano())
			for i := 0; i < numWrites; i++ {
				/* Do a read */
				if i > 0 && i%WriteRatioHigh == 0 {
					k := currentReadKey
					currentReadKey = (currentReadKey + 1) % numWrites
					v, ok := m.Get(k)
					if ok {
						expectedV := fmt.Sprintf("%12d", k)
						if v != expectedV {
							b.Error("Wrong value for key", k, ". Expect ", expectedV, ". Got ", v)
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
	})
}

/*
 * 3.1. Lots of writes to uniformly random keys, lots of uniformly random reads,
 * fits into memory
 */
func benchmarkLotsWritesLotsReads(m IMap, b *testing.B, numWrites int) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			rand.Seed(time.Now().UTC().UnixNano())
			for i := 0; i < numWrites; i++ {
				/* Do a read */
				if i > 0 && i%WriteRatioLow == 0 {
					k := rand.Int()
					v, ok := m.Get(k)
					if ok {
						expectedV := fmt.Sprintf("%12d", k)
						if v != expectedV {
							b.Error("Wrong value for key", k, ". Expect ", expectedV, ". Got ", v)
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
	})
}

/*
 * 3.2. Lots of writes to normally distributed random keys, lots of normally
 * distributed random reads, fits into memory
 */
func benchmarkLotsWritesLotsReadsNormalDist(m IMap, b *testing.B, numWrites int) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			rand.Seed(time.Now().UTC().UnixNano())
			for i := 0; i < numWrites; i++ {
				/* Do a read */
				if i > 0 && i%WriteRatioLow == 0 {
					k := getNextNormalRandom(numWrites)
					v, ok := m.Get(k)
					if ok {
						expectedV := fmt.Sprintf("%12d", k)

						if v != expectedV {
							b.Error("Wrong value for key", k, ". Expect ", expectedV, ". Got ", v)
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
	})
}

/*
 * 3.3. Interleaved sequential writes and reads
 */
func benchmarkLotsWritesLotsReadsSequential(m IMap, b *testing.B, numWrites int) {
	currentKey := 0
	runtime.GOMAXPROCS(runtime.NumCPU())
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			rand.Seed(time.Now().UTC().UnixNano())
			for i := 0; i < numWrites; i++ {
				/* Write if i is even, read if i is odd */
				if i%2 == 0 {
					/* Do a read */
					k := currentKey
					v, ok := m.Get(k)
					if ok {
						expectedV := fmt.Sprintf("%12d", k)
						if v != expectedV {
							b.Error("Wrong value for key", k, ". Expect ", expectedV, ". Got ", v)
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
	})
}

/*
 *  4.1. Initialize a large table (fitting into memory)
 *     (do not test the initialization part),
 *	   then lots of uniformly random reads ->
 *     cache behavior when reading from an unchanging table
 */
func benchmarkLotsReads(m IMap, b *testing.B, numKeys, numReads int) {
	/* Initialize the map */
	InitializeMap(numKeys, m)
	b.ResetTimer()
	runtime.GOMAXPROCS(runtime.NumCPU())
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			rand.Seed(time.Now().UTC().UnixNano())
			for i := 0; i < numReads; i++ {
				k := rand.Intn(numKeys)
				v, ok := m.Get(k)
				if ok {
					expectedV := fmt.Sprintf("%12d", k)
					if v != expectedV {
						b.Error("Wrong value for key", k, ". Expect ", expectedV, ". Got ", v)
					}
				} else {
					b.Error("Failed to get key ", k)
				}
			}
		}
	})
}

/*
 *  4.2. Initialize a large table (fitting into memory)
 *     (do not test the initialization part),
 *	   then lots of normally distributed random reads ->
 *     cache behavior when reading from an unchanging table
 */
func benchmarkLotsReadsNormalDist(m IMap, b *testing.B, numKeys, numReads int) {
	/* Initialize the map */
	InitializeMap(numKeys, m)
	b.ResetTimer()
	runtime.GOMAXPROCS(runtime.NumCPU())
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			rand.Seed(time.Now().UTC().UnixNano())
			for i := 0; i < numReads; i++ {
				k := getNextNormalRandom(numKeys)
				v, ok := m.Get(k)
				if ok {
					expectedV := fmt.Sprintf("%12d", k)
					if v != expectedV {
						b.Error("Wrong value for key", k, ". Expect ", expectedV, ". Got ", v)
					}
				} else {
					b.Error("Failed to get key ", k)
				}
			}
		}
	})
}

func getNextNormalRandom(numKeys int) int {
	mean := numKeys / 2
	stdDev := numKeys / 6
	var next int
	for {
		next = int(rand.NormFloat64()*float64(stdDev) + float64(mean))
		if next >= 0 && next < numKeys {
			return next
		}
	}
}

/*
 *  4.3. Initialize a large table (fitting into memory)
 *     (do not test the initialization part),
 *	   then lots of sequential reads ->
 *     cache behavior when reading from an unchanging table
 */
func benchmarkLotsReadsSequential(m IMap, b *testing.B, numKeys, numReads int) {
	currentKey := 0
	/* Initialize the map */
	InitializeMap(numKeys, m)
	b.ResetTimer()
	runtime.GOMAXPROCS(runtime.NumCPU())
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			rand.Seed(time.Now().UTC().UnixNano())
			for i := 0; i < numReads; i++ {
				k := currentKey
				currentKey = (currentKey + 1) % numKeys
				v, ok := m.Get(k)
				if ok {
					expectedV := fmt.Sprintf("%12d", k)
					if v != expectedV {
						b.Error("Wrong value for key", k, ". Expect ", expectedV, ". Got ", v)
					}
				} else {
					b.Error("Failed to get key ", k)
				}
			}
		}
	})
}

func Writer(do, done chan bool, m IMap, nKeys, numWrites int, b *testing.B) {
	<-do
	for i := 0; i < numWrites; i++ {
		k := rand.Intn(nKeys)
		v := fmt.Sprintf("%12d", k)
		m.Put(k, v)
	}
	done <- true
}

func Reader(do, done chan bool, m IMap, nKeys, numReads int, b *testing.B) {
	<-do
	for i := 0; i < numReads; i++ {
		k := rand.Intn(nKeys)
		v, ok := m.Get(k)
		if ok {
			expectedV := fmt.Sprintf("%12d", k)
			if v != expectedV {
				b.Error("Wrong value for key", k, ". Expect ", expectedV, ". Got ", v)
			}
		} else {
			b.Error("Failed to get key ", k)
		}
	}
	done <- true
}

/*
 *  8/9/10. n1 concurrent writers, n2 readers
 */
func benchmarkConcurrentWriterReaders(numWriters, numReaders int, m IMap, b *testing.B) {
	/* Set the maximum number of CPUs that can be executing simultaneously */
	runtime.GOMAXPROCS(runtime.NumCPU())
	do := make(chan bool)
	done := make(chan bool)
	InitializeMap(NumKeysInSmallMap, m)

	/* Start writers */
	for i := 0; i < numWriters; i++ {
		go Writer(do, done, m, NumKeysInSmallMap, NumIterationInConcurrentReadWrite, b)
	}
	/* Start readers */
	for i := 0; i < numReaders; i++ {
		go Reader(do, done, m, NumKeysInSmallMap, NumIterationInConcurrentReadWrite, b)
	}
	close(do)
	b.ResetTimer()
	/* Readers/Writers finish */
	for i := 0; i < numWriters+numReaders; i++ {
		<-done
	}
}

/*
 *  11. Write a lot, then delete a lot, then write a lot, etc. ->
 *      helps test resize behavior (assuming the implementation properly frees
 *      the memory and resizes the data structures)
 */
func benchmarkConcurrentWriteDeleteWrite(m IMap, b *testing.B) {
	b.ResetTimer()
	runtime.GOMAXPROCS(runtime.NumCPU())
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			rand.Seed(time.Now().UTC().UnixNano())
			for i := 0; i < NumWriteDeleteIter; i++ {
				for i := 0; i < NumKeysInSmallMap; i++ {
					k := i
					v := fmt.Sprintf("%12d", k)
					m.Put(k, v)
				}
				for i := 0; i < NumKeysInSmallMap; i++ {
					m.Remove(i)
				}
			}
		}
	})
}

/* 0. ==========================Basic Put/Get Test============================*/
func BenchmarkNativeMapPutGetBasic(b *testing.B) {
	benchmarkPutGetBasic(nativemap.NewNativeMap(), b)
}

func BenchmarkLockMapPutGetBasic(b *testing.B) {
	benchmarkPutGetBasic(lockmap.NewLockMap(), b)
}

func BenchmarkRWLockMapPutGetBasic(b *testing.B) {
	benchmarkPutGetBasic(rwlockmap.NewRWLockMap(), b)
}

func BenchmarkParallelMapPutGetBasic(b *testing.B) {
	benchmarkPutGetBasic(pmap.NewParallelMap(), b)
}

func BenchmarkGotomicMapPutGetBasic(b *testing.B) {
	benchmarkPutGetBasic(gotomic.NewGotomicMap(), b)
}

func BenchmarkConcurrentMapPutGetBasic(b *testing.B) {
	benchmarkPutGetBasic(concurrent.NewConcurrentMap(), b)
}

/* 1. =======================Lots of concurrent writes======================= */
func BenchmarkLockMapLotsWrite(b *testing.B) {
	benchmarkConcurrentWrites(lockmap.NewLockMap(), b, NumWritesInWriteOnlyTestSmall)
}

func BenchmarkRWLockMapLotsWrite(b *testing.B) {
	benchmarkConcurrentWrites(rwlockmap.NewRWLockMap(), b, NumWritesInWriteOnlyTestSmall)
}

func BenchmarkParallelMapLotsWrite(b *testing.B) {
	benchmarkConcurrentWrites(pmap.NewParallelMap(), b, NumWritesInWriteOnlyTestSmall)
}

func BenchmarkGotomicMapLotsWrite(b *testing.B) {
	benchmarkConcurrentWrites(gotomic.NewGotomicMap(), b, NumWritesInWriteOnlyTestSmall)
}

func BenchmarkConcurrentMapLotsWrite(b *testing.B) {
	benchmarkConcurrentWrites(concurrent.NewConcurrentMap(), b, NumWritesInWriteOnlyTestSmall)
}

/* 2. ==================Lots of concurrent writes, few reads================= */
func BenchmarkLockMapLotsWritesFewReads(b *testing.B) {
	benchmarkLotsWritesFewReads(lockmap.NewLockMap(), b, NumWritesInRWTestSmall)
}

func BenchmarkRWLockMapLotsWritesFewReads(b *testing.B) {
	benchmarkLotsWritesFewReads(rwlockmap.NewRWLockMap(), b, NumWritesInRWTestSmall)
}

func BenchmarkParallelMapLotsWritesFewReads(b *testing.B) {
	benchmarkLotsWritesFewReads(pmap.NewParallelMap(), b, NumWritesInRWTestSmall)
}

func BenchmarkGotomicMapLotsWritesFewReads(b *testing.B) {
	benchmarkLotsWritesFewReads(gotomic.NewGotomicMap(), b, NumWritesInRWTestSmall)
}

func BenchmarkConcurrentMapLotsWritesFewReads(b *testing.B) {
	benchmarkLotsWritesFewReads(concurrent.NewConcurrentMap(), b, NumWritesInRWTestSmall)
}

/* 3. ================Lots of concurrent writes, lots of reads=============== */
func BenchmarkLockMapLotsWritesLotsReads(b *testing.B) {
	benchmarkLotsWritesLotsReads(lockmap.NewLockMap(), b, NumWritesInRWTestSmall)
}

func BenchmarkRWLockMapLotsWritesLotsReads(b *testing.B) {
	benchmarkLotsWritesLotsReads(lockmap.NewLockMap(), b, NumWritesInRWTestSmall)
}

func BenchmarkParallelMapLotsWritesLotsReads(b *testing.B) {
	benchmarkLotsWritesLotsReads(pmap.NewParallelMap(), b, NumWritesInRWTestSmall)
}

func BenchmarkGotomicMapLotsWritesLotsReads(b *testing.B) {
	benchmarkLotsWritesLotsReads(gotomic.NewGotomicMap(), b, NumWritesInRWTestSmall)
}

func BenchmarkConcurrentMapLotsWritesLotsReads(b *testing.B) {
	benchmarkLotsWritesLotsReads(concurrent.NewConcurrentMap(), b, NumWritesInRWTestSmall)
}

/* 4. =======================Lots of concurrent reads======================== */
func BenchmarkNativeMapLotsReads(b *testing.B) {
	benchmarkLotsReads(nativemap.NewNativeMap(), b, NumKeysInBigMap, NumReadsInReadOnlyTestSmall)
}

func BenchmarkLockMapLotsReads(b *testing.B) {
	benchmarkLotsReads(lockmap.NewLockMap(), b, NumKeysInBigMap, NumReadsInReadOnlyTestSmall)
}

func BenchmarkRWLockMapLotsReads(b *testing.B) {
	benchmarkLotsReads(rwlockmap.NewRWLockMap(), b, NumKeysInBigMap, NumReadsInReadOnlyTestSmall)
}

func BenchmarkParallelMapLotsReads(b *testing.B) {
	benchmarkLotsReads(pmap.NewParallelMap(), b, NumKeysInBigMap, NumReadsInReadOnlyTestSmall)
}

func BenchmarkGotomicMapLotsReads(b *testing.B) {
	benchmarkLotsReads(gotomic.NewGotomicMap(), b, NumKeysInBigMap, NumReadsInReadOnlyTestSmall)
}

func BenchmarkConcurrentMapLotsReads(b *testing.B) {
	benchmarkLotsReads(concurrent.NewConcurrentMap(), b, NumKeysInBigMap, NumReadsInReadOnlyTestSmall)
}

/* 8. ====================100 concurrent writers, 10 readers==================*/
func BenchmarkLockMapConcurrentWriterReaders1(b *testing.B) {
	benchmarkConcurrentWriterReaders(100, 10, lockmap.NewLockMap(), b)
}

func BenchmarkRWLockMapConcurrentWriterReaders1(b *testing.B) {
	benchmarkConcurrentWriterReaders(100, 10, rwlockmap.NewRWLockMap(), b)
}

func BenchmarkParallelMapConcurrentWriterReaders1(b *testing.B) {
	benchmarkConcurrentWriterReaders(100, 10, pmap.NewParallelMap(), b)
}

func BenchmarkGotomicMapConcurrentWriterReaders1(b *testing.B) {
	benchmarkConcurrentWriterReaders(100, 10, gotomic.NewGotomicMap(), b)
}

func BenchmarkConcurrentMapConcurrentWriterReaders1(b *testing.B) {
	benchmarkConcurrentWriterReaders(100, 10, concurrent.NewConcurrentMap(), b)
}

/* 9. ====================10 concurrent writers, 100 readers==================*/
func BenchmarkLockMapConcurrentWriterReaders2(b *testing.B) {
	benchmarkConcurrentWriterReaders(10, 100, lockmap.NewLockMap(), b)
}

func BenchmarkRWLockMapConcurrentWriterReaders2(b *testing.B) {
	benchmarkConcurrentWriterReaders(10, 100, rwlockmap.NewRWLockMap(), b)
}

func BenchmarkParallelMapConcurrentWriterReaders2(b *testing.B) {
	benchmarkConcurrentWriterReaders(10, 100, pmap.NewParallelMap(), b)
}

func BenchmarkGotomicMapConcurrentWriterReaders2(b *testing.B) {
	benchmarkConcurrentWriterReaders(10, 100, gotomic.NewGotomicMap(), b)
}

func BenchmarkConcurrentMapConcurrentWriterReaders2(b *testing.B) {
	benchmarkConcurrentWriterReaders(10, 100, concurrent.NewConcurrentMap(), b)
}

/* 10. ====================1 concurrent writers, 100 readers==================*/
func BenchmarkLockMapConcurrentWriterReaders3(b *testing.B) {
	benchmarkConcurrentWriterReaders(1, 100, lockmap.NewLockMap(), b)
}

func BenchmarkRWLockMapConcurrentWriterReaders3(b *testing.B) {
	benchmarkConcurrentWriterReaders(1, 100, rwlockmap.NewRWLockMap(), b)
}

func BenchmarkParallelMapConcurrentWriterReaders3(b *testing.B) {
	benchmarkConcurrentWriterReaders(1, 100, pmap.NewParallelMap(), b)
}

func BenchmarkGotomicMapConcurrentWriterReaders3(b *testing.B) {
	benchmarkConcurrentWriterReaders(1, 100, gotomic.NewGotomicMap(), b)
}

func BenchmarkConcurrentMapConcurrentWriterReaders3(b *testing.B) {
	benchmarkConcurrentWriterReaders(1, 100, concurrent.NewConcurrentMap(), b)
}

/* 11. =======================Write, delete, write=========================== */
func BenchmarkLockMapWriteDeleteWrite(b *testing.B) {
	benchmarkConcurrentWriteDeleteWrite(lockmap.NewLockMap(), b)
}

func BenchmarkRWLockMapWriteDeleteWrite(b *testing.B) {
	benchmarkConcurrentWriteDeleteWrite(rwlockmap.NewRWLockMap(), b)
}

func BenchmarkParallelMapWriteDeleteWrite(b *testing.B) {
	benchmarkConcurrentWriteDeleteWrite(pmap.NewParallelMap(), b)
}

func BenchmarkGotomicMapWriteDeleteWrite(b *testing.B) {
	benchmarkConcurrentWriteDeleteWrite(gotomic.NewGotomicMap(), b)
}

func BenchmarkConcurrentMapWriteDeleteWrite(b *testing.B) {
	benchmarkConcurrentWriteDeleteWrite(concurrent.NewConcurrentMap(), b)
}
