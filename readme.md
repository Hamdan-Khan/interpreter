# Interpreter

A tree-walk interpreter for lox programming language written in Go.

It uses the grammar defined in crafting interpreters book by Rob Nystrom. You can find the grammar [here](https://craftinginterpreters.com/appendix-a.html). 

Currently, it doesn't support classes but I plan to add them in the future. This is more of a learning project for me than an attempt to create a perfectly working interpreter. I have documented the learning and the process of working of a tree-walk interpreter in this [blog](https://hamdan-khan.github.io/--).

## Syntax

### Functions
```lox
fun sum(a, b) {
    return a + b;
}
```
### Print
```lox
print "Hello world";
```
### Variables
```lox
var a = 1;
```
### If
```lox
if (true) {
    print "true";
} else {
    print "false";
}
```
### While
```lox
while (condition) {
    print "loop";
}
```
### For
```lox
for (var i = 0; i < 10; i = i + 1) {
    print i;
}
```


## Usage

```bash
go run . test.txt
```
