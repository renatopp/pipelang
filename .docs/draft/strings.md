# Streams

Chars(this)
Lines()
Words()

# Access

Get(this, n) string (char)
Cut(this, from, to) string
CutTo(this, n) string
SubFrom(this, n) string
SubFirst(this, n) string
SubLast(this, n) string

# Formatters

Ellipsis(this, size) String
PadCenter(this, size) String
PadLeft(this, size) String
PadRight(this, size) String
Repeat(this, n) string
Reverse(this) string
Trim(this, sub) string
TrimStart(this, sub) string
TrimEnd(this, sub) string
ToAscii()

## Casing

ToLower(this) string
ToUpper(this) string
ToTitle(this) string
ToSnake(this) string
ToKebab(this) string
ToDot(this) string
ToCamel(this) string
ToParagraph(this) string

## Partials

Or(this, default) string
Join(this, array) string

# Checkers

Contains(this, search) Boolean
ContainsAny(this, ...search) Boolean
EndsWith(this, suffix) bool
EndsWithAny(this, ...suffix) bool
StartsWith(this, prefix) bool
StartsWithAny(this, ...prefix) bool
IsEmpty(this) bool
Size() int
