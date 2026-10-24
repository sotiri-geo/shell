# Status log

Recording where I was up to and what needs to be done.


## 15-10-26

* Extracted core prompt logic into run() helper which accepts both a io.Reader, io.Writer to leverage dependecy injection pattern for testing
* Currently the test is forever looping with the "$ " prompt literal. Investigate whether we are correctly exiting the function upon io.EOF being 
  returned when we read a new line.
* Add another test to show multi word commands also return the correct error message.

## 17-10-26

### Completed

* Changed from using `fmt.Print()` to `fmt.Fprint(writer, line)`. By using dependency injection with an io.Writer it 
  allowed me to use a `bytes.Buffer` for testing to write to and when using the application switch to stdout. By default fmt.Print() would write to 
  stdout causing issues with always printing out the `$ ` prompt to stdout.
* Termintating at `io.EOF` and writing remaining bytes to writer such that it renders correctly
* Added tests for cli to test entrypoint.
* Implemented test that builds and runs the binary and invokes the `exit` command checking that the subshell exits with status 0 (successful).

### Context

**Exit builtin**
The exit utility shall cause the shell to exit from its current execution environment with the exit status specified by the unsigned decimal integer n. If the current execution environment is a subshell environment, the shell shall exit from the subshell environment with the specified exit status and continue in the environment from which that subshell environment was invoked; otherwise, the shell utility shall terminate with the specified exit status. If n is specified, but its value is not between 0 and 255 inclusively, the exit status is undefined.

**Inspecting executables**

We can check which executables are builtin or external by using `type` command. To see all matches (both builtin and external) use `-a` flag

```sh
type -a echo
echo is a shell builtin
echo is /bin/echo
```

## 18-10-26

### Context

* The echo builtin command writes its args to stdout followed by a new line. If no args sent by user, it just writes a new line `\n`.
* See here for full description https://pubs.opengroup.org/onlinepubs/9699919799/utilities/echo.html

Note: Remember that we need to maintain the prompt on a loop, so echo commands when successful executed cannot exit the program.

### Completed

* Implement the echo builtin command with non empty args

### Refactor

* When the setup, execution and assertions are the same, but the only thing thats different is the input/output, convert to table tests.
* The run() helper is where the responsibilities are accumulating.
* Keep the run() function managing the REPL loop and have a separate execute() func that manages the switching between different commands.
* Have main() just managing the IO i.e. stdin and stdout
* A useful reference here is to separate the loops lifecycle in run() to the individual command behaviour. This gives smaller pieces to understand and test. The latter can be managed in a separate helper.
* Handling io.EOF in two ways: if it returns this err with line length of zero, means we have nothing left to read, so early return to caller. Also, we may end up with a situation where the reader returns io.EOF with some text read prior, which will need to be written to out. Then execute() gets called.
* Handle empty command, should loop back and allow the user to type another command

## 24-10-26


