-- Problem 1: Find the sum of all the multiples of 3 or 5 below 1000
solution := range(1000)
| filter num: num%3==0 or num%5==0
| sum

println(solution)