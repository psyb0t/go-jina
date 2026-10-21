# Changelog

All notable changes per release. Versions follow [semver](https://semver.org).

## v0.2.1, 2026-10-21

- Regeneration no longer fails on the Reader API's new `/echo/{rest}` endpoint. Its request bodies declare an empty `allOf`, which made oapi-codegen panic, and its `rest` path parameter has no schema. The spec rewrite now drops empty `allOf` arrays and gives schemaless parameters `type: string`. The generated packages are unchanged in this release; the next regeneration picks up the endpoint.
- The scheduled issue-mirror pull runs at a staggered time instead of a shared minute.
- Generated from the same spec versions as v0.2.0.

## v0.2.0, 2026-10-19

- In `jinareader`, `Page` on `GetParams`, `GetTargetUrlParams`, `PostParams`, `PostTargetUrlParams`, `PostJSONBody`, `PostFormdataBody`, `PostMultipartBody` and their `PostTargetUrl...` equivalents changed from `*float64` to `*PageSelection0x0f`. `PageSelection0x0f` takes a page number, a string range such as "1,3,5-7", or a list of numbers and ranges, matching the Reader API's page selection syntax. Code that sets `Page` to a `*float64` must change to build a `PageSelection0x0f` instead.
- Adding that page selection type shifted the position-based names oapi-codegen gives every other generated anonymous and union type in `jinareader`. Most of the package's other exported identifiers were renamed one slot up, for example `Cookie0x16` is now `Cookie0x17`, `Enum0x14` is now `Enum0x15`, and `FormattedPageDto0x41` is now `FormattedPageDto0x42`. Code that names one of these generated types directly must update it; code that only calls the client methods and reads named fields is unaffected.
- `jinareader` now generates from Jina Reader API spec version 0.5.0+eaf3e28, up from 0.5.0+a72eff0. The `jina` package and its spec version are unchanged.

## v0.1.2, 2026-10-17

- Documentation only. The agent skill explains how to look up every endpoint, parameter and model with `go doc` or on pkg.go.dev, and declares its license.
- The skill no longer lists the API credentials under OpenClaw's `requires.env`, which kept OpenClaw from loading it until they were set. It requires only `go`; the credentials are still documented as what your program reads.
- The README points to the full API reference, and the examples in the README and the skill were checked against the current code.

## v0.1.1, 2026-10-17

- Documentation only. The README and `AGENTS.md` describe the daily regeneration workflow: what it checks, how it fixes failing checks, how it picks the next version and which paths a fix may not touch.

## v0.1.0, 2026-10-16

- Added `jina`, a REST client generated with oapi-codegen from Jina AI's OpenAPI document at api.jina.ai (Jina Search Foundation API, spec version 2026.10.13.1624, OpenAPI 3.1.0). It covers all 10 paths: embeddings, reranking, models, chat completions and batch embedding jobs.
- Added `jinareader` in `reader/`, a REST client generated from the Reader API's OpenAPI document at r.jina.ai (Jina Reader API, spec version 0.5.0+a72eff0, OpenAPI 3.0.0). It covers all 3 paths: the service index and URL reads with the URL in a JSON body (`/`), URL reads with the URL in the path, and the short-lived files a read can produce (`/ephemeral/{encodedFid}`).
- Number fields are `float64` instead of oapi-codegen's default `float32`.
- `make generate` downloads the current specs, rewrites the parts oapi-codegen cannot generate (null branches in `anyOf`, a duplicate enum type name, union request bodies that would encode as `{}`, a `url` parameter that shadows `net/url`, and untyped parameters and maps), and regenerates both clients. The specs are not committed.
- Added the `go-jina` agent skill. Install it with `claude plugin install go-jina@psyb0t` (Claude Code), `codex plugin add go-jina@psyb0t` (Codex) or `openclaw skills install @psyb0t/go-jina` (ClawHub).
