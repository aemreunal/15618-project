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

const (
	NumKeysInLargeMap                 = 1024 * 1024 * 1024 *2 // 2 G
	NumWritesInWriteOnlyTestLarge     = 1024 * 1024 * 1024 * 2 // 2 G
	NumWritesInRWTestLarge            = 1024 * 1024 * 1024 * 2 // 2 G
	NumReadsInReadOnlyTestLarge       = 1024 * 1024 * 1024 * 2 // 2 G
)

/* 5. ================1, 2, 3, 4 but very large tables (>16 GB)===============*/

/************************** 5_1 Lots writes ***********************************/
func BenchmarkLockMapLotsWriteLarge(b *testing.B) {
	benchmarkConcurrentWrites(lockmap.NewLockMap(), b, NumWritesInWriteOnlyTestLarge)
}

func BenchmarkRWLockMapLotsWriteLarge(b *testing.B) {
	benchmarkConcurrentWrites(rwlockmap.NewRWLockMap(), b, NumWritesInWriteOnlyTestLarge)
}

func BenchmarkParallelMapLotsWriteLarge(b *testing.B) {
	benchmarkConcurrentWrites(pmap.NewParallelMap(), b, NumWritesInWriteOnlyTestLarge)
}

func BenchmarkGotomicMapLotsWriteLarge(b *testing.B) {
	benchmarkConcurrentWrites(gotomic.NewGotomicMap(), b, NumWritesInWriteOnlyTestLarge)
}

func BenchmarkConcurrentMapLotsWriteLarge(b *testing.B) {
	benchmarkConcurrentWrites(concurrent.NewConcurrentMap(), b, NumWritesInWriteOnlyTestLarge)
}
/************************** 5_2 Lots of concurrent writes, few reads ***********************************/
func BenchmarkLockMapLotsWritesFewReadsLarge(b *testing.B) {
	benchmarkLotsWritesFewReads(lockmap.NewLockMap(), b, NumWritesInRWTestLarge)
}

func BenchmarkRWLockMapLotsWritesFewReadsLarge(b *testing.B) {
	benchmarkLotsWritesFewReads(rwlockmap.NewRWLockMap(), b, NumWritesInRWTestLarge)
}

func BenchmarkParallelMapLotsWritesFewReadsLarge(b *testing.B) {
	benchmarkLotsWritesFewReads(pmap.NewParallelMap(), b, NumWritesInRWTestLarge)
}

func BenchmarkGotomicMapLotsWritesFewReadsLarge(b *testing.B) {
	benchmarkLotsWritesFewReads(gotomic.NewGotomicMap(), b, NumWritesInRWTestLarge)
}

func BenchmarkConcurrentMapLotsWritesFewReadsLarge(b *testing.B) {
	benchmarkLotsWritesFewReads(concurrent.NewConcurrentMap(), b, NumWritesInRWTestLarge)
}

/************************** 5_3 Lots of concurrent writes, lots reads ***********************************/
func BenchmarkLockMapLotsWritesLotsReadsLarge(b *testing.B) {
	benchmarkLotsWritesLotsReads(lockmap.NewLockMap(), b, NumWritesInRWTestLarge)
}

func BenchmarkRWLockMapLotsWritesLotsReadsLarge(b *testing.B) {
	benchmarkLotsWritesLotsReads(lockmap.NewLockMap(), b, NumWritesInRWTestLarge)
}

func BenchmarkParallelMapLotsWritesLotsReadsLarge(b *testing.B) {
	benchmarkLotsWritesLotsReads(pmap.NewParallelMap(), b, NumWritesInRWTestLarge)
}

func BenchmarkGotomicMapLotsWritesLotsReadsLarge(b *testing.B) {
	benchmarkLotsWritesLotsReads(gotomic.NewGotomicMap(), b, NumWritesInRWTestLarge)
}

func BenchmarkConcurrentMapLotsWritesLotsReadsLarge(b *testing.B) {
	benchmarkLotsWritesLotsReads(concurrent.NewConcurrentMap(), b, NumWritesInRWTestLarge)
}

/************************** 5_4 Lots of concurrent reads ***********************************/
func BenchmarkNativeMapLotsReadsLarge(b *testing.B) {
	benchmarkLotsReads(nativemap.NewNativeMap(), b, NumKeysInLargeMap, NumReadsInReadOnlyTestLarge)
}

func BenchmarkLockMapLotsReadsLarge(b *testing.B) {
	benchmarkLotsReads(lockmap.NewLockMap(), b, NumKeysInLargeMap, NumReadsInReadOnlyTestLarge)
}

func BenchmarkRWLockMapLotsReadsLarge(b *testing.B) {
	benchmarkLotsReads(rwlockmap.NewRWLockMap(), b, NumKeysInLargeMap, NumReadsInReadOnlyTestLarge)
}

func BenchmarkParallelMapLotsReadsLarge(b *testing.B) {
	benchmarkLotsReads(pmap.NewParallelMap(), b, NumKeysInLargeMap, NumReadsInReadOnlyTestLarge)
}

func BenchmarkGotomicMapLotsReadsLarge(b *testing.B) {
	benchmarkLotsReads(gotomic.NewGotomicMap(), b, NumKeysInLargeMap, NumReadsInReadOnlyTestLarge)
}

func BenchmarkConcurrentMapLotsReadsLarge(b *testing.B) {
	benchmarkLotsReads(concurrent.NewConcurrentMap(), b, NumKeysInLargeMap, NumReadsInReadOnlyTestLarge)
}

/* 6. 1, 2, 3, 4 but with a particular set of keys read/wrote more frequently */

/* 7. ========1, 2, 3, 4 but with reading and writing sequential keys=========*/


