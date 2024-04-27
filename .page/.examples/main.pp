fn fib(x) {
	match x {
		0: 0,
		1: 1,
		_: fib(x - 1) + fib(x - 2)
	}
}

solution := range()
| map       fib
| takeWhile x: x <= 4_000_000
| filter    x: x.IsEven()
| sum

println(solution)
