# Status log

Recording where I was up to and what needs to be done.


## 15-10-26

* Extracted core prompt logic into run() helper which accepts both a io.Reader, io.Writer to leverage dependecy injection pattern for testing
* Currently the test is forever looping with the "$ " prompt literal. Investigate whether we are correctly exiting the function upon io.EOF being 
  returned when we read a new line.
* Add another test to show multi word commands also return the correct error message.
