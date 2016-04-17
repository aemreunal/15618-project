package benchmark

import (
	"runtime"
	"testing"
	"math/rand"
	"time"
)

const (
	StrKeyLen = 8
	NumWritesInWriteOnlyTest = 1024*1024*16  // 16 M
	NumWritesInRWTest = 1024*1024*16 				 // 16 M
	NumReadsInReadOnlyTest = 1024*1024*16    // 16 M
	NumIterationInConcurrentReadWrite = 1024*16
	NumWriteDeleteIter = 100
	NumKeysInBigMap = 1024*1024*16
	NumKeysInSmallMap = 1024*16
	WriteRatioHigh = 1000
	WriteRatioLow = 2
	chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVXYZabcdefghijklmnopqrstuvwxyz"
)

type IMap interface {
	Get(k interface{}) (interface{}, bool)
	Put(k,v interface{}) interface{}
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
		m.Put(i, i);
	}
}

func benchmarkPutGetBasic(m IMap, b *testing.B) {
	for i := 0; i < b.N; i++ {
		m.Put(i, i)
		j, ok := m.Get(i)
		if !ok {
			b.Error("Failed to get key ", i);
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
	runtime.GOMAXPROCS(runtime.NumCPU())
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
	runtime.GOMAXPROCS(runtime.NumCPU())
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			rand.Seed(time.Now().UTC().UnixNano())
			for i := 0; i < NumWritesInRWTest; i++ {
				/* Do a read */
				if i > 0 && i % WriteRatioHigh == 0 {
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
				if i > 0 && i % WriteRatioLow == 0 {
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
	for i := 0; i < numWriters + numReaders; i++ {
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

func BenchmarkNativeMapPutGetBasic(b *testing.B) {
	benchmarkPutGetBasic(NewNativeMap(), b)
}

func BenchmarkLockMapPutGetBasic(b *testing.B) {
	benchmarkPutGetBasic(NewLockMap(), b)
}

func BenchmarkRWLockMapPutGetBasic(b *testing.B) {
	benchmarkPutGetBasic(NewRWLockMap(), b)
}

/* 1. =======================Lots of concurrent writes======================= */
func BenchmarkLockMapLotsWriteInMem(b *testing.B) {
	benchmarkConcurrentWrites(NewLockMap(), b)
}

func BenchmarkRWLockMapLotsWriteInMem(b *testing.B) {
	benchmarkConcurrentWrites(NewRWLockMap(), b)
}

/* 2. ==================Lots of concurrent writes, few reads================= */
func BenchmarkLockMapLotsWritesFewReads(b *testing.B) {
	benchmarkLotsWritesFewReads(NewLockMap(), b)
}

func BenchmarkRWLockMapLotsWritesFewReads(b *testing.B) {
	benchmarkLotsWritesFewReads(NewRWLockMap(), b)
}

/* 3. ================Lots of concurrent writes, lots of reads=============== */
func BenchmarkLockMapLotsWritesLotsReads(b *testing.B) {
	benchmarkLotsWritesLotsReads(NewLockMap(), b)
}

func BenchmarkRWLockMapLotsWritesLotsReads(b *testing.B) {
	benchmarkLotsWritesLotsReads(NewLockMap(), b)
}

/* 4. =======================Lots of concurrent reads======================== */
func BenchmarkNativeMapLotsReads(b *testing.B) {
	m := NewNativeMap()
	/* Initialize a big map */
	InitializeMap(NumKeysInBigMap, m)
	benchmarkLotsReads(NumKeysInBigMap, m, b)
}

func BenchmarkLockMapLotsReads(b *testing.B) {
	m := NewLockMap()
	/* Initialize a big map */
	InitializeMap(NumKeysInBigMap, m)
	benchmarkLotsReads(NumKeysInBigMap, m, b)
}

func BenchmarkRWLockMapLotsReads(b *testing.B) {
	m := NewRWLockMap()
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
	benchmarkConcurrentWriterReaders(100, 10, NewLockMap(), b)
}

func BenchmarkRWLockMapConcurrentWriterReaders1(b *testing.B) {
	benchmarkConcurrentWriterReaders(100, 10, NewRWLockMap(), b)
}

/* 9. ====================10 concurrent writers, 100 readers==================*/
func BenchmarkLockMapConcurrentWriterReaders2(b *testing.B) {
	benchmarkConcurrentWriterReaders(10, 100, NewLockMap(), b)
}

func BenchmarkRWLockMapConcurrentWriterReaders2(b *testing.B) {
	benchmarkConcurrentWriterReaders(10, 100, NewRWLockMap(), b)
}

/* 10. ====================1 concurrent writers, 100 readers==================*/
func BenchmarkLockMapConcurrentWriterReaders3(b *testing.B) {
	benchmarkConcurrentWriterReaders(1, 100, NewLockMap(), b)
}

func BenchmarkRWLockMapConcurrentWriterReaders3(b *testing.B) {
	benchmarkConcurrentWriterReaders(1, 100, NewRWLockMap(), b)
}

/* 11. =======================Write, delete, write=========================== */
func BenchmarkLockMapWriteDeleteWrite(b *testing.B) {
	benchmarkConcurrentWriteDeleteWrite(NewLockMap(), b)
}

func BenchmarkRWLockMapWriteDeleteWrite(b *testing.B) {
	benchmarkConcurrentWriteDeleteWrite(NewRWLockMap(), b)
}
