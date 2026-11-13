# Personal Log

This contains any tips and tricks that I've learnt in order to handle issues more effectively.

* Break the problem down into compontents
* Try to write up a small contract that will hold regardless of the implementation e.g. "when a user does X, then Y happens"
* Try and get one complete behaviour working early before building out every handler
* Write a small failing test on the boundary between the components. Keep distinction between the responsibilites
* Add a small helper to make test path
* Keep feedback loop short and get a single behaviour running as early as possible


## Testing

When refactoring tests after you've refactored code, be mindful that the boundary of existing tests may change. In relation to what
we've done in this project:

* Initially execute() contained the logic for routing so we tested how we handle `type` in there i.e. main_test.go
* Then we extracted this logic into its own handler `handleType`
* Tested this through builtins_test.go, for the builtin and external command paths
* Now decided that the existing tests in main_test.go are no longer relevant to test the behaviour of the type handler through there
* We should instead focus that the routing is done correctly, i.e. make sure the responsibilites are being tested:
    * Does the user input reach the right handler with the right args


### Extending table tests vs t.Run()

A good measure of whether we should extend an existing table test vs creating another sub-test using t.Run() is whether the subject under
test requires a different setup/assertion pattern to the existing test cases. If you need a different env setup (or assertion pattern) but 
the test is closely related to the other tests, just add another t.Run() under that test group.

## Refactoring

* Let the test cases guide you on potential areas that may need refactoring 
* We've initially seen early during the dispatching of builtin handlers in the `execute()` func, that builtin names are managed in two places: One in the `execute()` body and another in `handleType` map. This risk of drift causes refactoring pressure, that is visible when writing tests across the boundaries. Because these lists of builtin names may become difficult to keep consistent.

## Identifying responsiblities

* When working on parsing single quotes, the question came up as to where this responsibility belonged.
* What should the command handler need to know about how the user wrote its original args. Should it need to be able to understand how to parsing single quotes, or is this agnositic to all commands. In this case, does the responsibility belong to the caller?
* Big question: How can you test that the shell ACTUALLY preserved the boundary within single quotes. Doing thing inside echo command can mask the underlinning operations. i.e. echo 'helloworld' and echo hello''world would return the same response, so it doesn't PROVE that the parsing actually occurred.
