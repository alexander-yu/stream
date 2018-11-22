package moment

var factorials = []int{1, 1, 2, 6, 24, 120, 720, 5040}

func factorial(n int) int {
	for n >= len(factorials) {
		next := factorials[len(factorials)-1] * len(factorials)
		factorials = append(factorials, next)
	}

	return factorials[n]
}

func sign(n int) int {
	if n%2 == 0 {
		return 1
	}

	return -1
}

func binom(n, k int) int {
	if k == 0 || k == n {
		return 1
	}

	return factorial(n) / (factorial(k) * factorial(n-k))
}
