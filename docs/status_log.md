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

### Todo 

* 
