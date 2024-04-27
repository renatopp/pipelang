
fn fizzBuzz(x) {
	match x%3, x%5 {
		0, 0: "FizzBuzz",
		0, _: "Fizz",
		_, 0: "Buzz",
		_, _: String(x),
	}
}

range(100)
| map  fizzBuzz
| each x: println(x)
