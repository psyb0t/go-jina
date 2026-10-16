# Testing

Refines `~/.claude/rule-details/go/writing-tests.md`.

- **TST1** Tests are black-box, in packages `jina_test` and `jinareader_test`, and call the generated client against an `httptest` server.
- **TST2** Each test asserts the HTTP method, the path, the `Authorization: Bearer` header and, for a request body, the decoded JSON, then checks the decoded `JSON200`.
- **TST3** A spec patch that changes what the client sends (a hoisted request body, a renamed parameter, a retyped parameter) gets a test that would fail without it.
- **TST4** `example_test.go` holds the README usage example verbatim. It talks to the live API, so it only has to compile.
- **TST5** `make test` runs `go test -race` in the dev image.
