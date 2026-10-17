# go-jina

[![Go Reference](https://pkg.go.dev/badge/github.com/psyb0t/go-jina.svg)](https://pkg.go.dev/github.com/psyb0t/go-jina)
[![CI](https://github.com/psyb0t/go-jina/actions/workflows/pipeline.yml/badge.svg?branch=main)](https://github.com/psyb0t/go-jina/actions/workflows/pipeline.yml)
[![version](https://raw.githubusercontent.com/psyb0t/go-jina/badges/version.svg)](https://github.com/psyb0t/go-jina/tags)
[![license](https://raw.githubusercontent.com/psyb0t/go-jina/badges/license.svg)](LICENSE)
[![imported by](https://raw.githubusercontent.com/psyb0t/go-jina/badges/importers.svg)](https://github.com/psyb0t/go-jina/blob/badges/importers.md)

Unofficial Go client for the [Jina AI APIs](https://jina.ai), generated with [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen) from the OpenAPI documents Jina AI publishes. I am not affiliated with Jina AI.

There are two packages:

| Package | Import path | API |
|---|---|---|
| `jina` | `github.com/psyb0t/go-jina` | `https://api.jina.ai`: embeddings, reranking, models, chat completions (document OCR) and batch embedding jobs |
| `jinareader` | `github.com/psyb0t/go-jina/reader` | `https://r.jina.ai`: the Reader API, which fetches a URL and returns it as markdown, text, HTML or JSON |

Everything in the `*.gen.go` files is generated. Don't edit them by hand, regenerate them.

## Install

```sh
go get github.com/psyb0t/go-jina
```

## Usage

Every endpoint in the specs has a typed method on `ClientWithResponses`. Responses decode into `JSON200` and the other status fields.

```go
package main

import (
	"context"
	"net/http"
	"os"

	"github.com/psyb0t/go-jina"
)

func main() {
	auth := jina.WithRequestEditorFn(func(_ context.Context, req *http.Request) error {
		req.Header.Set("Authorization", "Bearer "+os.Getenv("JINA_API_KEY"))

		return nil
	})

	client, err := jina.NewClientWithResponses("https://api.jina.ai", auth)
	if err != nil {
		panic(err)
	}

	var input jina.EmbeddingsV3Request_Input
	if err := input.FromEmbeddingsV3RequestInput0("hello world"); err != nil {
		panic(err)
	}

	var body jina.EmbeddingRequest
	if err := body.FromEmbeddingsV3Request(jina.EmbeddingsV3Request{
		Model: jina.EmbeddingsV3RequestModelJinaEmbeddingsV3,
		Input: input,
	}); err != nil {
		panic(err)
	}

	res, err := client.EmbeddingsV1EmbeddingsPostWithResponse(context.Background(), body)
	if err != nil {
		panic(err)
	}

	_ = res.JSON200
}
```

The same code is in `example_test.go`, so `make test` checks that it compiles.

The embedding and rerank request bodies are unions with one variant per model family (`EmbeddingsV3Request`, `ClipV2Request`, `TextRerankerRequest`, `RerankerV3Request` and so on). Fill the variant for your model and set it on `EmbeddingRequest` or `RerankingRequest` with the matching `From...` method. Fields that accept more than one shape, such as `input` or a response's `embedding`, work the same way with `From...` and `As...` methods.

Method names come from the spec's operation IDs, which spell out the path and method: `EmbeddingsV1EmbeddingsPost`, `RerankV1RerankPost`, `ListModelsV1ModelsGet`, `CreateBatchJobV1BatchEmbeddingsPost` and so on.

### Reader

```go
reader, err := jinareader.NewClientWithResponses("https://r.jina.ai", readerAuth)
if err != nil {
	panic(err)
}

accept := "application/json"

res, err := reader.GetTargetUrlWithResponse(
	ctx,
	"https://example.com",
	&jinareader.GetTargetUrlParams{Accept: &accept},
)
```

`res.JSON200.Data` holds the page title, URL and content. Without `Accept: application/json` the API answers in plain text, which is in `res.Body`. The Reader options (`X-Respond-With`, `X-Target-Selector`, `X-With-Links-Summary` and the rest) are fields on the params struct. `PostWithResponse` takes the URL and options in a JSON body instead.

## Auth

Both APIs take a Jina AI API key as a bearer token: `Authorization: Bearer <key>`. Get one from the [Jina AI API dashboard](https://jina.ai/api-dashboard/key-manager). Set it with a request editor, as in the example above. `jinareader` has its own `WithRequestEditorFn` with the same shape.

`GET /v1/models` and `GET /v1/models/{model_id}` work without a key. The Reader API also answers without a key, at a lower rate limit.

## Regenerating

The specs are not committed. `make generate` downloads the current OpenAPI documents into the gitignored `api/` directory and regenerates both clients:

```sh
make generate
```

`gen.go` and `reader/gen.go` declare the steps, so `make generate` is just `go generate ./...`:

- `scripts/fetch-spec.sh` downloads `https://api.jina.ai/openapi.json` to `api/openapi.json` and `https://r.jina.ai/openapi.json` to `api/reader.json`.
- `scripts/fix-spec.sh` runs each spec through a jq filter and writes `api/openapi.fixed.json` and `api/reader.fixed.json`.
  - `scripts/fix-jina-spec.jq`: the main spec is OpenAPI 3.1 written by pydantic. It turns `anyOf: [X, {type: null}]` into X plus `nullable: true`. It moves the `embedding_type` enum, which is allowed as one value or a list, into a shared `EmbeddingType` schema, because inline oapi-codegen gives both forms the same type name and stops. It also moves the embedding and rerank request bodies into the `EmbeddingRequest` and `RerankingRequest` schemas. Inline, the generated body type loses the union's JSON encoder and the client would send `{}`.
  - `scripts/fix-reader-spec.jq`: renames the `{url}` path parameter to `{targetUrl}`, because a Go parameter named `url` shadows `net/url` in the generated client. It drops `nullable` from untyped `additionalProperties`, which makes oapi-codegen write code that does not compile. It types untyped query and header parameters such as `jsonSchema` as strings, because otherwise the client always sends them and panics when they are unset.
- oapi-codegen runs with `config.yaml` (and `reader/config.yaml` plus `reader/config.models.yaml`). The configs map `number` to `float64`, because the specs give scores and timeouts no format and oapi-codegen would otherwise use `float32`. The Reader spec repeats around 100 parameters inline on every operation, so its models go in `reader/models.gen.go` and the client in `reader/client.gen.go`.

`.github/workflows/regen.yml` does the same every day in CI. When the new spec changes the client, it runs `make lint`, `make test` and `make sec`. If a check fails, Claude edits the code and the checks run again, up to three times, and a review step reads the result. A passing change is committed with a CHANGELOG entry and tagged as the next release: a patch version for fixes and additions, a minor version for breaking changes.

## Development

Every command runs in the `Dockerfile.dev` image. The host needs only Docker and Make.

```sh
make help      # list targets
make dep       # go mod tidy + vendor
make lint      # shellcheck, shfmt, go fix, golangci-lint
make test      # go test -race
make sec       # gitleaks, govulncheck, semgrep
make generate  # download the specs and regenerate
```

## Agent integrations

The [skill](.agents/skills/go-jina) teaches a coding agent how to use go-jina when it writes Go. It works in any agent that reads `.agents/skills/`, and installs natively in the clients below.

### Claude Code

```bash
claude plugin marketplace add psyb0t/agents
claude plugin install go-jina@psyb0t
```

### Codex

```bash
codex plugin marketplace add psyb0t/agents
codex plugin add go-jina@psyb0t
```

Codex also picks the skill up with no install in any checkout of this repo.

### OpenClaw

The skill is published to ClawHub on every release:

```bash
openclaw skills install @psyb0t/go-jina
```

## License

[MIT](LICENSE). The generated code comes from Jina AI's published API specifications. The specifications themselves are not redistributed here.
