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
  * We are more concerned with WHAT was passed into echo i.e. did 'helloworld' and hello''world both resolve as helloworld before being
    dipatched to the echoHandler?

## Terminology

With reference to the execute command, we've discussed this difference in the switch statement part whether to call it a router or a dispatcher:
* Router - emphasises choosing the desitination (where is should go based off the command name)
* Dispatcher - emphasises handing off the work to the designated handler to be invoked i.e. "echo" matches to the "echoHandler"

Both of these often happen at the same time. If you split these two into separaet functions - a lookup could return the handler (router) and the dispatcher could call that handler. However, in this codebase I would say "execute dispatches commands", the term describes what the code does, regardless of the implementation.

## Effectively naming tests

Name each table case after the behaviour you’d want to see in the failure log. The input shows how you trigger it, and the assertion shows the exact result. The name tells the reader why the case exists.

For your parser tests:

| Less helpful | Better | What the better name tells you |
|---|---|---|
| WithThreeBackslashes | KeepsEscapedSpacesInOneArgument | The test checks whether escaped spaces split an argument. |
| WithAdjacentQuotes | ConcatenatesAdjacentQuotedFragments | The fragments should form a single argument. |
| WithUnclosedSingleQuote | RejectsUnclosedSingleQuote | The parser should return an error. |
| WithUnclosedSingleQuote under TestRun | ContinuesAfterUnclosedSingleQuote | The shell should execute the next command. |
| WithWhitespace | PreservesWhitespaceInsideDoubleQuotes | Both the kind of whitespace behaviour and where it applies are clear. |

The two unmatched quote cases are a good example. The same input can exercise different responsibilities. A parser test can prove the error is returned. A REPL test earns ContinuesAfterUnclosedSingleQuote by sending a second command and checking its output. If it only checks the diagnostic, ReportsUnclosedSingleQuote is the more accurate name.

When you’re stuck on a name, write this sentence first: “If this test fails, it means the parser no longer ___.” For three\ \ \ spaces, you might fill in “keeps escaped spaces in one argument.” That gives you KeepsEscapedSpacesInOneArgument. If you can only fill in “handles three backslashes,” ask what distinguishes the correct result from the likely bug. That distinction is what belongs in the name.


## Cursor based struct

Your Idea: Struct + Helper Methods (The Scanner / Lexer Pattern)

Your proposed refactor is a textbook idiomatic Go pattern: an unexported scanner/parser struct behind a plain function API.
```go
    // Public/package API stays identical — all existing tests remain your safety net.
    func parseArgs(input string) ([]string, error) {
        p := &argParser{input: input}
        return p.parse()
    }

    type argParser struct {
        input string
        pos   int
        buf   strings.Builder
        args  []string
    }
```
 
Why this is idiomatic in Go

 • Zero API churn: Callers don't need to construct a parser struct; it's purely an internal implementation detail to avoid threading (input, pos, &buf) through every helper (like getClosingIndex currently has to).
 • Single cursor + tiny bounds-safe primitives: Give the struct 2–3 tiny navigation methods:
 ```go
    func (p *argParser) eof() bool     { return p.pos >= len(p.input) }
    func (p *argParser) peek() byte    { return p.input[p.pos] } // or (byte, bool)
    func (p *argParser) advance() byte { ch := p.input[p.pos]; p.pos++; return ch }
```

Then each syntactic rule becomes a small, self-contained method:
• readSingleQuoted() error
• readDoubleQuoted() error
• readEscape() error
• flushArg() / skipSpaces()


