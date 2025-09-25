package sieve

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

// BenchmarkNthPrimeCached benchmarks calls to NthPrime with caching.
func BenchmarkNthPrimeCached(b *testing.B) {
	sieve := NewSieve()

	// Fill cache with enough primes for all benchmarks
	maxIndex := int64(10000000)
	for i := int64(0); i <= maxIndex; i += 100000 {
		sieve.NthPrime(i) // pre-fill cache
	}

	b.ResetTimer()

	benchmarks := []struct {
		name  string
		index int64
	}{
		{"Small", 19},
		{"Medium", 500},
		{"Large", 1000000},
		{"ExtraLarge", 10000000},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = sieve.NthPrime(bm.index)
			}
		})
	}
}

// BenchmarkNthPrimeUncached benchmarks calls to NthPrime without caching.
func BenchmarkNthPrimeUncached(b *testing.B) {
	benchmarks := []struct {
		name  string
		index int64
	}{
		{"Small", 19},
		{"Medium", 500},
		{"Large", 1000000},
		{"ExtraLarge", 10000000},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				sieve := NewSieve()
				_ = sieve.NthPrime(bm.index)
			}
		})
	}
}

func TestNthPrime(t *testing.T) {
	sieve := NewSieve()

	assert.Equal(t, int64(-1), sieve.NthPrime(-3)) // added for testing negative numbers
	assert.Equal(t, int64(2), sieve.NthPrime(0))
	assert.Equal(t, int64(71), sieve.NthPrime(19))
	assert.Equal(t, int64(541), sieve.NthPrime(99))
	assert.Equal(t, int64(3581), sieve.NthPrime(500))
	assert.Equal(t, int64(7793), sieve.NthPrime(986))
	assert.Equal(t, int64(17393), sieve.NthPrime(2000))
	assert.Equal(t, int64(15485867), sieve.NthPrime(1000000))
	assert.Equal(t, int64(179424691), sieve.NthPrime(10000000))
	assert.Equal(t, int64(2038074751), sieve.NthPrime(100000000)) // not required, just a fun challenge
}

func FuzzNthPrime(f *testing.F) {
	sieve := NewSieve()

	f.Fuzz(func(t *testing.T, n int64) {
		if !big.NewInt(sieve.NthPrime(n)).ProbablyPrime(0) {
			t.Errorf("the sieve produced a non-prime number at index %d", n)
		}
	})
}
