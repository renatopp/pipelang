package internal

import "slices"

var Keywords = []string{
	"if",
	"else",
	"for",
	"with",
	"match",

	"return",
	"raise",
	"yield",
	"break",
	"continue",
	"defer",

	"true",
	"false",

	"data",
	"fn",
	"as",
	"is",
	"in",
}

var Assignments = []string{
	":=",
	"=",
	"+=",
	"-=",
	"*=",
	"/=",
	"%=",
	"^=",
	"..=",
}

var PrefixOperators = []string{
	"+",
	"-",
	"not",
	// "~",
}

var InfixOperators = []string{
	"+",
	"-",
	"*",
	"/",
	"%",
	"^",

	"==",
	"!=",
	"<",
	">",
	"<=",
	">=",
	"<=>",

	"??",
	"..",

	"and",
	"or",
	"xor",
}

var Operators = slices.Concat(PrefixOperators, InfixOperators)
