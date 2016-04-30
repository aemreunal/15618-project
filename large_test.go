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
	"concurrent"
)

const (
	NumKeysInLargeMap                 = 1024 * 1024 * 1024 *2 // 2 G
	NumWritesInWriteOnlyTestLarge     = 1024 * 1024 * 1024 * 2 // 2 G
	NumWritesInRWTestLarge            = 1024 * 1024 * 1024 * 2 // 2 G
	NumReadsInReadOnlyTestLarge       = 1024 * 1024 * 1024 * 2 // 2 G
)

/* 4. =======================Lots of concurrent reads======================== */
func BenchmarkNativeMapLotsReads(b *testing.B) {
	benchmarkLotsReads(nativemap.NewNativeMap(), b)
}

func BenchmarkLockMapLotsReads(b *testing.B) {
	benchmarkLotsReads(lockmap.NewLockMap(), b)
}

func BenchmarkRWLockMapLotsReads(b *testing.B) {
	benchmarkLotsReads(rwlockmap.NewRWLockMap(), b)
}

func BenchmarkParallelMapLotsReads(b *testing.B) {
	benchmarkLotsReads(pmap.NewParallelMap(), b)
}

func BenchmarkGotomicMapLotsReads(b *testing.B) {
	benchmarkLotsReads(gotomic.NewGotomicMap(), b)
}

func BenchmarkConcurrentMapLotsReads(b *testing.B) {
	benchmarkLotsReads(concurrent.NewConcurrentMap(), b)
}

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

/************************** 5_2 Lots of concurrent writes, lots reads ***********************************/
func BenchmarkLockMapLotsWritesLotsReadsLarge(b *testing.B) {
	benchmarkLotsWritesLotsReads(lockmap.NewLockMap(), b, NumWritesInRWTestLarge)
}

func BenchmarkRWLockMapLotsWritesLotsReadsLarge(b *testing.B) {
	benchmarkLotsWritesLotsReads(lockmap.NewLockMap(), b, NumWritesInRWTestSmall)
}

func BenchmarkParallelMapLotsWritesLotsReadsLarge(b *testing.B) {
	benchmarkLotsWritesLotsReads(pmap.NewParallelMap(), b, NumWritesInRWTestSmall)
}

func BenchmarkGotomicMapLotsWritesLotsReadsLarge(b *testing.B) {
	benchmarkLotsWritesLotsReads(gotomic.NewGotomicMap(), b, NumWritesInRWTestSmall)
}

func BenchmarkConcurrentMapLotsWritesLotsReadsLarge(b *testing.B) {
	benchmarkLotsWritesLotsReads(concurrent.NewConcurrentMap(), b, NumWritesInRWTestSmall)
}
/* 6. 1, 2, 3, 4 but with a particular set of keys read/wrote more frequently */
/* 7. ========1, 2, 3, 4 but with reading and writing sequential keys=========*/

/* TODO */

