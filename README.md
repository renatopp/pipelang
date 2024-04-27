# PIPE LANG

PIPE Lang is an **experimental** language designed to handle streams of data in a compact and readable style. Focused on daily automation tasks, the goal is simplify the scripts that requires sequential (and eventually some parallel) steps.

## Why should I use it?

You shouldn't. Feel free to try it, but at this point, the language is very unstable and more experimental than practical.

## What does it look like?

```haskell
-- Euler problem 2: Find the sum of the even-valued terms in the Fibonacci sequence whose values do not exceed four million
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
```

You may try it out in [https://pipe.r2p.dev](https://pipe.r2p.dev).

## The CLI

```bash
# Opens the REPL tool
pipe shell

# Evaluate inline
pipe eval 'print("hi")'

# Run file
pipe run file.pp
```

## Features

### The Type System

PIPE is a strong dynamic typed language with minimal builtin types:

```haskell
-- They are all `Number` type
num1 := 100
num2 := 100.0
num3 := 1e10

-- They are all `Boolean` type
bool1 := true
bool2 := false

-- They are all `String` type
str1 := 'Single-quoted strings'
str2 := "Single-quoted strings"
str3 := `Raw Strings
Can be Multiline`

-- `List` type
list1 := List{}
list2 := [1, 2, 3, 4, 5]
list3 := ['a', 1, [3, true]]

-- `Dict` type
dict1 := Dict {}
dict2 := { a=1, b=2, 3=4 }

-- `Maybe`
maybe := Maybe(2)
maybe.Ok()
maybe.Error()
maybe.Value()
maybe.Result()

-- `Stream` as the iterator object
list2.Elements() -- = <Stream>
```

The language uses Go short declaration style `:=`. Declared variables may be reassigned but only if the types match. Variables can be re-declared any time.

```haskell
var := 'Im a string'
var = 2               -- will raise an error!
var := 2              -- redeclare var as a number, ok
```

Tuples are used in assignments in an overly complex dynamics:

```haskell
a, b    := 1, 2      -- a=1; b=2
a       := 1, 2      -- a=1; ignoring 2
a, b    := 1         -- error!
a, ...b := 1         -- a=1; b=[]
a, ...b := 1, 2, 3   -- a=1; b=[2, 3]
a, b    := [1, 2]... -- a=1; b=2
```

Type conversion can be done explicitly:

```haskell
String(3)
Number('3')
Bool(3)
Stream([1, 2])
```


### Operations

Notice that PIPE is strong typed, so operations can only be performed when involves two variables of the same type, with exception to equality, logical and concat operators

```haskell
-- Arithmetic
a + b
a - b
a / b
a * b
a % b -- mod
a ^ b -- pow

-- Relational
a < b
a > b
a <= b
a >= b
a == b
a != b
a <=> b -- spaceship operator, returning -1, 0, 1 

-- Logical
!a
a and b
a or b
a xor b

-- Other
a .. b -- concat as string
```

### Functions

Pretty much everything in the language is an expression, which returns something. An empty block `{}`, such as in ifs, for and functions, generates false as default, otherwise it will return the last executed expression.

```haskell
-- `return` keyword is optional. 
fn add(a, b) {
  a + b
}

-- Function definition is just an expression
add := fn(a, b) {
  a + b
}

-- Functions are high order and respect closures
fn mult(x) {
  fn (a, b) { (a + b)*x }
}
mult(2)(1, 5) -- = 12

-- Lambda functions are definied by `:`
lambda1 := a : a*2
lambda2 := (a, b): a + b

-- Shortcuts
fn (a, b) {} -- Nameless
fn ping {}   -- Parameterless
fn {}        -- Nameless parameterless 

-- Generator functions return a `Stream` object, that is only
-- evaluated when queried.
fn OneTwoThree {
  yield 1
  yield 2
  yield 3
}
stream := OneTwoThree()
stream.Next() -- Maybe(1)
stream.Next() -- Maybe(2)
stream.Next() -- Maybe(3)
stream.Next() -- Maybe(Error)
```

### Function Chaining (AKA pipe expressions)

Function chaining (or pipe expressions) are the core of the language, it uses the power of generator functions to create processors that evaluate streams of data sequentially in a lazy way.

```haskell
[1, 2, 3, 4]
| filter x: x.IsEven()
| map    x: x*2
| sum
-- = 12
```

The code above is equivalent to `sum( map( filter( [1,2,3,4] , x:x.IsEven()) , x:x*2) )`.

Note that pipes are just a compact way to call functions. This means that `a | b` is equivalent to `b(a)`. The pipe expression has an special notation, where the first identifier after the pipe is the function to be called, and the remaining of the elements in the line are the arguments (without call parenthesis to reduce clutter):

- `a | b 1, 2, 3` is equivalent `b(a, 1, 2, 3)`
- `a | b 2, (a, b): a+b` is equivalent `b(a, 2, <lambda>)`

Most of builtin functions that operate in pipes converts the first argument to a stream, forcing the generator, thus, forcing it to be lazy.

### Flow Controls

Flow controls are a mixture of go, python and rust:

```haskell
-- Simple conditional
if a == b { ... }

-- Conditional with additional expressions
-- Notice that only the last expression is considered as conditional
if a := b(); a { ... }

-- Pattern matching
match x, y {
	     0, 0: println('zero zero')
	_ as a, 0: println(a, 'zero')
			 _, _: println('im done')
}

-- Infinite loop
for { ... }

-- Simple for
for a == b { ... }

-- For in stream
for a in stream { ... }
```

### Error Handling

Any function can throw errors by using the `raise` keyword:

```haskell
fn explode {
	raise 'error'
}
```

In order to capture this error, PIPE uses a wrap operator `?`:

```haskell
result := explode()?

if result.Error() {
	-- treat error
}

val := result.Value()
-- use the non-error value here
```

### Custom Data Types

You may define a custom data type, which are custom structures which you may mix with other structures. Note that this is more like a mixin pattern than inheritance.

```haskell
data BaseNode {
	id     = ''
	parent = empty(BaseNode) -- returns a Maybe<BaseNode>(empty_error)

	-- methods MUST have a `this` variable as first argument
	fn String(this) {
		return id
	}
}

data NumberNode(BaseNode) {
	value = 0

	fn String(this) {
		return sprintf('<%s, %d>', BaseNode.String(this), this.value)
	}
}

number := NumberNode { id='xyz', value=5 }
number.String() -- '<xyz, 5>'
```

