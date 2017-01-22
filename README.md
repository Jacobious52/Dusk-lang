# Dusk-lang
Dusk Language interpreter

Dusk is a small project I've been working on in my spare time. For the purposes of learning to design/lex/parse/evaluate a programming that's a little more than a 'toy language'

### So far Dusk features are:

- Dynamic/Interpreted
- Simple syntax
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
##Hellow world

```
// readln is a builtin function.
// when a function call has no args you can optionally use '!' syntax instead of '()'
let name = readln! 

print(name)
```
