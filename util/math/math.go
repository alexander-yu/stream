package math

var factorials = []int{1, 1, 2, 6, 24, 120, 720, 5040}

func factorial(n int) int {
	for n >= len(factorials) {
		next := factorials[len(factorials)-1] * len(factorials)
		factorials = append(factorials, next)
	}

	return factorials[n]
}

// Sign returns the sign of an integer (-1 if negative, 1 otherwise).
func Sign(n int) int {
	if n%2 == 0 {
		return 1
	}

	return -1
}

// Binom returns the binomial coefficient.
func Binom(n, k int) int {
	if k == 0 || k == n {
		return 1
	}

	return factorial(n) / (factorial(k) * factorial(n-k))
}
