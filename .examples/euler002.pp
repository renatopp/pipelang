-- Problem 2: Find the sum of the even-valued terms in the Fibonacci sequence whose values do not exceed four million
fn fib(x) {
  return match x {
    0: 0;
		1: 1
    x: fib(x-1) + fib(x-2)
  }
}

range()
| map       fib
| takeWhile x: x <= 4_000_000
| filter    x: x%2 == 0
| sum