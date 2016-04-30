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

/* 7. ========1, 2, 3, 4 but with reading and writing sequential keys=========*/

/************************** 7_1 Lots writes ***********************************/
func BenchmarkLockMapLotsWriteSeqKeys(b *testing.B) {
	benchmarkConcurrentWrites(lockmap.NewLockMap(), b, NumWritesInWriteOnlyTestSmall)
}

func BenchmarkRWLockMapLotsWriteSeqKeys(b *testing.B) {
	benchmarkConcurrentWrites(rwlockmap.NewRWLockMap(), b, NumWritesInWriteOnlyTestSmall)
}

func BenchmarkParallelMapLotsWriteSeqKeys(b *testing.B) {
	benchmarkConcurrentWrites(pmap.NewParallelMap(), b, NumWritesInWriteOnlyTestSmall)
}

func BenchmarkGotomicMapLotsWriteSeqKeys(b *testing.B) {
	benchmarkConcurrentWrites(gotomic.NewGotomicMap(), b, NumWritesInWriteOnlyTestSmall)
}

func BenchmarkConcurrentMapLotsWriteSeqKeys(b *testing.B) {
	benchmarkConcurrentWrites(concurrent.NewConcurrentMap(), b, NumWritesInWriteOnlyTestSmall)
}
/************************** 7_2 Lots of concurrent writes, few reads ***********************************/
func BenchmarkLockMapLotsWritesFewReadsSeqKeys(b *testing.B) {
	benchmarkLotsWritesFewReads(lockmap.NewLockMap(), b, NumWritesInRWTestSmall)
}

func BenchmarkRWLockMapLotsWritesFewReadsSeqKeys(b *testing.B) {
	benchmarkLotsWritesFewReads(rwlockmap.NewRWLockMap(), b, NumWritesInRWTestSmall)
}

func BenchmarkParallelMapLotsWritesFewReadsSeqKeys(b *testing.B) {
	benchmarkLotsWritesFewReads(pmap.NewParallelMap(), b, NumWritesInRWTestSmall)
}

func BenchmarkGotomicMapLotsWritesFewReadsSeqKeys(b *testing.B) {
	benchmarkLotsWritesFewReads(gotomic.NewGotomicMap(), b, NumWritesInRWTestSmall)
}

func BenchmarkConcurrentMapLotsWritesFewReadsSeqKeys(b *testing.B) {
	benchmarkLotsWritesFewReads(concurrent.NewConcurrentMap(), b, NumWritesInRWTestSmall)
}

/************************** 7_3 Lots of concurrent writes, lots reads ***********************************/
func BenchmarkLockMapLotsWritesLotsReadsSeqKeys(b *testing.B) {
	benchmarkLotsWritesLotsReads(lockmap.NewLockMap(), b, NumWritesInRWTestSmall)
}

func BenchmarkRWLockMapLotsWritesLotsReadsSeqKeys(b *testing.B) {
	benchmarkLotsWritesLotsReads(lockmap.NewLockMap(), b, NumWritesInRWTestSmall)
}

func BenchmarkParallelMapLotsWritesLotsReadsSeqKeys(b *testing.B) {
	benchmarkLotsWritesLotsReads(pmap.NewParallelMap(), b, NumWritesInRWTestSmall)
}

func BenchmarkGotomicMapLotsWritesLotsReadsSeqKeys(b *testing.B) {
	benchmarkLotsWritesLotsReads(gotomic.NewGotomicMap(), b, NumWritesInRWTestSmall)
}

func BenchmarkConcurrentMapLotsWritesLotsReadsSeqKeys(b *testing.B) {
	benchmarkLotsWritesLotsReads(concurrent.NewConcurrentMap(), b, NumWritesInRWTestSmall)
}

/************************** 7_4 Lots of concurrent reads ***********************************/
func BenchmarkNativeMapLotsReadsSeqKeys(b *testing.B) {
	benchmarkLotsReads(nativemap.NewNativeMap(), b, NumKeysInBigMap, NumReadsInReadOnlyTestSmall)
}

func BenchmarkLockMapLotsReadsSeqKeys(b *testing.B) {
	benchmarkLotsReads(lockmap.NewLockMap(), b, NumKeysInBigMap, NumReadsInReadOnlyTestSmall)
}

func BenchmarkRWLockMapLotsReadsSeqKeys(b *testing.B) {
	benchmarkLotsReads(rwlockmap.NewRWLockMap(), b, NumKeysInBigMap, NumReadsInReadOnlyTestSmall)
}

func BenchmarkParallelMapLotsReadsSeqKeys(b *testing.B) {
	benchmarkLotsReads(pmap.NewParallelMap(), b, NumKeysInBigMap, NumReadsInReadOnlyTestSmall)
}

func BenchmarkGotomicMapLotsReadsSeqKeys(b *testing.B) {
	benchmarkLotsReads(gotomic.NewGotomicMap(), b, NumKeysInBigMap, NumReadsInReadOnlyTestSmall)
}

func BenchmarkConcurrentMapLotsReadsSeqKeys(b *testing.B) {
	benchmarkLotsReads(concurrent.NewConcurrentMap(), b, NumKeysInBigMap, NumReadsInReadOnlyTestSmall)
}
