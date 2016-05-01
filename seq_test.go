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

/* 7. ========1, 2, 3, 4 but with reading and writing sequential keys=========*/

/************************** 7_1 Lots writes ***********************************/
func BenchmarkLockMapLotsWriteSeqKeys(b *testing.B) {
	benchmarkConcurrentWritesSequential(lockmap.NewLockMap(), b, NumWritesInWriteOnlyTestSmall)
}

func BenchmarkRWLockMapLotsWriteSeqKeys(b *testing.B) {
	benchmarkConcurrentWritesSequential(rwlockmap.NewRWLockMap(), b, NumWritesInWriteOnlyTestSmall)
}

func BenchmarkParallelMapLotsWriteSeqKeys(b *testing.B) {
	benchmarkConcurrentWritesSequential(pmap.NewParallelMap(), b, NumWritesInWriteOnlyTestSmall)
}

func BenchmarkGotomicMapLotsWriteSeqKeys(b *testing.B) {
	benchmarkConcurrentWritesSequential(gotomic.NewGotomicMap(), b, NumWritesInWriteOnlyTestSmall)
}

func BenchmarkConcurrentMapLotsWriteSeqKeys(b *testing.B) {
	benchmarkConcurrentWritesSequential(concurrent.NewConcurrentMap(), b, NumWritesInWriteOnlyTestSmall)
}

/************************** 7_2 Lots of concurrent writes, few reads ***********************************/
func BenchmarkLockMapLotsWritesFewReadsSeqKeys(b *testing.B) {
	benchmarkLotsWritesFewReadsSequential(lockmap.NewLockMap(), b, NumWritesInRWTestSmall)
}

func BenchmarkRWLockMapLotsWritesFewReadsSeqKeys(b *testing.B) {
	benchmarkLotsWritesFewReadsSequential(rwlockmap.NewRWLockMap(), b, NumWritesInRWTestSmall)
}

func BenchmarkParallelMapLotsWritesFewReadsSeqKeys(b *testing.B) {
	benchmarkLotsWritesFewReadsSequential(pmap.NewParallelMap(), b, NumWritesInRWTestSmall)
}

func BenchmarkGotomicMapLotsWritesFewReadsSeqKeys(b *testing.B) {
	benchmarkLotsWritesFewReadsSequential(gotomic.NewGotomicMap(), b, NumWritesInRWTestSmall)
}

func BenchmarkConcurrentMapLotsWritesFewReadsSeqKeys(b *testing.B) {
	benchmarkLotsWritesFewReadsSequential(concurrent.NewConcurrentMap(), b, NumWritesInRWTestSmall)
}

/************************** 7_3 Lots of concurrent writes, lots reads ***********************************/
func BenchmarkLockMapLotsWritesLotsReadsSeqKeys(b *testing.B) {
	benchmarkLotsWritesLotsReadsSequential(lockmap.NewLockMap(), b, NumWritesInRWTestSmall)
}

func BenchmarkRWLockMapLotsWritesLotsReadsSeqKeys(b *testing.B) {
	benchmarkLotsWritesLotsReadsSequential(lockmap.NewLockMap(), b, NumWritesInRWTestSmall)
}

func BenchmarkParallelMapLotsWritesLotsReadsSeqKeys(b *testing.B) {
	benchmarkLotsWritesLotsReadsSequential(pmap.NewParallelMap(), b, NumWritesInRWTestSmall)
}

func BenchmarkGotomicMapLotsWritesLotsReadsSeqKeys(b *testing.B) {
	benchmarkLotsWritesLotsReadsSequential(gotomic.NewGotomicMap(), b, NumWritesInRWTestSmall)
}

func BenchmarkConcurrentMapLotsWritesLotsReadsSeqKeys(b *testing.B) {
	benchmarkLotsWritesLotsReadsSequential(concurrent.NewConcurrentMap(), b, NumWritesInRWTestSmall)
}

/************************** 7_4 Lots of concurrent reads ***********************************/
func BenchmarkNativeMapLotsReadsSeqKeys(b *testing.B) {
	benchmarkLotsReadsSequential(nativemap.NewNativeMap(), b, NumKeysInBigMap, NumReadsInReadOnlyTestSmall)
}

func BenchmarkLockMapLotsReadsSeqKeys(b *testing.B) {
	benchmarkLotsReadsSequential(lockmap.NewLockMap(), b, NumKeysInBigMap, NumReadsInReadOnlyTestSmall)
}

func BenchmarkRWLockMapLotsReadsSeqKeys(b *testing.B) {
	benchmarkLotsReadsSequential(rwlockmap.NewRWLockMap(), b, NumKeysInBigMap, NumReadsInReadOnlyTestSmall)
}

func BenchmarkParallelMapLotsReadsSeqKeys(b *testing.B) {
	benchmarkLotsReadsSequential(pmap.NewParallelMap(), b, NumKeysInBigMap, NumReadsInReadOnlyTestSmall)
}

func BenchmarkGotomicMapLotsReadsSeqKeys(b *testing.B) {
	benchmarkLotsReadsSequential(gotomic.NewGotomicMap(), b, NumKeysInBigMap, NumReadsInReadOnlyTestSmall)
}

func BenchmarkConcurrentMapLotsReadsSeqKeys(b *testing.B) {
	benchmarkLotsReadsSequential(concurrent.NewConcurrentMap(), b, NumKeysInBigMap, NumReadsInReadOnlyTestSmall)
}
