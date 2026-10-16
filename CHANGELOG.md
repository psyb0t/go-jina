# Changelog

All notable changes per release. Versions follow [semver](https://semver.org).

## v0.1.0, 2026-10-16

- Added `jina`, a REST client generated with oapi-codegen from Jina AI's OpenAPI document at api.jina.ai (Jina Search Foundation API, spec version 2026.10.13.1624, OpenAPI 3.1.0). It covers all 10 paths: embeddings, reranking, models, chat completions and batch embedding jobs.
- Added `jinareader` in `reader/`, a REST client generated from the Reader API's OpenAPI document at r.jina.ai (Jina Reader API, spec version 0.5.0+a72eff0, OpenAPI 3.0.0). It covers all 3 paths: the service index and URL reads with the URL in a JSON body (`/`), URL reads with the URL in the path, and the short-lived files a read can produce (`/ephemeral/{encodedFid}`).
- Number fields are `float64` instead of oapi-codegen's default `float32`.
- `make generate` downloads the current specs, rewrites the parts oapi-codegen cannot generate (null branches in `anyOf`, a duplicate enum type name, union request bodies that would encode as `{}`, a `url` parameter that shadows `net/url`, and untyped parameters and maps), and regenerates both clients. The specs are not committed.
- Added the `go-jina` agent skill. Install it with `claude plugin install go-jina@psyb0t` (Claude Code), `codex plugin add go-jina@psyb0t` (Codex) or `openclaw skills install @psyb0t/go-jina` (ClawHub).
