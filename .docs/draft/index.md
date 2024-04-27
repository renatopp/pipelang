# Drafting the language

## Done

```

<expresion> = <prefix><value>
            = <value><op><value>

```

- Basics
  - Modules
  - Imports
  - Types, Checks and Conversions
  - Assignments and Operators
- Functions
  - High order
  - Lambda
  - Arguments
  - Pipes
  - Comments
- Control Flow
  - Cases
  - Ifs
  - Loops

---

Core:

- Fast to develop
- Discardable
- Intuitive
- Batteries included

---

## The Basics

The language will have:

**Comments**

```
-- This is a comment
```

**Types**

```
number   := 10.2
sting    := 'Hello, World'
booleans := false
lists    := List { 1, 2, 3 }
dicts    := Dict { a=1, b=2 }
maybe    := Maybe()
```

**Functions**

Creating functions:

```
function := fn() {
  return 'Hello'
}

fn function() {} -- translated to function := fn () {}
```

Function arguments:

```
fn {}
fn () {}
fn (a) {}
fn (a, _) {}
fn (a, ...b) {}
fn (a, ...b, c) {}
```

Lambda functions:

```
: 1         -- fn() { 1 }
x:x*2       -- fn(x) { x*2 }
(x, y): x+y -- fn(x, y) { x + y }
x: {}       -- fn(x) {}
```

**Pipes**

```
value
| function x, y, z
| other

-- equivalent to other(function(value, x, y, z))
```

**Operations and Assignments**

Assignments:

```
-- go like
a := 2
a = 4

-- tuple
(a, b) = (2, 3)
(a, ...b) = (2, 3, 4, 5)
(a, b) = list...
(a, ...b, c) = list...

-- (a, b) = (2, 3) is equivalent to:
-- a = 2
-- b = 2

-- (a, ...b) = (1, 2, 3) is equivalent to:
-- a = 1
-- b = List{ 2, 3 }

-- discard

(a, _) = (2, 3)

```

Operations:

```
-- arithmetic
1 + 1
1 - 1
1 * 1
1 / 1
1 % 1
1 ^ 1
1 // 1

-- equality
1 == 1
1 != 1
1 < 1
1 <= 1
1 > 1
1 >= 1

-- logic
not 1
1 and 1
1 or 1
1 xor 1
```

Control flow

```
if expression_list { ... }
elif expression_list { ... }
else { ... }

case {
  (1, 1): {}
  (1, _): {}
  (_, 1): {}
  (..._): {}
}

for expression {}
for {}
```

Other expression

```
a as value -- value := a
a in value -- a = Stream(value).Next()
a is value -- returns bool contextually
```

**Modules**

Exports only uppercase variables, same as go

```
math.Sqrt -- builtin auto imported
import(x) -- carrega lib de um path (pode ser uma url)

package := import(x)
import(x) as a
```

**Custom Types**

```
type A String {
  fn CustomFunction() {
    println(1)
  }
}

-- same for any type
```

Data Types

```
type Node data {
  value = 2
  other = 3

  fn Sample() {}
}

Node { sample }
```

## Expression Templates

`[claim]` is

- `if` `[expression]`
- `if` `[expression] [block]`
- `for` `[expression] [block]`
- `return` `[expression]`
- `raise` `[expression]`
- `yield` `[expression]`
- `break`
- `continue`
- `[expression]`
- `[claim]` `[eol]` `[claim]`
- empty

and

`[expression]` is

- `[identifier]`
- `[literal]`
- `[tuple]`
- `case` `[expression] [block]`
- `type [identifier] [expression]`
- `type [identifier] [expression] [block]`
- `[expression]` `[tuple]`
- `[expression]` `[block]`
- `[pre operator]` `[expression]`
- `[expression]` `[post operator]`
- `[expression]` `[operator]` `[expression]`
- `[identifier tuple]` `[assignment]` `[expression]`

and

`[block]` is

- `{` `[claim]` `}`

and

`[identifier tuple]` is

- `(` `)`
- `(` `[identifier]` [`,` `[identifier]`]\* `,`? `)`
- `[identifier]`
- `...[identifier]`

and

`[tuple]` is

- `(` `)`
- `(` `[expression]` [`,` `[expression]`]\* `,`? `)`
