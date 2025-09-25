package sieve

import (
	"fmt"
	"math"
)

type Sieve interface {
	NthPrime(n int64) int64
}

type Prime struct {
	Primes []int64
}

func NewSieve() Sieve {
	return &Prime{}
}

// Generate all primes up to limit using Sieve of Eratosthenes (https://en.wikipedia.org/wiki/Sieve_of_Eratosthenes)
func generatePrimesUpTo(limit int64) []int64 {
	// There are no primes below 2 so return empty array
	if limit < 2 {
		return []int64{}
	}

	// Initialize array for sieving
	isPrime := make([]bool, limit+1)
	for i := 2; i <= int(limit); i++ {
		isPrime[i] = true
	}

	// This is the core sieve algorithm
	for i := 2; i*i <= int(limit); i++ {
		if isPrime[i] {
			for j := i * i; j <= int(limit); j += i {
				isPrime[j] = false
			}
		}
	}

	// Append all remaining true values, marking them as prime
	primes := make([]int64, 0)
	for i := 2; i <= int(limit); i++ {
		if isPrime[i] {
			primes = append(primes, int64(i))
		}
	}
	return primes
}

// NthPrime returns the nth prime (0-indexed)
func (p *Prime) NthPrime(n int64) int64 {
	// A quick validation check
	if n < 0 {
		fmt.Println("Number must be non-negative")
		return -1
	}

	// Rosser's theorem is not valid for this (n == 0), so return the first index
	if n == 0 {
		return 2
	}

	// Return if the number is already cached
	if int(n) < len(p.Primes) {
		return p.Primes[int(n)]
	}

	nf := float64(n + 1)
	// The following uses Rosser's theorem (https://en.wikipedia.org/wiki/Rosser%27s_theorem) for getting the upper limit
	maxPrime := int64(nf*(math.Log(nf)+math.Log(math.Log(nf)))+1) + 10
	segmentSize := int64(1_000_000) // 1 million bools per segment (~1MB per segment. bool in a slice == 1 byte)

	// Start from 2 unless primes have already been cached. If there's already a cache, resume just after the last cached.
	start := int64(2)
	if len(p.Primes) > 0 {
		start = p.Primes[len(p.Primes)-1] + 1
	}

	// Get base primes up to sqrt(maxPrimes)
	// - We only need primes up to the square root of maxPrime to find all primes. Anything over the sqrt(maxPrimes) adds computation and memory usage.
	sqrtMax := int64(math.Sqrt(float64(maxPrime))) + 1
	basePrimes := generatePrimesUpTo(sqrtMax)

	// Ensure p.Primes contains all base primes and append them if they don't exist
	for _, bp := range basePrimes {
		if len(p.Primes) == 0 || p.Primes[len(p.Primes)-1] < bp {
			p.Primes = append(p.Primes, bp)
		}
	}

	// Segmented sieve loop
	for low := start; low <= maxPrime; low += segmentSize {

		// This is for the last number in the segment
		high := low + segmentSize - 1

		// Make sure high doesn't exceed maxPrime
		if high > maxPrime {
			high = maxPrime
		}

		// Create the segment to be sieved by the base primes
		segment := make([]bool, high-low+1)
		for i := range segment {
			segment[i] = true
		}

		// Cross off multiples of base primes
		for _, bp := range basePrimes {
			// Find the first number in this segment that bp divides evenly so we can start marking multiples of bp as non-prime
			startIdx := ((low + bp - 1) / bp) * bp

			// Mark numbers in this segment as non-prime
			for j := startIdx; j <= high; j += bp {
				// Use j-low to adjust to the correct index since it is a segment
				segment[j-low] = false
			}
		}

		// Append primes from segment
		for i := low; i <= high; i++ {
			// Use [i-low] to adjust to the correct index since it is a segment
			if segment[i-low] {
				p.Primes = append(p.Primes, i)
				if int64(len(p.Primes)) > n {
					// return nth prime if found
					return p.Primes[int(n)]
				}
			}
		}
	}

	return p.Primes[int(n)]
}
