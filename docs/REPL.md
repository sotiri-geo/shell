# What is the REPL

It's acronym stands for "read-eval-print-loop", also termed language shell. It's a simple interactive computer programming environment that takes a single user inputs, executes them and returns the result to the user. Normally a progam written in a REPL environment is executed piecewise.

# Overview

In a REPL env:

* User enters one or more expressions and the REPL evaluates them and displays the results.

An example would be:

1. Read - show `$ ` prompt and read the user's command
2. Eval - interpret and execute it, such as pwd
3. Print - display the result, such as the current directory 
4. Loop - Show the prompt again and wait for another command

```sh
$ pwd                  ← Read the command
/Users/sotiri/project  ← Execute it and display the result
$                      ← Loop back to the prompt
```

In terms of what this might look like in a Go progam:


```go
for {
    // Read input
    // Interpret and execute the command
    // Display the result
}
```
