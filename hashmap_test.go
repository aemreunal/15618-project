package benchmark

import (
	"math/rand"
	"runtime"
	"testing"
	"time"
	"nativemap"
	"lockmap"
	"rwlockmap"
	"pmap"
	"gotomic"
)

const (
	StrKeyLen                         = 8
	NumWritesInWriteOnlyTest          = 1024 * 1024 * 16 // 16 M
	NumWritesInRWTest                 = 1024 * 1024 * 16 // 16 M
	NumReadsInReadOnlyTest            = 1024 * 1024 * 16 // 16 M
	NumIterationInConcurrentReadWrite = 1024 * 16
	NumWriteDeleteIter                = 1
	NumKeysInBigMap                   = 1024 * 1024 * 16
	NumKeysInSmallMap                 = 1024 * 16
	WriteRatioHigh                    = 1000
	WriteRatioLow                     = 2
	chars                             = "0123456789ABCDEFGHIJKLMNOPQRSTUVXYZabcdefghijklmnopqrstuvwxyz"
)

type IMap interface {
	Get(k interface{}) (interface{}, bool)
	Put(k, v interface{}) interface{}
	Remove(k interface{}) (interface{}, bool)
}

/* Generate a random string of strlen */
func RandomString(strlen int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

func InitializeMap(nKeys int64, m IMap) {
	for i := int64(0); i < nKeys; i++ {
		m.Put(i, i)
	}
}

func benchmarkPutGetBasic(m IMap, b *testing.B) {
	for i := 0; i < b.N; i++ {
		m.Put(i, i)
		j, ok := m.Get(i)
		if !ok {
			b.Error("Failed to get key ", i)
		}
		if j != i {
			b.Error("Should be ", i, ". Got ", j)
		}
	}
}

/* 1. Lots of writes to uniformly random keys, no reads, fits to memory ->
helps test cache misses for those keys
*/
func benchmarkConcurrentWrites(m IMap, b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			rand.Seed(time.Now().UTC().UnixNano())
			for i := 0; i < NumWritesInWriteOnlyTest; i++ {
				k := rand.Int63()
				v := k
				m.Put(k, v)
			}
		}
	})
}

/*
2. Lots of writes to uniformly random keys, few reads, fits to memory
*/
func benchmarkLotsWritesFewReads(m IMap, b *testing.B) {
	keys := []int64{}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			rand.Seed(time.Now().UTC().UnixNano())
			for i := 0; i < NumWritesInRWTest; i++ {
				/* Do a read */
				if i > 0 && i%WriteRatioHigh == 0 {
					idx := rand.Int63n(int64(len(keys)))
					v, ok := m.Get(keys[idx])
					if ok {
						if v != keys[idx] {
							b.Error("Wrong value for key", keys[idx], ". Expect ", keys[idx], ". Got ", v)
						}
					} else {
						b.Error("Failed to get key ", keys[idx])
					}
				} else {
					/* Do write */
					k := rand.Int63()
					v := k
					keys = append(keys, k)
					m.Put(k, v)
				}
			}
		}
	})
}

/*
3. Lots of writes to uniformly random keys, lots of uniformly random reads,
fits into memory
*/

func benchmarkLotsWritesLotsReads(m IMap, b *testing.B) {
	keys := []int64{}
	runtime.GOMAXPROCS(runtime.NumCPU())
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			rand.Seed(time.Now().UTC().UnixNano())
			for i := 0; i < NumWritesInRWTest; i++ {
				/* Do a read */
				if i > 0 && i%WriteRatioLow == 0 {
					idx := rand.Int63n(int64(len(keys)))
					v, ok := m.Get(keys[idx])
					if ok {
						if v != keys[idx] {
							b.Error("Wrong value for key", keys[idx], ". Expect ", keys[idx], ". Got ", v)
						}
					} else {
						b.Error("Failed to get key ", keys[idx])
					}
				} else {
					/* Do write */
					k := rand.Int63()
					v := k
					keys = append(keys, k)
					m.Put(k, v)
				}
			}
		}
	})
}

/*
*  4. Initialize a large table (fitting into memory)
*     (do not test the initialization part),
*	   then lots of uniformly random reads ->
*     cache behavior when reading from an unchanging table
 */
func benchmarkLotsReads(nKeys int64, m IMap, b *testing.B) {
	b.ResetTimer()
	runtime.GOMAXPROCS(runtime.NumCPU())
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			rand.Seed(time.Now().UTC().UnixNano())
			for i := 0; i < NumReadsInReadOnlyTest; i++ {
				k := rand.Int63n(nKeys)
				v, ok := m.Get(k)
				if ok {
					if v != k {
						b.Error("Wrong value for key", k, ". Expect ", k, ". Got ", v)
					}
				} else {
					b.Error("Failed to get key ", k)
				}
			}
		}
	})
}

func Writer(do, done chan bool, m IMap, nKeys, numWrites int64, b *testing.B) {
	<-do
	for i := int64(0); i < numWrites; i++ {
		k := rand.Int63n(nKeys)
		v := k
		m.Put(k, v)
	}
	done <- true
}

func Reader(do, done chan bool, m IMap, nKeys, numReads int64, b *testing.B) {
	<-do
	for i := int64(0); i < numReads; i++ {
		k := rand.Int63n(nKeys)
		v, ok := m.Get(k)
		if ok {
			if v != k {
				b.Error("Wrong value for key", k, ". Expect ", k, ". Got ", v)
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
					m.Put(i, i)
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

/* 1. =======================Lots of concurrent writes======================= */
func BenchmarkLockMapLotsWrite(b *testing.B) {
	benchmarkConcurrentWrites(lockmap.NewLockMap(), b)
}

func BenchmarkRWLockMapLotsWrite(b *testing.B) {
	benchmarkConcurrentWrites(rwlockmap.NewRWLockMap(), b)
}

func BenchmarkParallelMapLotsWrite(b *testing.B) {
	benchmarkConcurrentWrites(pmap.NewParallelMap(), b)
}

func BenchmarkGotomicMapLotsWrite(b *testing.B) {
	benchmarkConcurrentWrites(gotomic.NewGotomicMap(), b)
}
/* 2. ==================Lots of concurrent writes, few reads================= */
func BenchmarkLockMapLotsWritesFewReads(b *testing.B) {
	benchmarkLotsWritesFewReads(lockmap.NewLockMap(), b)
}

func BenchmarkRWLockMapLotsWritesFewReads(b *testing.B) {
	benchmarkLotsWritesFewReads(rwlockmap.NewRWLockMap(), b)
}

func BenchmarkParallelMapLotsWritesFewReads(b *testing.B) {
	benchmarkLotsWritesFewReads(pmap.NewParallelMap(), b)
}

func BenchmarkGotomicMapLotsWritesFewReads(b *testing.B) {
	benchmarkLotsWritesFewReads(gotomic.NewGotomicMap(), b)
}

/* 3. ================Lots of concurrent writes, lots of reads=============== */
func BenchmarkLockMapLotsWritesLotsReads(b *testing.B) {
	benchmarkLotsWritesLotsReads(lockmap.NewLockMap(), b)
}

func BenchmarkRWLockMapLotsWritesLotsReads(b *testing.B) {
	benchmarkLotsWritesLotsReads(lockmap.NewLockMap(), b)
}

func BenchmarkParallelMapLotsWritesLotsReads(b *testing.B) {
	benchmarkLotsWritesLotsReads(pmap.NewParallelMap(), b)
}

func BenchmarkGotomicMapLotsWritesLotsReads(b *testing.B) {
	benchmarkLotsWritesLotsReads(gotomic.NewGotomicMap(), b)
}
/* 4. =======================Lots of concurrent reads======================== */
func BenchmarkNativeMapLotsReads(b *testing.B) {
	m := nativemap.NewNativeMap()
	/* Initialize a big map */
	InitializeMap(NumKeysInBigMap, m)
	benchmarkLotsReads(NumKeysInBigMap, m, b)
}

func BenchmarkLockMapLotsReads(b *testing.B) {
	m := lockmap.NewLockMap()
	/* Initialize a big map */
	InitializeMap(NumKeysInBigMap, m)
	benchmarkLotsReads(NumKeysInBigMap, m, b)
}

func BenchmarkRWLockMapLotsReads(b *testing.B) {
	m := rwlockmap.NewRWLockMap()
	/* Initialize a big map */
	InitializeMap(NumKeysInBigMap, m)
	benchmarkLotsReads(NumKeysInBigMap, m, b)
}

func BenchmarkParallelMapLotsReads(b *testing.B) {
	m := pmap.NewParallelMap()
	/* Initialize a big map */
	InitializeMap(NumKeysInBigMap, m)
	benchmarkLotsReads(NumKeysInBigMap, m, b)
}

func BenchmarkGotomicMapLotsReads(b *testing.B) {
	m := gotomic.NewGotomicMap()
	/* Initialize a big map */
	InitializeMap(NumKeysInBigMap, m)
	benchmarkLotsReads(NumKeysInBigMap, m, b)
}
/* 5. ================1, 2, 3, 4 but very large tables (>16 GB)===============*/
/* 6. 1, 2, 3, 4 but with a particular set of keys read/wrote more frequently */
/* 7. ========1, 2, 3, 4 but with reading and writing sequential keys=========*/

/* TODO */

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

/* 9. ====================10 concurrent writers, 100 readers==================*/
func BenchmarkLockMapConcurrentWriterReaders2(b *testing.B) {
	benchmarkConcurrentWriterReaders(10, 100, lockmap.NewLockMap(), b)
}

func BenchmarkRWLockMapConcurrentWriterReaders2(b *testing.B) {
	benchmarkConcurrentWriterReaders(10, 100, rwlockmap.NewRWLockMap(), b)
}

/* 10. ====================1 concurrent writers, 100 readers==================*/
func BenchmarkLockMapConcurrentWriterReaders3(b *testing.B) {
	benchmarkConcurrentWriterReaders(1, 100, lockmap.NewLockMap(), b)
}

func BenchmarkRWLockMapConcurrentWriterReaders3(b *testing.B) {
	benchmarkConcurrentWriterReaders(1, 100, rwlockmap.NewRWLockMap(), b)
}

/* 11. =======================Write, delete, write=========================== */
func BenchmarkLockMapWriteDeleteWrite(b *testing.B) {
	benchmarkConcurrentWriteDeleteWrite(lockmap.NewLockMap(), b)
}

func BenchmarkRWLockMapWriteDeleteWrite(b *testing.B) {
	benchmarkConcurrentWriteDeleteWrite(rwlockmap.NewRWLockMap(), b)
}
