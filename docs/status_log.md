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

## Context

Ref: https://app.codecrafters.io/courses/shell/stages/ez5

We are required to implement the `type` command in our interpreter.

docs: https://pubs.opengroup.org/onlinepubs/9799919799/utilities/type.html

definition: The type utility shall indicate how each argument would be interpreted if used as a command name. Typically identifies the operand as either shell builtin, function, alias, or keyword, and where applicable, display its pathname.

## Example 

For builtin commands (like echo, exit, and type), print <command> is a shell builtin.
For unrecognized commands that don't match any builtin, print <command>: not found.

$ type echo
echo is a shell builtin
$ type exit
exit is a shell builtin
$ type type
type is a shell builtin
$ type invalid_command
invalid_command: not found
$

## Completed

* Implemented type command
* Extacted command handlers out into separate file builtin, but kept under the same package main importing and testing under same build.
* Duplication exists: command names between the switch statement in `execute` and the handler `handleType` have both got references to builtin. Decide how to sync the two. Maybe use a map to iterate and have that source centralised as a constant.
* Add tests for builtin handlers like `echo`, `type` that will focus on command behaviour. Keep existing run tests as the integration test, checking loop behaviour.


## 25-10-26

## Context

We need to extend the interpreter to use the environment variable `PATH`, which contains a list of directories to search for executable programs. Normally what happens, when a user enters a command that is not a builtin command, a shell searches the directories in `PATH` for an executable file with that name.

## Example of what PATH prints to stdout

# Print the PATH variable
$ echo $PATH
/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin


Note: directory paths are separated with colon `:`.

A Shell will search through the paths defined from left to right (ordering is important). It will stop at the first directory that it discovers the named executable the user requests. When a shell exhausts every directory path and cannot find the executable command, then the following error is returned:

$ myprogram
bash: myprogram: command not found

Usually indicating to the user that the binary isn't installed or is in a directory not included in `PATH` variable.

GNU docs: https://www.gnu.org/software/bash/manual/bash.html#index-PATH

## Completed

We are required to extend the type command such that:

* When we encounter builtin command, report it as builtin (<commad> is a shell builtin) and stop.
* If the command is not a builtin, the shell should go through every directory in `PATH`:
  * Check if a file with the command name exists
  * Check if the file has execute permissions 
  * If the file exists and has execute permissions, print <command> is <full_path> and stop
  * If the file exists but lacks execute permissions, skip it and continue to the next directory
  * If no executable is found in any directory for the given command name, then print <command>: not found

$ type grep
grep is /usr/bin/grep
$ type invalid_command
invalid_command: not found
$ type echo
echo is a shell builtin

1. Small contract:
   - Builtins take precedence over executables in PATH.
   - Search directories in their listed order.
   - Return the first matching executable file.
   - Skip unsuitable candidates and continue searching.
   - When nothing matches, print the existing “not found” response and continue the shell.


## Refactor (optional)

* Add test coverage for user calling type on a executable command through run() in main_test.go


## 28-10-26

## Completed

* Refactor the tests found for run() and execute() in main_test.go to now test the responsibilities of the routing 
* Remove any redundtant tests found already in the builtin_test.go
* Added files external.go and external_test.go to extract helper to separate the responibility of running external command
* Running end to end test through run() to prove behaviour works - I.e. Running external command

## Context

Next, we need to extend our project to handle running external programs. When a command isn't a builtin, the shell should:

* Determine if a given command is executable
* If so, execute the command
* Pass any user args from the commandline to the program

E.g. `custom_exec arg1 arg2` should eval: `custom_exec` check under `PATH` if its an executable. Then run this with `arg1`, `arg2`.

Contract: 

* If the first arg is an executable it should execute it with any subsequent user args
* Skip if first arg is not an executable found in PATH

## 30-10-26

## Completed

* Add test to check when external command doesn't exist - fallback
* Add a test to check external command exists but not on PATH - fallback
* Add a test to check multiple args are accepted with helper - Do implementation with TDD
* findExecutable to handle relative paths "." with "./executable" 
* Add findExectuable with leading and trailing path list separators i.e. ":file/to/path" and "file/to/path:"
* Wire it up with the router in execute()
* fetch the base in runCommand to only run cmd() with the executable name rather than full path. This is because the dir 
  already lives in PATH env variable.

## 31-10-26

## Completed

* fetch the base in runCommand to only run cmd() with the executable name rather than full path. This is because the dir 
  already lives in PATH env variable.

## Context

We require to implement the `pwd` builtin command which prints the full absolute path of the current working directory. Normally, 
when a shell starts, its current working dir is typically the dir from which is was executed. The `pwd` implementation needs to 
retrieve this info from the operating system and print it.

## Example

$ pwd
/home/user/projects
$ pwd
/usr/local/bin

## Contract

* When a user executes pwd command then abs path of the current working dirctory should be returned
* This is a builtin command and should be prioritised over external executables

## TODO 

* Add a failing test on the boundary - in execute - handle the logic when a user calls `pwd`
* Consider refactoring this out as a handler to manage the logic of discovering working directory 
  * Consider starting with a failing test
* Consider refactoring execute into routing both builtin and external - Ask codex whether this would be a good approach
