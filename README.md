
# RuneVM: A Runey Scripting Language
![](editor/vscode/rune/logo.png)

RuneVM is an interpreter for the versatile scripting language "Rune".

Rune is a dynamically typed scripting language designed to be easy to use and implement while remaining practical for real-world tasks. While performance isn't its primary focus, Rune excels in simplicity and flexibility.

Rune stands out for being:
- **Intuitive:** Designed to be user-friendly and straightforward.
- **Embeddable:** Seamlessly integrates into your Go projects.
- **Extensible:** Easily extendable with custom functions and variables.
- **Portable:** Simple to port to other languages.
- **Versatile:** Can be used standalone as a script runner with the [Rune](dist/) binary or embedded as a scripting language within your project.

## Everything you would definetely need for scripting, but not more and not less:
- **Datatypes:** `int`, `float`, `string`, `bool`
- **Branching:** `if`, `elif`, `else` and optional `then`
- **Looping:** `while`
- **Operators:** `=`, `||`, `&&`, `<`, `>`, `<=`, `>=`, `==`, `!=`, `+`, `-`, `*`, `/`, `%`
- **Arrays**: `array`
- **Dictionaries/Maps**: `table`
- **Compound types and field access**: e.g. `person1 = new(person)` and  `person1.name = "John"` for more info see: [Field access](#field-access)
- **Import:** See: [Code modularization](#code-modularization)
- **Builtin functions**: Look at the [builtin functions](#builin-functions) for more info
- **Editor plugins**: See: [editor plugins](#editor-plugins)
- **Simple interop between Rune and Go**: See [interop](#interop-between-rune-and-go)


## Pre compiled binaries
In the [dist directory](dist/), you will find pre compiled binaries for the three major platforms.

## Build from source
Simply run the `build.bat` on Windows or the `build.sh` on Linux or Mac from the [dist directory](dist/). This will build the Rune binary for all three platforms in there respective directories.

## Using the Language

### Writing Scripts

Rune scripts have the `.rune` extension. Here is a sample script (`test.rune`):

```js
myvar = 12;

printer = fun(toPrint) {
    println(toPrint)
};

if myvar == 10 then
    printer("then")
elif myvar == 12 then
    printer("elif 12\n")
elif myvar == 11 then
    printer("elif 11")
else printer("else");

while myvar > 0 do {
    myvar = myvar  -1;
    printer(myvar)
};

if myvar < 0 {
    printer("Is smaller than 0")
} elif myvar > 0 {
    printer("Is biger or than 0")
} else {
    printer("Is exactly 0")
};

if true then printer("Nice") else printer("Also nice");

person1 = table{}
person1.name = "John"

person1.sayHello = fun(self) {
    println(self.name)
}
person1.sayHello() # prints: "John"

person2 = new(person1)
person2.name = "Jenny"

person1.sayHello() # prints: "John"
person2.sayHello() # prints: "Jenny"

mytable = table{"key 1": false, "key 2": 10, "key 3": 12.4} # spaces in keys will be removed automatically
println(mytable) # prints: {"key1": false, "key2": 10, "key3": 12.4}
println(mytable["key2"]) # prints: 10
mytable["key2"] = 11
println(mytable["key2"]) # prints: 11
mytable.key2 = 12
println(mytable["key2"]) # prints: 12
```

### Running Scripts

To run a script, use the following command:

```
rune path/to/your/script.rune
```

## Code modularization

Rune supports an `import` statement to include and execute other Rune scripts within the current script. This allows for better modularization and reuse of code. The `import` statement takes a file path (without the `.rune` extension) and imports the contents of the specified file into the current script.

### Syntax

```js
import "path/to/script"
```

- `import` is an expression which should evaluate to a string that specifies the path to the script file.
- The `.rune` extension is automatically added by the parser.

### Example
Assume we have two Rune scripts:
1. `main.rune`
```js
import "other"

println("This is the main script.")
```
2. `other.rune`
```js
println("This is the imported script.")
```

When you run `main.rune`, the output will be:
```
This is the imported script.
This is the main script.
```

## Conditional Imports
Since `import` is an expression in Rune, it can be used conditionally within the script:
```js
condition = 1

if condition == 2 then
    import "test2"
else
    import "test3"

println("This is the main script.")
```

In this example, depending on the value of `condition`, either `test2.rune` or `test3.rune` will be imported and executed.

Or this works also:
```js
toImport = "test2"

if true then toImport = "test3"

import toImport

println("This is the main script.")
```
In this example, `test3.rune` will be imported and executed.

### Notes
- All imports share the same global scope, meaning variables and functions defined in the imported script are accessible in the main script and vice versa.
- The imported script is executed immediately at the point of the import statement, and any side effects (such as variable assignments or function definitions) will affect the global environment.
- Import paths are relative to the location of the script file executing the import.
- It is idiomatic in Rune to have a single main.rune script that imports all necessary files, rather than scattering import statements throughout various Rune scripts. This approach ensures a clear and organized entry point for the program.
- Cyclic imports will produce an error

## Embedding RuneVM in Your Project

To embed RuneVM in your own Go project, follow these steps:

1. Add `runevm` to your project:
    ```
    go get github.com/RednibCoding/runevm
    ```

2. Import RuneVM in your Go code:
    ```go
    import "github.com/RednibCoding/runevm"
    ```

### Example Usage

Here is an example of how to use RuneVM in your Go project:

```go
package main

import (
	"fmt"
	"os"

	"github.com/RednibCoding/runevm"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Println("USAGE: rune <sourcefile>")
		os.Exit(1)
	}
	source, err := os.ReadFile(args[1])
	if err != nil {
		fmt.Printf("ERROR: Can't find source file '%s'.\n", args[1])
		os.Exit(1)
	}

	filepath := args[1]
	vm := Runevm.NewRuneVM()
	vm.Run(string(source), filepath)
}
```

## Interop between Rune and Go

### Functions

Functions can be defined and added to the VM using the `Set` method. Functions must have the following signature:

```go
func(args ...interface{}) interface{}
```

Example:

```go
func customPrintFunction(args ...interface{}) interface{} {
    for _, arg := range args {
        switch v := arg.(type) {
        case string:
            fmt.Print(v)
        case int:
            fmt.Print(v)
        case float64:
            fmt.Print(v)
        default:
            return fmt.Errorf("unsupported argument type")
        }
    }
    fmt.Println()
    return nil
}

vm.Set("print", customPrintFunction)

```
>**Note:** It is important to always check the number of arguments and their types, as you don't know what errors users might make in their scripts:
```go

func customPrintFunction(args ...interface{}) interface{} {
    // Check the number of arguments (assuming 2 are expected here)
	if len(args) != 2 {
		return fmt.Errorf("move requires exactly 2 arguments")
	}

	// Using type assertions to check if x and y are of type int
	x, ok1 := args[0].(int)
	y, ok2 := args[1].(int)

	if !ok1 || !ok2 {
		return fmt.Errorf("both arguments must be of type int")
	}

	fmt.Printf("x: %d, y: %d\n", x, y)

	Mover.Move(x, y)
	return nil
}
```


### Variables

Custom variables can be defined and added to the VM using the `SetXXX` methods: `SetInt`, `SetFloat`, `SetString`, `SetBool` and `SetArray`.

Variables can be of type `int`, `float64`, `string`, `bool`, `array` or `table`.

Example:

```go
vm.SetString("greetings", "Hello from VM!")
vm.SetInt("myInt", 10)
vm.SetArray("myArr", []string{"One", "Two", "Three"})
vm.SetArray("myArr2", []interface{}{"One", 10, false})
```

## Error Handling

If a custom function returns an error, the interpreter will report it with the line and column number where the error occurred.

Example:

```go
func customPanicPrinter(args []interface{}) error {
    for _, arg := range args {
        if arg == "panic" {
            return fmt.Errorf("intentional panic triggered")
        }
        fmt.Print(arg)
    }
    return nil
}
```

If the `customPanicPrinter` function encounters the string `"panic"`, it will return an error, and the interpreter will handle and report it.

```js
customPanicPrinter("panic")
```
output:
```
error (example.Rune:1:6): Error in function call: 'intentional panic triggered'
```

## Using Functions and Variables defined in Rune from Go
You can get function defined in `Rune` via the `GetFun` function:

Let's say in Rune you have the following function named "printer"
```js
printer = fun(printme) {
    println(printme)
}
```

You can get this function by first, running the script and then call `GetFun` afterwards:
```go
vm := Runevm.NewRuneVM()
vm.Run(string(source), filepath)

printerFunc, err := vm.GetFun("printer")
if err != nil {
    fmt.Println(err)
    return
}

// call the function
printerFunc("Hello From PrinterFunc")
```

Output:
```
Hello From PrinterFunc
```

Similarely you can retrieve the value of variables.

Lets say you have a string variable defined in Rune like so:

```js
toPrint = "I am the toPrint variable"
```

Then you can get the value of it by calling `GetString`:

```go
toPrint, err := vm.GetString("toPrint")
if err != nil {
    fmt.Println(err)
}

fmt.Println(toPrint)
```

output:
```
I am the toPrint variable
```

You can also execute a function defined in Rune with a variable defined in Rune as argument.

Lets combine the above:

```js
toPrint = "I am the toPrint variable"

printer = fun(printme) {
    println(printme)
}
```

And in go we can do:
```go
vm := Runevm.NewRuneVM()
vm.Run(string(source), filepath)

// Get the printer function from Rune
printerFunc, err := vm.GetFun("printer")
if err != nil {
    fmt.Println(err)
    return
}

// Get the 'toPrint' variable from Rune
toPrint, err := vm.GetString("toPrint")
if err != nil {
    fmt.Println(err)
}

// Use both
printerFunc(toPrint)
```

output:
```
I am the toPrint variable
```

Analogous to `GetString` are the functions: `GetInt`, `GetFloat`, `GetBool`, `GetArray` and `GetTable`.

### Getting Tables defined in Rune
Image you have the following Rune code:
```js
person = table{}
person.name = "John"
person.sayHello = fun(self) {
    println("Hello ", self.name)
}
```

You can get the `person` table in go by using the `GetTable` function:
```go
vm := Runevm.NewRuneVM()
vm.Run(string(source), filepath)

person, err := vm.GetTable("person")
if err != nil {
    panic(err.Error())
}
```
To get the function `sayHello` from person, you can use the `GetTableFun` function:
```go
person, sayHello, err := vm.GetTableFun("person", "sayHello")
if err != nil {
    panic(err.Error())
}
```
Now you have the reference to the sayHello function defined on the person object.

In order to call it, you have to pass the person object to the function (this acts as the self parameter). In Rune, this happens automatically, but since we call the function manually in go, we have to pass the self argument manually (otherwise you would get a Rune error: `Variable self is not an array or map`):
```js
sayHello(person)
```


# Rune Language Specification

The Rune language is a simple, dynamic scripting language. The following chapter describes the syntax and features of the Rune language, including how to define and call functions, use variables, control flow with `if` and `while` statements, data types, arrays, tables and more.

## Defining and Calling Functions

### Defining Functions

Functions in Rune are defined by assigning a function to a name.

Example:

```js
greet = func() {
    print("Greetings")
}
```

### Calling Functions

Functions are called using the name followed by parentheses and optional arguments.

Example:

```js
greet()
```

## Defining and Using Variables

### Defining Variables

Variables in Rune are defined by assigning a value to a name using the `=` operator. Variables are dynamically typed and can hold numbers, strings or booleans.

Example:

```js
message = "Hello, World!";
count = 42;
pi = 3.14;
isTrue = true;
```

### Using Variables

Variables can be used in expressions and statements.

Example:

```js
println(message)
print("The count is: ", count)
```

## Operators
Supported operators are: `=`, `||`, `&&`, `<`, `>`, `<=`, `>=`, `==`, `!=`, `+`, `-`, `*`, `/`, `%`

They do what you would expect. I guess no further explaination is needed.

## Control Flow

### If Statements

The `if` statement is used to execute a block of code conditionally. The `else` and `elif` keywords can be used for additional conditions.

Example:

```js
if count > 10 {
    print "Count is greater than 10\n"
} elif count == 10 {
    print "Count is 10\n"
} else {
    print "Count is less than 10\n"
}
```

### While Statements

The `while` statement is used to execute a block of code repeatedly as long as a condition is true.

Example:

```js
x = 0
while x < 5 {
    print "x is ", x, "\n"
    x = x + 1
}
```

## Data Types

### Number

Numbers in Rune can be of type `int` or `float64`. They are used in mathematical expressions and comparisons.

Example:

```js
a = 10
b = 3.14
result = a * b
```

### String

Strings in Rune are sequences of characters enclosed in double quotes `"`.

Example:

```js
greeting = "Hello, World!";
print(greeting);
```

### Bool

Booleans in Rune can have one of two values: `true` and `false`.

```js
RuneIsAwesome = true

if RuneIsAwesome then println("Rune is awesome!") else println("This should never print! ;)")
```

## Falsy Values

In Rune, the following values are considered falsy:
- The number `0`
- The empty string `""`
- The boolean value `false`

Any other value is considered truthy.

### Example of Falsy Values

```js
if 0 {
    println("This will not print")
} else {
    println("0 is falsy")
}

if "" {
    println("This will not print")
} else {
    println("Empty string is falsy")
}

if false {
    println("This will not print")
} else {
    println("false is falsy ;)")
}

```

## Arrays
In Rune you can define an array by binding it to a name:
```js
myArr = array{"hello", 10, false, 20, 30}
```
Arrays can have values of different types.

To access an array by index you can use the name followed by an index enclosed in square brackts (arrays are zero based):
```js
first = myArr[0]
println(first) # output: "hello"
```

## Tables
Tables in Rune are similar to hash maps or dictionaries in other languages with some added features.
In Rune you can define a table by binding it to a name:
```js
myTable = table{"key1": 1, "key2": false}
```
Tables can have values of different types, but keys must be of type `string`.

To access a table field by key you can use the name followed by an key enclosed in square brackts:
```js
second = myTable["key2"]
println(second) # output: false
```

>**Note**: keys are unique, this means adding a value with a key that already exists, the **existing value get overriden**:
```js
myTable = table{"uid": "10"}
println(myTable) # output: {"uid": 10}
myTable = append(myTable, "uid", "Hello World")
println(myTable) # output: {"uid": "Hello World"}
```

### Field Access
It is possible to access the fields of a table also via the `.`.
Both expressions are the same under the hood:
```js
myTable = table{"uid": "10"}
println(myTable["uid"]) # output: 10
println(myTable.uid) # output: 10
```

### The **`self`** argument
When defining a function on a table, a 'self' argument will be injected automatically.
'self' always refers to the table where the function was defined and called on.
```js
person = table{}
person.name = "John"
person.sayHello = fun(self) {
    println("Hello ", self.name)
}

person.sayHello() # prints: "Hello John"

person2 = new(person)
person2.name = "Jenny"
person2.sayHello() # prints: "Hello Jenny"
```

>**Copies vs References:** Remember that assigning a table to a variable creates a reference of it. Calling new however, creates a copy.
```js
person = table{}
person.name = "John"
person.sayHello = func(self) {
    prinln("Hello ", self.name)
}

person.sayHello() # prints: "Hello John"

person2 = person
person2.name = "Jenny"

person2.sayHello() # prints: "Hello Jenny"

person.sayHello() # ALSO PRINTS: "Hello Jenny"
```

In the above example, both variables `person` and `person2` refer to the **same underlying table**.

If you want to create a new copy of it and therefore get the (in this case) expected behavior, you can use the builtin `new` function:
```js
person = table{}
person.name = "John"
person.sayHello = func(self) {
    prinln("Hello ", self.name)
}

person.sayHello() # prints: "Hello John"

person2 = new(person) # create a copy of `person`
person2.name = "Jenny"

person2.sayHello() # prints: "Hello Jenny"

person.sayHello() # NOW CORRECTLY PRINTS: "Hello John"
```

>**Side note:** Arrays have the same behavior in terms of references, so the function `new` also works on arrays. 

### Whitespaces in keys
Whitespace in keys will be removed automatically to ensure they are accessable via the field access operator `.`.
```js
mytable = table{"key 1": false, "key 2": 10, "key 3": 3.14}
println(mytable) # prints: {"key1": false, "key2": 10, "key3": 3.14}
println(mytable.key3) # prints: 3.14
```

## Short-hand with `if` `then` and `elif` `then`
The `if`-statement can be written on one line with the `then` keyword:

Both if statements are the same:
```js
RuneIsAwesome = true;

if RuneIsAwesome then println("Rune is awesome") elif RuneIsAwesome != true then println("Oh to bad") else println("This should never print! :)");

if RuneIsAwesome {
    println("Rune is awesome")
}elif RuneIsAwesome != true {
    println("Oh to bad")
} else {
    println("This should never print! :)")
}
```

## Builin Functions

### print
- **Syntax**: `print(<arg1>, <arg2>, ...)`
- **Description**: Prints the given arguments to the standard out
- **Example**: `print("Hello, World times ", 10)`

### println
- **Syntax**: `println(<arg1>, <arg2>, ...)`
- **Description**: Prints the given arguments to the standard out and adds a newline character at the end
- **Example**: `print("Hello, World times ", 10)`

### wait
- **Syntax**: `wait(<milliseconds>)`
- **Description**: Waits the given amout of milliseconds
- **Example**: `wait(2000)`

### exit
- **Syntax**: `exit()`
- **Description**: Exits the execution immediately
- **Example**: `exit()`

### typeof
- **Syntax**: `typeof(<arg>)`
- **Description**: Returns the type name as a string of the given argument. Possible types are: `int`, `float`, `string`, `bool` and`unknown`.
- **Example**: `typeof(10) # returns "int"`

### append
- **Syntax**: `append(<array|table|string>, <value>)`
- **Description**: Appends the given value to the given array, table or string. Returns the new array, table or string.
- **Example**: `myArr = append(myArr, 10)`

### remove
- **Syntax**: `remove(<array|table|string>, <index>)`
- **Description**: Removed the given index from the given array, table or string. Returns the new array, table or string.
- **Example**: `myArr = remove(myArr, 2)`

### haskey
- **Syntax**: `haskey(<table>, <key>)`
- **Description**: Returns true if the given table has the given key, otherwise false.
- **Example**: `hasTheKey = haskey(myTable, "TheKey")`

### slice
- **Syntax**: `slice(<array|table|string>, <start>, <end>)`
- **Description**: Returns a slice of the given array, table, or string from the start index to the end index.
- **Example**: `slicedArray = slice(myArray, 1, 3)`

### slicefirst
- **Syntax**: `slicefirst(<array|table|string>, <end>)`
- **Description**: Returns a slice of the given array, table, or string from the start to the given end index.
- **Example**: `slicedArray = slicefirst(myArray, 2)`

### slicelast
- **Syntax**: `slicelast(<array|table|string>, <start>)`
- **Description**: Returns a slice of the given array, table, or string from the given start index to the end.
- **Example**: `slicedArray = sliceLast(myArray, 3)`

### len
- **Syntax**: `len(<array|table|string>)`
- **Description**: Returns the lenght of the given array, table or string.
- **Example**: `arrLen = len(myArr)`

### new
- **Syntax**: `new(<array|table>)`
- **Description**: Returns a deep copy of the given array or table.
- **Example**: `table2 = new(table1)`

## Editor Plugins
In the `editor` directory you will find plugins for different editors. Currently for _(help is welcome)_:
 - [VS Code](https://code.visualstudio.com/)

    ![](editor/vscode/rune/installed.png)

 The `readme.md` in each directory explains how to install them.
