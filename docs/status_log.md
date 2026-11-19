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
* Add a failing test on the boundary - in execute - handle the logic when a user calls `pwd`
* Add a test on handleType to ensure `pwd` is recognised as a builtin.
* Consider refactoring this out as a handler to manage the logic of discovering working directory 
  * Consider starting with a failing test
  * Then write a test for handlePwd in isolation
  * Wire it up into the router
* Test to validate that when we have a competeing pwd as an external executable, we still pick the builtin - validates builtin takes priority

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

## 06-11-26

## Context

We now need to implement the builtin `cd` command to handle absolute paths. This command changes directory.

The cd command can handle different types of arguments:

* Abs paths, like /usr/local/bin
* Relative paths, like ./, ../, ./dir 
* ~ character, which represents the user's home directory

For this stage we'll focus on abs paths. This starts with "/" and specifes a location from the root of the filesystem.

## Contract

* If the directory exists, change to that directory
* If the directory doesn't exist, print cd: <directory>: No such file or directory
* If cd fails to execute, should be a no-op


## Example

$ pwd
/home/user
$ cd /usr/local/bin
$ pwd
/usr/local/bin
$ cd /does_not_exist
cd: /does_not_exist: No such file or directory
$ pwd
/usr/local/bin

## Completed

* Add failing test on the boundary to where we dispatch builtin commands for `cd` i.e. in TestRoutingCommand
* Implement the failing test - happy path


## 08-11-26

## Completed 
* Add test for cd handler on path existing
* Implement handler in builtin file
* Extend test to non existent paths: This should verify the behaviour of writing the error msg to output. Current path should be unchanged.
* Implement that hanlding of non-existent paths
* Wire up into the dispatcher in the execute() helper
* Debug why $ prompt is missing in output: returning errors stops execution/exits shell. 
* non-existent directories when calling `cd` should be graceful and not exit program. Perhaps only failing to write error to output causes exit.
* Add tests for "./", "../" , "./dir" - relative paths 

## 09-11-26

## Contract

* When a user passing the "~" tilder character, `cd` should switch to the users home directory
* The home directory is defined by the `$HOME` env variable

## Completed
* Test that "~" changes dir to the default home directory
* Test when changing the $HOME variable, this still obeys


## 10-11-26

## Shell quoting

In bash, quoting is used to control how the shell interprets characters and words. Its mainly used to remove the special meaning of certain characters or words to the shell. This can include special treatment for special characters, to prevent reserved words from being recognised as such and prevent parameter expansion.

To highlight:
* Single quotes '' preseve all characters including variable expansion like $1, $2 etc
* Double quotes "" preserve most characters but allow certain special characters like $ and \ to be intepreted

## The backslash escape character

The backslash character `\` can escape a single character, removing its special meaing, example:

$ echo three   spaces
three spaces
$ echo three\ \ \ spaces
three   spaces

Here, the backslash escapes each space, making them part of the string instead of a separator. This means the backslash preseves the literal value of the next character that follows, removing any special meaning it has, with the exception of newline.

## 12-11-26

Tests will verify with the echo command. So we can do manual user testing to verify behaviour.

## Contract

* when parsing a command with single quotes, all characters are preseved
* when parsing a command with spaces and tabs within single quotes, they are not collapsed or used as delimiters
* when parsing a command with quoted strings next to each other, they are concatenated to a single argument

## Notes

My issue here is that we are dealing with parsing the argument rather than the command. So whilst we are going to be testing for echo 'hello  world', to preseve space, this can occur for every commands argument. 


## 13-11-26

## Completed
* Write a failing test that shows single quotes preserves whitespace during an echo command 

## 15-11-26

## Completed
* Expand the failing test to include other special characters like $, \, %, !, ~.
* Add a test to show that two quoted strings next to eachother collapse into one i.e. 'hello''world' and 'hello' 'world'
* Add test that shows empty quotes are ignored hello''world goes to helloworld
* Create one integration test that proves the command handlers are recieving the parsed args before processing
* Wire it all up to the execute

## 18-11-26

## Completed
* Parser preserves whitespace with single quotes 
* Parser concatenates empty single quotes between characters
* Parser ignores empty single quotes
* Refactor command handlers to accept pargs (parsed args as an array)
* Stop the parser from hanging. echo 'unclosed currently leaves the shell stuck. Decide what it should do for this stage, then check that it returns promptly. As you trace the loop, ask whether every branch advances the input or stops parsing. The end desition was to return early and return an error with no closing quote found.
* make the closing quote graceful and log out error but not exit program
* Currently hello'''world prints hello'world instead of returning early
* Test edgecase empty quote at end of args

## 19-11-26

## Completed

* Refactor parseArgs and simplify the logic. Adding description where behaviour is not obvious.

## Context - Double Quotes

GNU docs: https://www.gnu.org/software/bash/manual/bash.html#Double-Quotes

* Most characters with "" (Double-Quotes) are treated literally
* Some characters get interpreted - $, \ (escaping)

## Contract - Within Double-Quotes

* Consecutive whitespace (inc. spaces, tabs) must be preserved
* Characters that act as delimiters or special characters lose their special meaning inside double quotes and are treated literally.
* Double-Quoted strings placed adjacent to each other are concatenated to form a single argument.
* Double-Quoted and unquoted strings placed next to eachother are also concatenated

## TODO 

First exersize is to create a set of behavioural tests based off the contract defined above for double quotes. Implementation can be done incrementally with a TDD approach:

### Test construction

TestDoubleQuotes:

* WithMulipleWhiteSpacePreserved - "hello   world"
* WithAdjacentQuotesConcatenated - hello""world
* WithQuotedAndUnquotedAdjacent - "hello"world
* WithUnquotedAndQuotedAdjacent - hello"world"
* WithQuotedSeparatedByWhitespace - "hello" "world"
* ContainingSingleQuotes - "sotiri's test"

