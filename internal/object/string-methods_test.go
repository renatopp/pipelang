package object_test

import (
	"testing"

	"github.com/renatopp/pipelang/test/common"
)

func TestString_Index(t *testing.T) {
	common.AssertCode(t, `'renato'[0]`, `r`)
	common.AssertCode(t, `'renato'[-1]`, `o`)
	common.AssertCode(t, `'renato'[0, 5]`, `renat`)
	common.AssertCode(t, `'renato'[3, 5]`, `at`)
	common.AssertCode(t, `'áçãẞ🥸'[1, 3]`, `çã`)
	common.AssertCode(t, `'renato'[4, 10]`, `to`)
	common.AssertCode(t, `'renato'[-2, 10]`, `to`)
	common.AssertCode(t, `'renato'[-100, 0]`, ``)
	common.AssertCodeError(t, `'renato'[6]`)
	common.AssertCodeError(t, `'renato'[-7]`)
	common.AssertCodeError(t, `'renato'[-7, 30, 20]`)
}

func TestString_Size(t *testing.T) {
	common.AssertCode(t, `'renato'.Size()`, `6`)
	common.AssertCode(t, `''.Size()`, `0`)
	common.AssertCode(t, `'çãí🥸'.Size()`, `4`)
}

func TestString_Get(t *testing.T) {
	common.AssertCode(t, `'áçãẞ🥸'.Get(3)`, `ẞ`)
	common.AssertCode(t, `'renato'.Get(0)`, `r`)
	common.AssertCode(t, `'renato'.Get(-1)`, `o`)
	common.AssertCodeError(t, `'renato'.Get()`)
	common.AssertCodeError(t, `'renato'.Get(6)`)
	common.AssertCodeError(t, `'renato'.Get(-7)`)
	common.AssertCodeError(t, `'renato'.Get('a')`)
}

func TestString_GetOr(t *testing.T) {
	common.AssertCode(t, `'áçãẞ🥸'.GetOr(2, 'a')`, `ã`)
	common.AssertCode(t, `'renato'.GetOr(0, 'a')`, `r`)
	common.AssertCode(t, `'renato'.GetOr(100, 'a')`, `a`)
	common.AssertCodeError(t, `'renato'.GetOr(0)`)
	common.AssertCodeError(t, `'renato'.GetOr(0, 3)`)
	common.AssertCodeError(t, `'renato'.GetOr('a', 'a')`)
}

func TestString_Sub(t *testing.T) {
	common.AssertCode(t, `'áçãẞ🥸'.Sub(3, 10)`, `ẞ🥸`)
	common.AssertCode(t, `'renato'.Sub(0, 5)`, `renat`)
	common.AssertCode(t, `'renato'.Sub(3, 5)`, `at`)
	common.AssertCode(t, `'renato'.Sub(4, 10)`, `to`)
	common.AssertCode(t, `'renato'.Sub(-2, 10)`, `to`)
	common.AssertCode(t, `'renato'.Sub(-100, 0)`, ``)
	common.AssertCodeError(t, `'renato'.Sub(-1000)`)
	common.AssertCodeError(t, `'renato'.Sub('a')`)
}

func TestString_Cut(t *testing.T) {
	common.AssertCode(t, `'áçãẞ🥸'.Cut(3)`, `('áçã', 'ẞ🥸')`)
	common.AssertCode(t, `'renato'.Cut(1)`, `('r', 'enato')`)
	common.AssertCode(t, `'renato'.Cut(2)`, `('re', 'nato')`)
	common.AssertCode(t, `'renato'.Cut(-1)`, `('renat', 'o')`)
	common.AssertCode(t, `'renato'.Cut(-100)`, `('', 'renato')`)
	common.AssertCode(t, `'renato'.Cut(100)`, `('renato', '')`)
}

func TestString_CutFn(t *testing.T) {
	common.AssertCode(t, `'áçãẞ🥸'.CutFn(chr: chr == 'ẞ')`, `('áçã', 'ẞ🥸')`)
	common.AssertCode(t, `'renato'.CutFn(chr: chr == 'a')`, `('ren', 'ato')`)
	common.AssertCode(t, `'renato'.CutFn((_, idx): idx >= 3)`, `('ren', 'ato')`)
}

func TestString_Find(t *testing.T) {
	common.AssertCode(t, `'áçãẞ🥸'.Find('🥸')`, `4`)
	common.AssertCode(t, `'renato nato'.Find('a')`, `3`)
	common.AssertCode(t, `'renato nato'.Find('an')`, `-1`)
	common.AssertCode(t, `'renato nato'.Find('na')`, `2`)
}

func TestString_FindAny(t *testing.T) {
	common.AssertCode(t, `'áçãẞ🥸'.FindAny('cã2')`, `2`)
	common.AssertCode(t, `'renato nato'.FindAny('a')`, `3`)
	common.AssertCode(t, `'renato nato'.FindAny('an')`, `2`)
	common.AssertCode(t, `'renato nato'.FindAny('na')`, `2`)
}

func TestString_FindFn(t *testing.T) {
	common.AssertCode(t, `'áçãẞ🥸'.FindFn(chr: chr == 'ã')`, `2`)
	common.AssertCode(t, `'renato nato'.FindFn(chr: chr == 'a')`, `3`)
}

func TestString_FindLast(t *testing.T) {
	common.AssertCode(t, `'áçãẞ🥸'.FindLast('ẞ')`, `3`)
	common.AssertCode(t, `'renato nato'.FindLast('a')`, `8`)
	common.AssertCode(t, `'renato nato'.FindLast('an')`, `-1`)
	common.AssertCode(t, `'renato nato'.FindLast('na')`, `7`)
}

func TestString_FindLastAny(t *testing.T) {
	common.AssertCode(t, `'áçãẞ🥸'.FindLastAny('fẞa')`, `3`)
	common.AssertCode(t, `'renato nato'.FindLastAny('a')`, `8`)
	common.AssertCode(t, `'renato nato'.FindLastAny('an')`, `8`)
	common.AssertCode(t, `'renato nato'.FindLastAny('na')`, `8`)
}

func TestString_FindLastFn(t *testing.T) {
	common.AssertCode(t, `'áçãẞ🥸'.FindLastFn(chr: chr == 'ẞ')`, `3`)
	common.AssertCode(t, `'renato nato'.FindLastFn(chr: chr == 'a')`, `8`)
}

func TestString_Join(t *testing.T) {
	common.AssertCode(t, `'🎁'.Join([1, 2, 3])`, `1🎁2🎁3`)
	common.AssertCode(t, `', '.Join([])`, ``)
	common.AssertCode(t, `', '.Join([2])`, `2`)
	common.AssertCode(t, `', '.Join([1, 2, 'c'])`, `1, 2, c`)
	common.AssertCode(t, `'; '.Join([1, 2, 'c'])`, `1; 2; c`)
	common.AssertCode(t, `'| '.Join([1, 2, 'c'])`, `1| 2| c`)
	common.AssertCode(t, `','.Join([1, 2, 'c'])`, `1,2,c`)
	common.AssertCodeError(t, `', '.Join("a")`)
}

func TestString_JoinArgs(t *testing.T) {
	common.AssertCode(t, `'🎁'.JoinArgs(1, 2, 3)`, `1🎁2🎁3`)
	common.AssertCode(t, `', '.JoinArgs(1, 2, 'c')`, `1, 2, c`)
	common.AssertCode(t, `'; '.JoinArgs(1, 2, 'c')`, `1; 2; c`)
	common.AssertCode(t, `'| '.JoinArgs(1, 2, 'c')`, `1| 2| c`)
	common.AssertCode(t, `','.JoinArgs(1, 2, 'c')`, `1,2,c`)
	common.AssertCode(t, `','.JoinArgs('c')`, `c`)
}

func TestString_Split(t *testing.T) {
	common.AssertCode(t, `'1🎁2🎁3'.Split('🎁')`, `['1', '2', '3']`)
	common.AssertCode(t, `'hello,world'.Split(',')`, `['hello', 'world']`)
	common.AssertCode(t, `'hellopingworld'.Split('ping')`, `['hello', 'world']`)
	common.AssertCode(t, `'abcdefg'.Split('')`, `['a', 'b', 'c', 'd', 'e', 'f', 'g']`)
	common.AssertCode(t, `'.1.1.'.Split('.')`, `['', '1', '1', '']`)
	common.AssertCode(t, `''.Split('')`, `[]`)
	common.AssertCodeError(t, `''.Split(3)`)
}

func TestString_SplitAt(t *testing.T) {
	common.AssertCode(t, `'1🎁2🎁3'.SplitAt(2)`, `['1🎁', '2🎁3']`)
	common.AssertCode(t, `'hello,world'.SplitAt(6)`, `['hello,', 'world']`)
	common.AssertCode(t, `'hello,world'.SplitAt(100)`, `['hello,world', '']`)
	common.AssertCode(t, `'hello,world'.SplitAt(-2)`, `['hello,wor', 'ld']`)
	common.AssertCode(t, `''.SplitAt(2)`, `['', '']`)
}

func TestString_SplitFn(t *testing.T) {
	common.AssertCode(t, `'1🎁2🎁3'.SplitFn(chr: chr == '🎁')`, `['1', '🎁2', '🎁3']`)
	common.AssertCode(t, `'renato'.SplitFn(chr: chr == 'a')`, `['ren', 'ato']`)
	common.AssertCode(t, `'renato'.SplitFn((_, idx): idx >= 3)`, `['ren', 'a', 't', 'o']`)
}

func TestString_Fields(t *testing.T) {
	common.AssertCode(t, `'1🎁2🎁3'.Fields()`, `['1🎁2🎁3']`)
	common.AssertCode(t, `'1🎁   🎁3'.Fields()`, `['1🎁', '🎁3']`)
	common.AssertCode(t, `'hello, world'.Fields()`, `['hello,', 'world']`)
	common.AssertCode(t, `'hello ping world'.Fields()`, `['hello', 'ping', 'world']`)
	common.AssertCode(t, `'a b c d e f g'.Fields()`, `['a', 'b', 'c', 'd', 'e', 'f', 'g']`)
	common.AssertCode(t, `'.1    1.'.Fields()`, `['.1', '1.']`)
	common.AssertCode(t, `''.Fields()`, `[]`)
}

func TestString_ToLower(t *testing.T) {
	common.AssertCode(t, `'🎁'.ToLower()`, `🎁`)
	common.AssertCode(t, `'Renato'.ToLower()`, `renato`)
	common.AssertCode(t, `'RENATO'.ToLower()`, `renato`)
	common.AssertCode(t, `'ÇÀÃ.'.ToLower()`, `çàã.`)
	common.AssertCode(t, `''.ToLower()`, ``)
}

func TestString_ToUpper(t *testing.T) {
	common.AssertCode(t, `'Renato'.ToUpper()`, `RENATO`)
	common.AssertCode(t, `'renato'.ToUpper()`, `RENATO`)
	common.AssertCode(t, `'çàã.'.ToUpper()`, `ÇÀÃ.`)
	common.AssertCode(t, `''.ToUpper()`, ``)
}

func TestString_ToTitle(t *testing.T) {
	common.AssertCode(t, `'renato pereira'.ToTitle()`, `Renato Pereira`)
	common.AssertCode(t, `'RENATO'.ToTitle()`, `Renato`)
	common.AssertCode(t, `'çàã.'.ToTitle()`, `Çàã.`)
	common.AssertCode(t, `''.ToTitle()`, ``)
}

func TestString_ToSnake(t *testing.T) {
	common.AssertCode(t, `'renato pereira'.ToSnake()`, `renato_pereira`)
	common.AssertCode(t, `'RENATO     Per'.ToSnake()`, `renato_per`)
	common.AssertCode(t, `'çàã.'.ToSnake()`, ``)
	common.AssertCode(t, `''.ToSnake()`, ``)
}

func TestString_ToKebab(t *testing.T) {
	common.AssertCode(t, `'renato pereira'.ToKebab()`, `renato-pereira`)
	common.AssertCode(t, `'RENATO     Per'.ToKebab()`, `renato-per`)
	common.AssertCode(t, `'çàã.'.ToKebab()`, ``)
	common.AssertCode(t, `''.ToKebab()`, ``)
}

func TestString_ToCamel(t *testing.T) {
	common.AssertCode(t, `'renato pereira'.ToCamel()`, `renatoPereira`)
	common.AssertCode(t, `'RENATO     Per'.ToCamel()`, `renatoPer`)
	common.AssertCode(t, `'çàã.'.ToCamel()`, ``)
	common.AssertCode(t, `''.ToCamel()`, ``)
}

func TestString_ToPascal(t *testing.T) {
	common.AssertCode(t, `'renato pereira'.ToPascal()`, `RenatoPereira`)
	common.AssertCode(t, `'RENATO     Per'.ToPascal()`, `RenatoPer`)
	common.AssertCode(t, `'çàã.'.ToPascal()`, ``)
	common.AssertCode(t, `''.ToPascal()`, ``)
}

func TestString_ToDot(t *testing.T) {
	common.AssertCode(t, `'renato pereira'.ToDot()`, `renato.pereira`)
	common.AssertCode(t, `'RENATO     Per'.ToDot()`, `renato.per`)
	common.AssertCode(t, `'çàã.'.ToDot()`, ``)
	common.AssertCode(t, `''.ToDot()`, ``)
}

func TestString_ToTrain(t *testing.T) {
	common.AssertCode(t, `'renato pereira'.ToTrain()`, `Renato-Pereira`)
	common.AssertCode(t, `'RENATO     Per'.ToTrain()`, `Renato-Per`)
	common.AssertCode(t, `'çàã.'.ToTrain()`, ``)
	common.AssertCode(t, `''.ToTrain()`, ``)
}

func TestString_ToSentence(t *testing.T) {
	common.AssertCode(t, `'renato pereira'.ToSentence()`, `Renato pereira`)
	common.AssertCode(t, `'RENATO     Per'.ToSentence()`, `RENATO     Per`)
	common.AssertCode(t, `'çàã.'.ToSentence()`, `Çàã.`)
	common.AssertCode(t, `''.ToSentence()`, ``)
}

func TestString_Ellipsis(t *testing.T) {
	common.AssertCode(t, `'renato pereira'.Ellipsis(5)`, `renat...`)
	common.AssertCode(t, `'RENATO     Per'.Ellipsis(10)`, `RENATO    ...`)
	common.AssertCode(t, `'çàã.'.Ellipsis(2)`, `çà...`)
	common.AssertCode(t, `''.Ellipsis(3)`, ``)
}

func TestString_PadCenter(t *testing.T) {
	common.AssertCode(t, `'renáto'.PadCenter(10)`, `  renáto  `)
	common.AssertCode(t, `'renato'.PadCenter(5)`, `renato`)
}

func TestString_PadCenterWith(t *testing.T) {
	common.AssertCode(t, `'renáto'.PadCenterWith(10, '⚠️')`, `⚠️⚠️renáto⚠️⚠️`)
	common.AssertCode(t, `'renato'.PadCenterWith(5, '⚠️')`, `renato`)
}

func TestString_PadLeft(t *testing.T) {
	common.AssertCode(t, `'renáto'.PadLeft(10)`, `    renáto`)
	common.AssertCode(t, `'renato'.PadLeft(5)`, `renato`)
}

func TestString_PadLeftWith(t *testing.T) {
	common.AssertCode(t, `'renáto'.PadLeftWith(10, '⚠️')`, `⚠️⚠️⚠️⚠️renáto`)
	common.AssertCode(t, `'renato'.PadLeftWith(5, '⚠️')`, `renato`)
}

func TestString_PadRight(t *testing.T) {
	common.AssertCode(t, `'renáto'.PadRight(10)`, `renáto    `)
	common.AssertCode(t, `'renato'.PadRight(5)`, `renato`)
}

func TestString_PadRightWith(t *testing.T) {
	common.AssertCode(t, `'renáto'.PadRightWith(10, '⚠️')`, `renáto⚠️⚠️⚠️⚠️`)
	common.AssertCode(t, `'renato'.PadLeftWith(5, '⚠️')`, `renato`)
}

func TestString_Repeat(t *testing.T) {
	common.AssertCode(t, `'|'.Repeat(10)`, `||||||||||`)
	common.AssertCode(t, `'À'.Repeat(4)`, `ÀÀÀÀ`)
}

func TestString_Reverse(t *testing.T) {
	common.AssertCode(t, `'á🎁à'.Reverse()`, `à🎁á`)
}

func TestString_TrimSpace(t *testing.T) {
	common.AssertCode(t, `''.TrimSpace()`, ``)
	common.AssertCode(t, `'  abc  '.TrimSpace()`, `abc`)
	common.AssertCode(t, `'  🎁  '.TrimSpace()`, `🎁`)
}

func TestString_TrimSpaceLeft(t *testing.T) {
	common.AssertCode(t, `''.TrimSpaceLeft()`, ``)
	common.AssertCode(t, `'  abc  '.TrimSpaceLeft()`, `abc  `)
	common.AssertCode(t, `'  🎁  '.TrimSpaceLeft()`, `🎁  `)
}

func TestString_TrimSpaceRight(t *testing.T) {
	common.AssertCode(t, `''.TrimSpaceRight()`, ``)
	common.AssertCode(t, `'  abc  '.TrimSpaceRight()`, `  abc`)
	common.AssertCode(t, `'  🎁  '.TrimSpaceRight()`, `  🎁`)
}

func TestString_Trim(t *testing.T) {
	common.AssertCode(t, `''.Trim(' |[]')`, ``)
	common.AssertCode(t, `' |abc]  '.Trim(' |[]')`, `abc`)
	common.AssertCode(t, `'  [🎁]  '.Trim(' |[]')`, `🎁`)
}

func TestString_TrimLeft(t *testing.T) {
	common.AssertCode(t, `''.TrimLeft(' |[]')`, ``)
	common.AssertCode(t, `' |abc]  '.TrimLeft(' |[]')`, `abc]  `)
	common.AssertCode(t, `'  [🎁]  '.TrimLeft(' |[]')`, `🎁]  `)
}

func TestString_TrimRight(t *testing.T) {
	common.AssertCode(t, `''.TrimRight(' |[]')`, ``)
	common.AssertCode(t, `' |abc]  '.TrimRight(' |[]')`, ` |abc`)
	common.AssertCode(t, `'  [🎁]  '.TrimRight(' |[]')`, `  [🎁`)
}

func TestString_Replace(t *testing.T) {
	common.AssertCode(t, `''.Replace('a', 'no')`, ``)
	common.AssertCode(t, `'a b c a'.Replace('a', 'no')`, `no b c no`)
	common.AssertCode(t, `'🎁,🎁!ç🎁'.Replace('🎁', '🥸')`, `🥸,🥸!ç🥸`)
}

func TestString_ReplaceN(t *testing.T) {
	common.AssertCode(t, `''.ReplaceN('a', 'no', 1)`, ``)
	common.AssertCode(t, `'a b c a'.ReplaceN('a', 'no', 1)`, `no b c a`)
	common.AssertCode(t, `'🎁,🎁!ç🎁'.ReplaceN('🎁', '🥸', 1)`, `🥸,🎁!ç🎁`)
}

func TestString_Sort(t *testing.T) {
	common.AssertCode(t, `''.Sort()`, ``)
	common.AssertCode(t, `'a c b'.Sort()`, `  abc`)
	common.AssertCode(t, `'z23aç🎁'.Sort()`, `23azç🎁`)
}

func TestString_SortFn(t *testing.T) {
	common.AssertCode(t, `''.SortFn((a, b): { a <=> b })`, ``)
	common.AssertCode(t, `'a c b'.SortFn((a, b): { a <=> b })`, `  abc`)
	common.AssertCode(t, `'z23aç🎁'.SortFn((a, b): { a <=> b })`, `23azç🎁`)
}

func TestString_Contains(t *testing.T) {
	common.AssertCode(t, `''.Contains('a', 'z')`, `false`)
	common.AssertCode(t, `'a c b'.Contains('a', 'z')`, `true`)
	common.AssertCode(t, `'zoo'.Contains('a', 'z')`, `true`)
	common.AssertCode(t, `'bool'.Contains('a', 'z')`, `false`)
	common.AssertCode(t, `'bool'.Contains('🎁')`, `false`)
	common.AssertCode(t, `'bool🎁'.Contains('🎁')`, `true`)
	common.AssertCode(t, `'bool🎁'.Contains()`, `false`)
}

func TestString_ContainsChars(t *testing.T) {
	common.AssertCode(t, `''.ContainsChars('an')`, `false`)
	common.AssertCode(t, `'a c b'.ContainsChars('an')`, `true`)
	common.AssertCode(t, `'a c banana'.ContainsChars('an')`, `true`)
	common.AssertCode(t, `'z23aç🎁'.ContainsChars('🎁🎁')`, `true`)
	common.AssertCode(t, `'z23aç🎁🎁asdfasd'.ContainsChars('🎁🎁')`, `true`)
}

func TestString_StartsWith(t *testing.T) {
	common.AssertCode(t, `''.StartsWith('')`, `true`)
	common.AssertCode(t, `'abc'.StartsWith('')`, `true`)
	common.AssertCode(t, `''.StartsWith('a')`, `false`)
	common.AssertCode(t, `'a c b'.StartsWith('a', 'z')`, `true`)
	common.AssertCode(t, `'za c b'.StartsWith('a', 'zas')`, `false`)
	common.AssertCode(t, `'zas c b'.StartsWith('a', 'z')`, `true`)
	common.AssertCode(t, `'zas c b'.StartsWith('a', 'zas')`, `true`)
	common.AssertCode(t, `'c b'.StartsWith('a', 'z')`, `false`)
	common.AssertCode(t, `'🎁💡🔫🎁😳'.StartsWith('🎁🎁')`, `false`)
	common.AssertCode(t, `'🎁🎁💡🔫🎁😳'.StartsWith('🎁🎁')`, `true`)
	common.AssertCode(t, `'🎁🎁🎁💡🔫🎁😳'.StartsWith('🎁🎁')`, `true`)
}

func TestString_EndsWith(t *testing.T) {
	common.AssertCode(t, `''.EndsWith('')`, `true`)
	common.AssertCode(t, `'abc'.EndsWith('')`, `true`)
	common.AssertCode(t, `''.EndsWith('a')`, `false`)
	common.AssertCode(t, `'a c b'.EndsWith('b', 'z')`, `true`)
	common.AssertCode(t, `'za c bal'.EndsWith('a', 'ball')`, `false`)
	common.AssertCode(t, `'zas c ball'.EndsWith('a', 'ball')`, `true`)
	common.AssertCode(t, `'c b'.EndsWith('a', 'z')`, `false`)
	common.AssertCode(t, `'🎁💡🔫🎁😳'.EndsWith('🎁😳')`, `true`)
	common.AssertCode(t, `'🎁🎁💡🔫🎁😳😳'.EndsWith('🎁😳')`, `false`)
	common.AssertCode(t, `'🎁🎁🎁💡🔫🎁😳🎁'.EndsWith('🎁😳')`, `false`)
}

func TestString_IsEmpty(t *testing.T) {
	common.AssertCode(t, `''.IsEmpty()`, `true`)
	common.AssertCode(t, `'zas c ball'.IsEmpty()`, `false`)
	common.AssertCode(t, `'💡🔫🎁😳'.IsEmpty()`, `false`)
}

func TestString_Chars(t *testing.T) {
	common.AssertCode(t, `
	str := 'renato'

	result := ''
	calls := 0
	for c in str.Chars() {
		calls += 1
		result += c
	}
	
	[calls, result]
	`, `[6, 'renato']`)

	common.AssertCode(t, `
	str := 'renato'

	result := ''
	calls := 0
	for c in str {
		calls += 1
		result += c
	}
	
	[calls, result]
	`, `[6, 'renato']`)
}

func TestString_Lines(t *testing.T) {
	common.AssertCode(t, `
	str := '\nthis\nis\nlines\n'

	result := ''
	calls := 0
	for c in str.Lines() {
		calls += 1
		result += c
	}
	
	[calls, result]
	`, `[5, 'thisislines']`)
}

func TestString_Words(t *testing.T) {
	common.AssertCode(t, `
	str := 'this, is       a word!'

	result := ''
	calls := 0
	for c in str.Words() {
		calls += 1
		result += c
	}
	
	[calls, result]
	`, `[4, 'this,isaword!']`)
}
