solution := range()
| map       fib
| takeWhile x: x <= 4_000_000
| filter    x: x.IsEven()
| sum

-- Equivalent to:
sum(
  filter(
    takeWhile(
      map(range(), fib),
      x: x <= 4_000_000
    ),
    x => x.IsEven()
  )
)