fn OneTwoThree() {
	yield 1
	yield 2
	yield 3
}

OneTwoThree() | List
-- [1, 2, 3]
