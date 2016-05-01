package benchmark

import (
	"concurrent"
	"gotomic"
	"lockmap"
	"nativemap"
	"pmap"
	"rwlockmap"
	"testing"
)

/* 6. 1, 2, 3, 4 but with a particular set of keys read/wrote more frequently */

/************************** 6_1 Lots writes ***********************************/
func BenchmarkLockMapLotsWriteFreqKeys(b *testing.B) {
	benchmarkConcurrentWritesNormalDist(lockmap.NewLockMap(), b, NumWritesInWriteOnlyTestSmall)
}

func BenchmarkRWLockMapLotsWriteFreqKeys(b *testing.B) {
	benchmarkConcurrentWritesNormalDist(rwlockmap.NewRWLockMap(), b, NumWritesInWriteOnlyTestSmall)
}

func BenchmarkParallelMapLotsWriteFreqKeys(b *testing.B) {
	benchmarkConcurrentWritesNormalDist(pmap.NewParallelMap(), b, NumWritesInWriteOnlyTestSmall)
}

func BenchmarkGotomicMapLotsWriteFreqKeys(b *testing.B) {
	benchmarkConcurrentWritesNormalDist(gotomic.NewGotomicMap(), b, NumWritesInWriteOnlyTestSmall)
}

func BenchmarkConcurrentMapLotsWriteFreqKeys(b *testing.B) {
	benchmarkConcurrentWritesNormalDist(concurrent.NewConcurrentMap(), b, NumWritesInWriteOnlyTestSmall)
}

/************************** 6_2 Lots of concurrent writes, few reads ***********************************/
func BenchmarkLockMapLotsWritesFewReadsFreqKeys(b *testing.B) {
	benchmarkLotsWritesFewReadsNormalDist(lockmap.NewLockMap(), b, NumWritesInRWTestSmall)
}

func BenchmarkRWLockMapLotsWritesFewReadsFreqKeys(b *testing.B) {
	benchmarkLotsWritesFewReadsNormalDist(rwlockmap.NewRWLockMap(), b, NumWritesInRWTestSmall)
}

func BenchmarkParallelMapLotsWritesFewReadsFreqKeys(b *testing.B) {
	benchmarkLotsWritesFewReadsNormalDist(pmap.NewParallelMap(), b, NumWritesInRWTestSmall)
}

func BenchmarkGotomicMapLotsWritesFewReadsFreqKeys(b *testing.B) {
	benchmarkLotsWritesFewReadsNormalDist(gotomic.NewGotomicMap(), b, NumWritesInRWTestSmall)
}

func BenchmarkConcurrentMapLotsWritesFewReadsFreqKeys(b *testing.B) {
	benchmarkLotsWritesFewReadsNormalDist(concurrent.NewConcurrentMap(), b, NumWritesInRWTestSmall)
}

/************************** 6_3 Lots of concurrent writes, lots reads ***********************************/
func BenchmarkLockMapLotsWritesLotsReadsFreqKeys(b *testing.B) {
	benchmarkLotsWritesLotsReadsNormalDist(lockmap.NewLockMap(), b, NumWritesInRWTestSmall)
}

func BenchmarkRWLockMapLotsWritesLotsReadsFreqKeys(b *testing.B) {
	benchmarkLotsWritesLotsReadsNormalDist(lockmap.NewLockMap(), b, NumWritesInRWTestSmall)
}

func BenchmarkParallelMapLotsWritesLotsReadsFreqKeys(b *testing.B) {
	benchmarkLotsWritesLotsReadsNormalDist(pmap.NewParallelMap(), b, NumWritesInRWTestSmall)
}

func BenchmarkGotomicMapLotsWritesLotsReadsFreqKeys(b *testing.B) {
	benchmarkLotsWritesLotsReadsNormalDist(gotomic.NewGotomicMap(), b, NumWritesInRWTestSmall)
}

func BenchmarkConcurrentMapLotsWritesLotsReadsFreqKeys(b *testing.B) {
	benchmarkLotsWritesLotsReadsNormalDist(concurrent.NewConcurrentMap(), b, NumWritesInRWTestSmall)
}

/************************** 5_4 Lots of concurrent reads ***********************************/
func BenchmarkNativeMapLotsReadsFreqKeys(b *testing.B) {
	benchmarkLotsReadsNormalDist(nativemap.NewNativeMap(), b, NumKeysInBigMap, NumReadsInReadOnlyTestSmall)
}

func BenchmarkLockMapLotsReadsFreqKeys(b *testing.B) {
	benchmarkLotsReadsNormalDist(lockmap.NewLockMap(), b, NumKeysInBigMap, NumReadsInReadOnlyTestSmall)
}

func BenchmarkRWLockMapLotsReadsFreqKeys(b *testing.B) {
	benchmarkLotsReadsNormalDist(rwlockmap.NewRWLockMap(), b, NumKeysInBigMap, NumReadsInReadOnlyTestSmall)
}

func BenchmarkParallelMapLotsReadsFreqKeys(b *testing.B) {
	benchmarkLotsReadsNormalDist(pmap.NewParallelMap(), b, NumKeysInBigMap, NumReadsInReadOnlyTestSmall)
}

func BenchmarkGotomicMapLotsReadsFreqKeys(b *testing.B) {
	benchmarkLotsReadsNormalDist(gotomic.NewGotomicMap(), b, NumKeysInBigMap, NumReadsInReadOnlyTestSmall)
}

func BenchmarkConcurrentMapLotsReadsFreqKeys(b *testing.B) {
	benchmarkLotsReadsNormalDist(concurrent.NewConcurrentMap(), b, NumKeysInBigMap, NumReadsInReadOnlyTestSmall)
}
