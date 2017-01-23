# Dusk-lang
Dusk Language interpreter

Dusk is a small project I've been working on in my spare time. For the purposes of learning to design/lex/parse/evaluate a programming language that's a little more than a 'toy language'

### So far Dusk features are:

- Dynamic/Interpreted
- Simple syntax
- Optional semicolons for terminating statements
- Most statements are expressions
- 64 bit Integers
- 64 bit Floats
- Strings
- Arrays
- First class functions
- Closures
- Classes by closures and '.' operator access

### Planned features:
- Maps
- For - in loops
- Interpolated formatting of strings e.g `"hello, \{person.name}"`
- Pairs
- Modules
- Bytecode compilation and evaluatation

# Examples
## Hello World

```
let name = readln! 

println(name)
```
## Some examples
### let statements
```
// let statements are the basis of creating variables

let name = "friendo"             // bind a string literal
let hobby = 'unit testing'      // bind another string using single quote
let age = 20                   // bind a int64 literal
let height = 187.3            // bind a float64 literal

let array = [name, hobby, age, height] // put them into an array

```

### operations on strings and arrays
```
let a = 'buddy'
a[0]  // 'b'
a[-1] // 'y'  arrays/strings wrap negatively around back to 0

let chant = [2,4,6,8]
chant = push(chant, "who do we appreciate?") // appends the string to the back of chant and returns the new array

// alternately you can use the '+' or  '+=' to concat arrays
chant += ["infix operators!"]
```

### Builtin functions
```
// basic array functions
let a = [1,2,3,4]
len(a)     // 4
first(a)   // 1
last(a)    // 4
rest(a)    // 2,3,4
lead(a)    // 1,2,3
push(a, 5) // 1,2,3,4,5
alloc(256, 'a') // creates an array of 256 a's.. can be any value
set(a, 0, 6) // a[0] = 6,2,3,4

// basic string functions
let s = "hello, friend"
split(s, '')     // splits s into an array of it's characters ['h', 'e', 'l', 'l' ... ]
split(s, ', ')  // splits by ', '. ['hello', 'world']
join(a, '')    // joins an array into a string of it's objects
join(a, '.')  // joins with a '.' in between each element


// i/o functions
println
print
readln
read
readc
readall
```

### functions
```
// functions are literals aswell
// functions are defined with the '|' arg1, arg2 '|' syntax
// return statement is 'ret'

// if a function only has one statement or expression it can directly follow the definition on the same line
let birthday = || age += 1      // function with no args. increments age by 1

// if there is more than one line use braces starting on the same line
let shrink = |n| {
  ret height - n      // return statement
}

let grow = |n| {
  height + n          // return statement is optional if the last statement if the result to be returned
}

// usages

// when a function call has no args you can optionally use '!' syntax instead of '()'
birthday!   // optionally birthday()

let newHeight = grow(45)
newHeight = shrink(age) // assign newHeight to the result of shrink
```
### power operator!!
```
// power operator
4^3 // returns 4 to the power of 3
```
### If statements
```
// compact syntax
if true: "yes" else "no"

// full syntax
if true {
  ret "yes"
} else {
  ret "no"
}

// syntax can be mixed and matched
if true {
  if a == b: "yes" else {
    ret "no"
  }
} else "no"

// if else
if a == 1 {
  "one"
} else if a == 2 {
  "two"
} else "huge!"
```
### Closures
```
// define a function that takes one arg and returns a function that sums it's argument together

let newAdder = |a| {
  ret |b| {
    ret a + b
  }
}
// compact syntax of same definition let newAdder = |a| |b| a + b

let add2 = newAdder(2)

add2(4) // 6
```

### Classes
```
let person = |n| {                          // person acts like a construtor returning a 'new' person
  let name = n                              // local variable 'name'
  let sayhi = || println("hello, " + name)  // local function 'sayhi'
  ret || person                             // return a closure of self. This is the new class
}

let p1 = person("ted")
let p2 = person("bob")

p1.name              // prints "ted"
p2.name = "kyle"    // change p2's name to "kyle"
p2.sayhi!          // prints "hello, kyle"

```

## Building source
Place contents in `$GOPATH/src/jacob/dusk`
`go build`
`./dust` to start repl
`./dust file.dusk` to just run file


#### reference:
Ball, Thorsten. “Writing An Interpreter In Go.” 2016
My main guide for designing the internal interpreters structure. Used his tests for checking correctness.
It's a great book that I suggest if you want to write a programming language with more advanced features then most online literture. However always right the code yourself and change it so it fits with your idea of the program, adding and removing elements
