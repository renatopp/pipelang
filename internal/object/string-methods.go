package object

import (
	"regexp"
	"slices"
	"strings"
	"unicode/utf8"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func str_len(s string) int {
	return utf8.RuneCountInString(s)
}

func init() {
	StringTypeObj.AddMethod(String_Size)
	StringTypeObj.AddMethod(String_Get)
	StringTypeObj.AddMethod(String_GetOr)
	StringTypeObj.AddMethod(String_Sub)
	StringTypeObj.AddMethod(String_Cut)
	StringTypeObj.AddMethod(String_CutFn)
	StringTypeObj.AddMethod(String_Find)
	StringTypeObj.AddMethod(String_FindAny)
	StringTypeObj.AddMethod(String_FindFn)
	StringTypeObj.AddMethod(String_FindLast)
	StringTypeObj.AddMethod(String_FindLastAny)
	StringTypeObj.AddMethod(String_FindLastFn)
	StringTypeObj.AddMethod(String_Join)
	StringTypeObj.AddMethod(String_JoinArgs)
	StringTypeObj.AddMethod(String_Split)
	StringTypeObj.AddMethod(String_SplitAt)
	StringTypeObj.AddMethod(String_SplitFn)
	StringTypeObj.AddMethod(String_Fields)
	StringTypeObj.AddMethod(String_ToLower)
	StringTypeObj.AddMethod(String_ToUpper)
	StringTypeObj.AddMethod(String_ToTitle)
	StringTypeObj.AddMethod(String_ToSnake)
	StringTypeObj.AddMethod(String_ToKebab)
	StringTypeObj.AddMethod(String_ToCamel)
	StringTypeObj.AddMethod(String_ToPascal)
	StringTypeObj.AddMethod(String_ToDot)
	StringTypeObj.AddMethod(String_ToTrain)
	StringTypeObj.AddMethod(String_ToSentence)
	StringTypeObj.AddMethod(String_Ellipsis)
	StringTypeObj.AddMethod(String_PadCenter)
	StringTypeObj.AddMethod(String_PadCenterWith)
	StringTypeObj.AddMethod(String_PadLeft)
	StringTypeObj.AddMethod(String_PadLeftWith)
	StringTypeObj.AddMethod(String_PadRight)
	StringTypeObj.AddMethod(String_PadRightWith)
	StringTypeObj.AddMethod(String_Repeat)
	StringTypeObj.AddMethod(String_Reverse)
	StringTypeObj.AddMethod(String_TrimSpace)
	StringTypeObj.AddMethod(String_TrimSpaceLeft)
	StringTypeObj.AddMethod(String_TrimSpaceRight)
	StringTypeObj.AddMethod(String_Trim)
	StringTypeObj.AddMethod(String_TrimLeft)
	StringTypeObj.AddMethod(String_TrimRight)
	StringTypeObj.AddMethod(String_Replace)
	StringTypeObj.AddMethod(String_ReplaceN)
	StringTypeObj.AddMethod(String_Sort)
	StringTypeObj.AddMethod(String_SortFn)
	StringTypeObj.AddMethod(String_Contains)
	StringTypeObj.AddMethod(String_ContainsChars)
	StringTypeObj.AddMethod(String_StartsWith)
	StringTypeObj.AddMethod(String_EndsWith)
	StringTypeObj.AddMethod(String_IsEmpty)
	StringTypeObj.AddMethod(String_Chars)
	StringTypeObj.AddMethod(String_Lines)
	StringTypeObj.AddMethod(String_Words)
}

// ----------------------------------------------------------------------------
// Accessors
// ----------------------------------------------------------------------------
var String_Size = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*String)
		return NewNumber(float64(str_len(this.Value)))
	},
	`Size`,
	`Returns the number of characters in the string.`,
	P("this", V.Type(StringId)),
)

var String_Get = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*String)
		index := int(args[1].(*Number).Value)
		len := str_len(this.Value)

		if index < 0 {
			index = len + index
		}

		if index < 0 || index >= len {
			return scope.Interrupt(Raise("string index '%d' out of bounds", index))
		}

		runes := []rune(this.Value)
		return NewString(string(runes[index]))
	},
	`Get`,
	`Returns the character at the given index. Raise an error if the index is out of bounds.`,
	P("this", V.Type(StringId)),
	P("index", V.Type(NumberId)),
)

var String_GetOr = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*String)
		index := int(args[1].(*Number).Value)
		len := str_len(this.Value)

		if index < 0 {
			index = len + index
		}

		if index < 0 || index >= len {
			return args[2]
		}

		runes := []rune(this.Value)
		return NewString(string(runes[index]))
	},
	`GetOr`,
	`Returns the character at the given index. Returns the default value if the
index is out of bounds.`,
	P("this", V.Type(StringId)),
	P("index", V.Type(NumberId)),
	P("default", V.Type(StringId)),
)

var String_Sub = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*String)
		from := int(args[1].(*Number).Value)
		to := int(args[2].(*Number).Value)
		len := str_len(this.Value)

		if from < 0 {
			from = len + from
		}

		if to < 0 {
			to = len + to
		}

		if to < from || to < 0 || from >= len {
			return EmptyString
		}

		from = max(0, from)
		to = min(len, to)
		runes := []rune(this.Value)
		return NewString(string(runes[from:to]))
	},
	`Sub`,
	`Returns a substring starting from the given index up to the end index. Returns empty if the range is completely out of bounds.
	
You may use negative indexes to start from the end of the string.
	`,
	P("this", V.Type(StringId)),
	P("from", V.Type(NumberId)),
	P("to", V.Type(NumberId)),
)

var String_Cut = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*String)
		index := int(args[1].(*Number).Value)
		len := str_len(this.Value)

		if index < 0 {
			index = len + index
		}

		if index < 0 {
			return NewTuple(EmptyString, this.Copy())
		} else if index >= len {
			return NewTuple(this.Copy(), EmptyString)
		}

		runes := []rune(this.Value)
		return NewTuple(
			NewString(string(runes[:index])),
			NewString(string(runes[index:])),
		)
	},
	`Cut`,
	`Divide the string into two parts at the given index. Returns a tuple with the two parts.`,
	P("this", V.Type(StringId)),
	P("index", V.Type(NumberId)),
)

var String_CutFn = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*String)
		f := args[1].(*Function)

		runes := []rune(this.Value)
		for idx, chr := range runes {
			ret := scope.Eval().Call(scope, f, []Object{
				NewString(string(chr)),
				NewNumber(float64(idx)),
			})
			if isRaise(ret) {
				return ret
			}

			if ret.AsBool() {
				return NewTuple(
					NewString(string(runes[:idx])),
					NewString(string(runes[idx:])),
				)
			}
		}

		return NewTuple(this.Copy(), EmptyString)
	},
	`CutFn`,
	`Divide the string into two parts when the provided function returns true. Returns a tuple with the two parts.`,
	P("this", V.Type(StringId)),
	P("f", V.Type(FunctionId)),
)

var String_Find = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*String)
		sub := args[1].(*String)
		runes := []rune(this.Value)
		idx := 0
		for i, chr := range runes {
			if strings.HasPrefix(this.Value[idx:], sub.Value) {
				return NewNumber(float64(i))
			}
			idx += utf8.RuneLen(chr)
		}
		return NewNumber(float64(-1))
	},
	`Find`,
	`Find the first occurrence of the substring in the string. Returns the index of the first character of the substring, or -1 if not found.`,
	P("this", V.Type(StringId)),
	P("sub", V.Type(StringId)),
)

var String_FindAny = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*String)
		sub := args[1].(*String)
		runes := []rune(this.Value)
		idx := 0
		for i, chr := range runes {
			if strings.ContainsRune(sub.Value, chr) {
				return NewNumber(float64(i))
			}
			idx += utf8.RuneLen(chr)
		}
		return NewNumber(float64(-1))
	},
	`FindAny`,
	`Find the first occurrence of any of the characters in the string. Returns the index of the first character, or -1 if not found.`,
	P("this", V.Type(StringId)),
	P("sub", V.Type(StringId)),
)

var String_FindFn = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*String)
		f := args[1].(*Function)
		runes := []rune(this.Value)
		for idx, chr := range runes {
			ret := scope.Eval().Call(scope, f, []Object{
				NewString(string(chr)),
				NewNumber(float64(idx)),
			})
			if isRaise(ret) {
				return ret
			}

			if ret.AsBool() {
				return NewNumber(float64(idx))
			}
		}

		return NewNumber(-1)
	},
	`FindFn`,
	`Find the first occurrence of a character in the string that satisfies the function. Returns the index of the character, or -1 if not found.`,
	P("this", V.Type(StringId)),
	P("f", V.Type(FunctionId)),
)

var String_FindLast = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*String)
		sub := args[1].(*String)
		runes := []rune(this.Value)
		idx := len(this.Value)
		for i := len(runes) - 1; i >= 0; i-- {
			chr := runes[i]
			idx -= utf8.RuneLen(chr)
			if strings.HasPrefix(this.Value[idx:], sub.Value) {
				return NewNumber(float64(i))
			}
		}
		return NewNumber(float64(-1))
	},
	`FindLast`,
	`Find the last occurrence of the substring in the string. Returns the index of the first character of the substring, or -1 if not found.`,
	P("this", V.Type(StringId)),
	P("sub", V.Type(StringId)),
)

var String_FindLastAny = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*String)
		sub := args[1].(*String)
		runes := []rune(this.Value)
		for i := len(runes) - 1; i >= 0; i-- {
			chr := runes[i]
			if strings.ContainsRune(sub.Value, chr) {
				return NewNumber(float64(i))
			}
		}
		return NewNumber(float64(-1))
	},
	`FindLastAny`,
	`Find the last occurrence of any of the characters in the string. Returns the index of the first character, or -1 if not found.`,
	P("this", V.Type(StringId)),
	P("sub", V.Type(StringId)),
)

var String_FindLastFn = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*String)
		f := args[1].(*Function)
		runes := []rune(this.Value)
		for idx := len(runes) - 1; idx >= 0; idx-- {
			chr := runes[idx]
			ret := scope.Eval().Call(scope, f, []Object{
				NewString(string(chr)),
				NewNumber(float64(idx)),
			})
			if isRaise(ret) {
				return ret
			}

			if ret.AsBool() {
				return NewNumber(float64(idx))
			}
		}

		return NewNumber(-1)
	},
	`FindLastFn`,
	`Find the last occurrence of a character in the string that satisfies the function. Returns the index of the character, or -1 if not found.`,
	P("this", V.Type(StringId)),
	P("f", V.Type(FunctionId)),
)

// ----------------------------------------------------------------------------
// Conversions
// ----------------------------------------------------------------------------
var String_Join = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*String)
		list := args[1].(*List)
		elements := make([]string, len(list.Elements))
		for i, e := range list.Elements {
			elements[i] = e.AsString()
		}
		return NewString(string(strings.Join(elements, this.Value)))
	},
	`Join`,
	`Join the elements of a list into a new string, using the current string as the separator.`,
	P("this", V.Type(StringId)),
	P("list", V.Type(ListId)),
)

var String_JoinArgs = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*String)
		elements := make([]string, len(args)-1)
		for i, e := range args[1:] {
			elements[i] = e.AsString()
		}
		return NewString(string(strings.Join(elements, this.Value)))
	},
	`JoinArgs`,
	`Join the elements of a list into a new string, using the current string as the separator.`,
	P("this", V.Type(StringId)),
	P("args").AsSpread(),
)

var String_Split = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*String)
		sub := args[1].(*String)

		items := strings.Split(this.Value, sub.Value)
		objs := make([]Object, len(items))

		for i, item := range items {
			objs[i] = NewString(item)
		}

		return NewList(objs...)
	},
	`Split`,
	`Split the string into a list of substrings, given the separator.`,
	P("this", V.Type(StringId)),
	P("sub", V.Type(StringId)),
)

var String_SplitAt = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*String)
		index := int(args[1].(*Number).Value)
		len := str_len(this.Value)
		runes := []rune(this.Value)

		if index < 0 {
			index = len + index
		}

		if index < 0 {
			return NewList(EmptyString, this.Copy())
		} else if index >= len {
			return NewList(this.Copy(), EmptyString)
		}

		return NewList(
			NewString(string(runes[:index])),
			NewString(string(runes[index:])),
		)
	},
	`SplitAt`,
	`Split the string into a list of substrings, given the index.`,
	P("this", V.Type(StringId)),
	P("index", V.Type(NumberId)),
)

var String_SplitFn = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*String)
		f := args[1].(*Function)

		parts := []Object{}
		lastIdx := 0
		runes := []rune(this.Value)
		for idx, chr := range runes {
			ret := scope.Eval().Call(scope, f, []Object{
				NewString(string(chr)),
				NewNumber(float64(idx)),
			})
			if isRaise(ret) {
				return ret
			}

			if ret.AsBool() {
				parts = append(parts, NewString(string(runes[lastIdx:idx])))
				lastIdx = idx
			}
		}

		if lastIdx < len(this.Value) {
			parts = append(parts, NewString(string(runes[lastIdx:])))
		}

		return NewList(parts...)
	},
	`SplitFn`,
	`Split the string into a list of substrings, given the function.`,
	P("this", V.Type(StringId)),
	P("f", V.Type(FunctionId)),
)

var String_Fields = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*String)

		parts := strings.Fields(this.Value)
		objs := make([]Object, len(parts))
		for i, part := range parts {
			objs[i] = NewString(part)
		}

		return NewList(objs...)
	},
	`Fields`,
	`Split the string into a list of substrings, using white spaces as separators.`,
	P("this", V.Type(StringId)),
)

// ----------------------------------------------------------------------------
// Casings
// ----------------------------------------------------------------------------
var toLower = cases.Lower(language.English)
var String_ToLower = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*String)
		return NewString(toLower.String(this.Value))
	},
	`ToLower`,
	`Converts the string to lower case.`,
	P("this", V.Type(StringId)),
)

var toUpper = cases.Upper(language.English)
var String_ToUpper = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*String)
		return NewString(toUpper.String(this.Value))
	},
	`ToUpper`,
	`Converts the string to upper case.`,
	P("this", V.Type(StringId)),
)

var toTitle = cases.Title(language.English)
var String_ToTitle = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*String)
		return NewString(toTitle.String(this.Value))
	},
	`ToTitle`,
	`Converts the string to title case.`,
	P("this", V.Type(StringId)),
)

var snakeReplacer = regexp.MustCompile(`[^a-zA-Z0-9]+`)
var String_ToSnake = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*String)
		s := snakeReplacer.ReplaceAllString(this.Value, " ")
		words := strings.Fields(s)
		return NewString(strings.ToLower(strings.Join(words, "_")))
	},
	`ToSnake`,
	`Converts the string to snake case.`,
	P("this", V.Type(StringId)),
)

var String_ToKebab = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*String)
		s := snakeReplacer.ReplaceAllString(this.Value, " ")
		words := strings.Fields(s)
		return NewString(strings.ToLower(strings.Join(words, "-")))
	},
	`ToKebab`,
	`Converts the string to kebab case.`,
	P("this", V.Type(StringId)),
)

var String_ToCamel = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*String)
		s := snakeReplacer.ReplaceAllString(this.Value, " ")
		words := strings.Fields(s)
		if len(words) == 0 {
			return EmptyString
		}
		words[0] = toLower.String(words[0])
		for i, word := range words[1:] {
			words[i+1] = toTitle.String(word)
		}
		return NewString(strings.Join(words, ""))
	},
	`ToCamel`,
	`Converts the string to camel case.`,
	P("this", V.Type(StringId)),
)

var String_ToPascal = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*String)
		s := snakeReplacer.ReplaceAllString(this.Value, " ")
		words := strings.Fields(s)
		for i, word := range words {
			words[i] = toTitle.String(word)
		}
		return NewString(strings.Join(words, ""))
	},
	`ToPascal`,
	`Converts the string to pascal case.`,
	P("this", V.Type(StringId)),
)

var String_ToDot = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*String)
		s := snakeReplacer.ReplaceAllString(this.Value, " ")
		words := strings.Fields(s)
		return NewString(strings.ToLower(strings.Join(words, ".")))
	},
	`ToDot`,
	`Converts the string to dot case.`,
	P("this", V.Type(StringId)),
)

var String_ToTrain = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*String)
		s := snakeReplacer.ReplaceAllString(this.Value, " ")
		words := strings.Fields(s)
		for i, word := range words {
			words[i] = toTitle.String(word)
		}
		return NewString(strings.Join(words, "-"))
	},
	`ToTrain`,
	`Converts the string to train case.`,
	P("this", V.Type(StringId)),
)

var String_ToSentence = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*String)
		if len(this.Value) == 0 {
			return EmptyString
		}

		firstRune, size := utf8.DecodeRuneInString(this.Value)
		word := string(toUpper.String(string(firstRune))) + this.Value[size:]

		return NewString(word)
	},
	`ToSentence`,
	`Converts the string to sentence case.`,
	P("this", V.Type(StringId)),
)

// ----------------------------------------------------------------------------
// Formats
// ----------------------------------------------------------------------------
var String_Ellipsis = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*String)
		max := int(args[1].(*Number).Value)
		if str_len(this.Value) <= max {
			return this.Copy()
		}
		runes := []rune(this.Value)
		return NewString(string(runes[:max]) + "...")
	},
	`Ellipsis`,
	`Truncate the string to the given length and append an ellipsis.`,
	P("this", V.Type(StringId)),
	P("max", V.Type(NumberId)),
)

var String_PadCenter = F(
	func(scope *Scope, args ...Object) Object {
		return String_PadCenterWith.Call(scope, args[0], args[1], NewString(" "))
	},
	`PadCenter`,
	`Pad the string with spaces to the given size, centering the content.`,
	P("this", V.Type(StringId)),
	P("size", V.Type(NumberId)),
)

var String_PadCenterWith = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*String)
		size := int(args[1].(*Number).Value)
		pad := args[2].(*String).Value
		len := str_len(this.Value)

		if len >= size {
			return this.Copy()
		}

		left := (size - len) / 2
		right := size - len - left
		return NewString(strings.Repeat(pad, left) + this.Value + strings.Repeat(pad, right))
	},
	`PadCenterWith`,
	`Pad the string with the provided pad string to the given size, centering the content.`,
	P("this", V.Type(StringId)),
	P("size", V.Type(NumberId)),
	P("pad", V.Type(StringId)),
)

var String_PadLeft = F(
	func(scope *Scope, args ...Object) Object {
		return String_PadLeftWith.Call(scope, args[0], args[1], NewString(" "))
	},
	`PadLeft`,
	`Pad the string with spaces to the given size, aligning the content to the right.`,
	P("this", V.Type(StringId)),
	P("size", V.Type(NumberId)),
)

var String_PadLeftWith = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*String)
		size := int(args[1].(*Number).Value)
		pad := args[2].(*String).Value
		len := str_len(this.Value)

		if len >= size {
			return this.Copy()
		}

		return NewString(strings.Repeat(pad, size-len) + this.Value)
	},
	`PadLeftWith`,
	`Pad the string with the provided pad string to the given size, aligning the content to the right.`,
	P("this", V.Type(StringId)),
	P("size", V.Type(NumberId)),
	P("pad", V.Type(StringId)),
)

var String_PadRight = F(
	func(scope *Scope, args ...Object) Object {
		return String_PadRightWith.Call(scope, args[0], args[1], NewString(" "))
	},
	`PadRight`,
	`Pad the string with spaces to the given size, aligning the content to the left.`,
	P("this", V.Type(StringId)),
	P("size", V.Type(NumberId)),
)

var String_PadRightWith = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*String)
		size := int(args[1].(*Number).Value)
		pad := args[2].(*String).Value
		len := str_len(this.Value)

		if len >= size {
			return this.Copy()
		}

		return NewString(this.Value + strings.Repeat(pad, size-len))
	},
	`PadRightWith`,
	`Pad the string with the provided pad string to the given size, aligning the content to the left.`,
	P("this", V.Type(StringId)),
	P("size", V.Type(NumberId)),
	P("pad", V.Type(StringId)),
)

var String_Repeat = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*String)
		times := int(args[1].(*Number).Value)
		return NewString(strings.Repeat(this.Value, times))
	},
	`Repeat`,
	`Repeat the string the given number of times.`,
	P("this", V.Type(StringId)),
	P("times", V.Type(NumberId)),
)

var String_Reverse = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*String)
		runes := []rune(this.Value)
		for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
			runes[i], runes[j] = runes[j], runes[i]
		}
		return NewString(string(runes))
	},
	`Reverse`,
	`Reverse the string.`,
	P("this", V.Type(StringId)),
)

var String_TrimSpace = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*String)
		return NewString(strings.TrimSpace(this.Value))
	},
	`TrimSpace`,
	`Remove leading and trailing white spaces from the string.`,
	P("this", V.Type(StringId)),
)

var String_TrimSpaceLeft = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*String)
		return NewString(strings.TrimLeft(this.Value, " \t\n\r"))
	},
	`TrimSpaceLeft`,
	`Remove leading white spaces from the string.`,
	P("this", V.Type(StringId)),
)

var String_TrimSpaceRight = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*String)
		return NewString(strings.TrimRight(this.Value, " \t\n\r"))
	},
	`TrimSpaceRight`,
	`Remove trailing white spaces from the string.`,
	P("this", V.Type(StringId)),
)

var String_Trim = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*String)
		sub := args[1].(*String)
		return NewString(strings.Trim(this.Value, sub.Value))
	},
	`Trim`,
	`Remove the given characters from the beginning and end of the string.`,
	P("this", V.Type(StringId)),
	P("sub", V.Type(StringId)),
)

var String_TrimLeft = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*String)
		sub := args[1].(*String)
		return NewString(strings.TrimLeft(this.Value, sub.Value))
	},
	`TrimLeft`,
	`Remove the given characters from the beginning of the string.`,
	P("this", V.Type(StringId)),
	P("sub", V.Type(StringId)),
)

var String_TrimRight = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*String)
		sub := args[1].(*String)
		return NewString(strings.TrimRight(this.Value, sub.Value))
	},
	`TrimRight`,
	`Remove the given characters from the end of the string.`,
	P("this", V.Type(StringId)),
	P("sub", V.Type(StringId)),
)

var String_Replace = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*String)
		old := args[1].(*String)
		new := args[2].(*String)
		return NewString(strings.ReplaceAll(this.Value, old.Value, new.Value))
	},
	`Replace`,
	`Replace all occurrences of the old substring with the new one.`,
	P("this", V.Type(StringId)),
	P("old", V.Type(StringId)),
	P("new", V.Type(StringId)),
)

var String_ReplaceN = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*String)
		old := args[1].(*String)
		new := args[2].(*String)
		n := int(args[3].(*Number).Value)
		return NewString(strings.Replace(this.Value, old.Value, new.Value, n))
	},
	`ReplaceN`,
	`Replace the first n occurrences of the old substring with the new one.`,
	P("this", V.Type(StringId)),
	P("old", V.Type(StringId)),
	P("new", V.Type(StringId)),
	P("n", V.Type(NumberId)),
)

var String_Sort = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*String)
		runes := []rune(this.Value)
		slices.Sort(runes)
		return NewString(string(runes))
	},
	`Sort`,
	`Sort the characters of the string.`,
	P("this", V.Type(StringId)),
)

var String_SortFn = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*String)
		fn := args[1].(*Function)
		runes := []rune(this.Value)
		slices.SortFunc(runes, func(a, b rune) int {
			ret := scope.Eval().Call(scope, fn, []Object{
				NewString(string(a)),
				NewString(string(b)),
			})
			if isRaise(ret) {
				return 0
			}

			n, ok := ret.(*Number)
			if !ok {
				return 0
			}

			return int(n.Value)
		})
		return NewString(string(runes))
	},
	`SortFn`,
	`Sort the characters of the string given the function.`,
	P("this", V.Type(StringId)),
	P("f", V.Type(FunctionId)),
)

// ----------------------------------------------------------------------------
// Checkers
// ----------------------------------------------------------------------------
var String_Contains = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*String)
		searches := args[1:]
		for _, e := range searches {
			search := e.(*String)
			if strings.Contains(this.Value, search.Value) {
				return True
			}
		}
		return False
	},
	`Contains`,
	`Check if the string contains any of the given substrings.`,
	P("this", V.Type(StringId)),
	P("search", V.Type(StringId)).AsSpread(),
)

var String_ContainsChars = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*String)
		search := args[1].(*String)
		return NewBoolean(strings.ContainsAny(this.Value, search.Value))
	},
	`ContainsChars`,
	`Check if the string contains any of the characters in the given string.`,
	P("this", V.Type(StringId)),
	P("search", V.Type(StringId)),
)

var String_StartsWith = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*String)
		searches := args[1:]
		for _, e := range searches {
			search := e.(*String)
			if strings.HasPrefix(this.Value, search.Value) {
				return True
			}
		}
		return False
	},
	`StartsWith`,
	`Check if the string starts with any of the given substrings.`,
	P("this", V.Type(StringId)),
	P("search", V.Type(StringId)).AsSpread(),
)

var String_EndsWith = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*String)
		searches := args[1:]
		for _, e := range searches {
			search := e.(*String)
			if strings.HasSuffix(this.Value, search.Value) {
				return True
			}
		}
		return False
	},
	`EndsWith`,
	`Check if the string ends with any of the given substrings.`,
	P("this", V.Type(StringId)),
	P("search", V.Type(StringId)).AsSpread(),
)

var String_IsEmpty = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*String)
		return NewBoolean(this.Value == "")
	},
	`IsEmpty`,
	`Check if the string is empty.`,
	P("this", V.Type(StringId)),
)

// ----------------------------------------------------------------------------
// Streams
// ----------------------------------------------------------------------------
var String_Chars = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*String)
		runes := []rune(this.Value)
		idx := 0
		return NewInternalStream(func(s *Scope) Object {
			if idx >= len(runes) {
				return nil
			}
			r := runes[idx]
			idx++
			return YieldWith(NewString(string(r)))
		}, scope)
	},
	`Chars`,
	`Create a stream of characters from the string.`,
	P("this", V.Type(StringId)),
)

var String_Lines = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*String)
		runes := []rune(this.Value)
		pivot := 0
		idx := 0
		return NewInternalStream(func(s *Scope) Object {
			for {
				if idx >= len(runes) {
					if pivot < idx {
						line := runes[pivot:]
						pivot = idx + 1
						return YieldWith(NewString(string(line)))
					} else if pivot == idx {
						pivot = idx + 1
						return YieldWith(EmptyString)
					}
					return nil
				}

				chr := runes[idx]
				idx++

				if chr == '\n' {
					line := runes[pivot : idx-1]
					pivot = idx
					return YieldWith(NewString(string(line)))
				}
			}
		}, scope)
	},
	`Lines`,
	`Create a stream of lines from the string.`,
	P("this", V.Type(StringId)),
)

var String_Words = F(
	func(scope *Scope, args ...Object) Object {
		this := args[0].(*String)
		runes := []rune(this.Value)
		pivot := 0
		idx := 0
		return NewInternalStream(func(s *Scope) Object {
			for {
				if idx >= len(runes) {
					if pivot < idx {
						word := runes[pivot:]
						pivot = idx + 1
						return YieldWith(NewString(string(word)))
					}
					return nil
				}

				chr := runes[idx]
				idx++

				if chr == ' ' || chr == '\t' || chr == '\n' || chr == '\r' {
					if pivot < idx-1 {
						word := runes[pivot : idx-1]
						pivot = idx
						println(string(word))
						return YieldWith(NewString(string(word)))
					}
					pivot = idx
				}
			}
		}, scope)
	},
	`Words`,
	`Create a stream of words from the string.`,
	P("this", V.Type(StringId)),
)
