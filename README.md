# Dusk-lang
Dusk Language interpreter

Dusk is a small project I've been working on in my spare time. For the purposes of learning to design/lex/parse/evaluate a programming that's a little more than a 'toy language'

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
- Pairs
- Maps
- Modules
- Bytecode compilation and evaluatation

#Examples
##Hello World

```
let name = readln! 

println(name)
```

## Some examples
###let statements
```
// let statements are the basis of creating variables

let name = "friendo"    // bind a string literal
let age = 20           // bind a int64 literal
let height = 187.3    // bind a float64 literal

```
### functions
```
// functions are literals aswell
// functions are defined with the '|' arg1, arg2 '|' syntax

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
###Closures
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

###Classes
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


