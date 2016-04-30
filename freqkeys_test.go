package benchmark

import (
	"testing"
	"nativemap"
	"lockmap"
	"rwlockmap"
	"pmap"
	"gotomic"
	"concurrent"
)

/* 6. 1, 2, 3, 4 but with a particular set of keys read/wrote more frequently */

/************************** 6_1 Lots writes ***********************************/
func BenchmarkLockMapLotsWriteFreqKeys(b *testing.B) {
	benchmarkConcurrentWrites(lockmap.NewLockMap(), b, NumWritesInWriteOnlyTestSmall)
}

func BenchmarkRWLockMapLotsWriteFreqKeys(b *testing.B) {
	benchmarkConcurrentWrites(rwlockmap.NewRWLockMap(), b, NumWritesInWriteOnlyTestSmall)
}

func BenchmarkParallelMapLotsWriteFreqKeys(b *testing.B) {
	benchmarkConcurrentWrites(pmap.NewParallelMap(), b, NumWritesInWriteOnlyTestSmall)
}

func BenchmarkGotomicMapLotsWriteFreqKeys(b *testing.B) {
	benchmarkConcurrentWrites(gotomic.NewGotomicMap(), b, NumWritesInWriteOnlyTestSmall)
}

func BenchmarkConcurrentMapLotsWriteFreqKeys(b *testing.B) {
	benchmarkConcurrentWrites(concurrent.NewConcurrentMap(), b, NumWritesInWriteOnlyTestSmall)
}
/************************** 6_2 Lots of concurrent writes, few reads ***********************************/
func BenchmarkLockMapLotsWritesFewReadsFreqKeys(b *testing.B) {
	benchmarkLotsWritesFewReads(lockmap.NewLockMap(), b, NumWritesInRWTestSmall)
}

func BenchmarkRWLockMapLotsWritesFewReadsFreqKeys(b *testing.B) {
	benchmarkLotsWritesFewReads(rwlockmap.NewRWLockMap(), b, NumWritesInRWTestSmall)
}

func BenchmarkParallelMapLotsWritesFewReadsFreqKeys(b *testing.B) {
	benchmarkLotsWritesFewReads(pmap.NewParallelMap(), b, NumWritesInRWTestSmall)
}

func BenchmarkGotomicMapLotsWritesFewReadsFreqKeys(b *testing.B) {
	benchmarkLotsWritesFewReads(gotomic.NewGotomicMap(), b, NumWritesInRWTestSmall)
}

func BenchmarkConcurrentMapLotsWritesFewReadsFreqKeys(b *testing.B) {
	benchmarkLotsWritesFewReads(concurrent.NewConcurrentMap(), b, NumWritesInRWTestSmall)
}

/************************** 6_3 Lots of concurrent writes, lots reads ***********************************/
func BenchmarkLockMapLotsWritesLotsReadsFreqKeys(b *testing.B) {
	benchmarkLotsWritesLotsReads(lockmap.NewLockMap(), b, NumWritesInRWTestSmall)
}

func BenchmarkRWLockMapLotsWritesLotsReadsFreqKeys(b *testing.B) {
	benchmarkLotsWritesLotsReads(lockmap.NewLockMap(), b, NumWritesInRWTestSmall)
}

func BenchmarkParallelMapLotsWritesLotsReadsFreqKeys(b *testing.B) {
	benchmarkLotsWritesLotsReads(pmap.NewParallelMap(), b, NumWritesInRWTestSmall)
}

func BenchmarkGotomicMapLotsWritesLotsReadsFreqKeys(b *testing.B) {
	benchmarkLotsWritesLotsReads(gotomic.NewGotomicMap(), b, NumWritesInRWTestSmall)
}

func BenchmarkConcurrentMapLotsWritesLotsReadsFreqKeys(b *testing.B) {
	benchmarkLotsWritesLotsReads(concurrent.NewConcurrentMap(), b, NumWritesInRWTestSmall)
}

/************************** 5_4 Lots of concurrent reads ***********************************/
func BenchmarkNativeMapLotsReadsFreqKeys(b *testing.B) {
	benchmarkLotsReads(nativemap.NewNativeMap(), b, NumKeysInBigMap, NumReadsInReadOnlyTestSmall)
}

func BenchmarkLockMapLotsReadsFreqKeys(b *testing.B) {
	benchmarkLotsReads(lockmap.NewLockMap(), b, NumKeysInBigMap, NumReadsInReadOnlyTestSmall)
}

func BenchmarkRWLockMapLotsReadsFreqKeys(b *testing.B) {
	benchmarkLotsReads(rwlockmap.NewRWLockMap(), b, NumKeysInBigMap, NumReadsInReadOnlyTestSmall)
}

func BenchmarkParallelMapLotsReadsFreqKeys(b *testing.B) {
	benchmarkLotsReads(pmap.NewParallelMap(), b, NumKeysInBigMap, NumReadsInReadOnlyTestSmall)
}

func BenchmarkGotomicMapLotsReadsFreqKeys(b *testing.B) {
	benchmarkLotsReads(gotomic.NewGotomicMap(), b, NumKeysInBigMap, NumReadsInReadOnlyTestSmall)
}

func BenchmarkConcurrentMapLotsReadsFreqKeys(b *testing.B) {
	benchmarkLotsReads(concurrent.NewConcurrentMap(), b, NumKeysInBigMap, NumReadsInReadOnlyTestSmall)
}


/* 7. ========1, 2, 3, 4 but with reading and writing sequential keys=========*/


