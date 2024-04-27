


fn diff(n) {
	squareSum := range(n+1) | map x: x^2 | sum
	sumSquare := range(n+1) | sum | map x: x^2 | Number
	return sumSquare - squareSum
}
