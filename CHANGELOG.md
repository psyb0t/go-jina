# Changelog

All notable changes per release. Versions follow [semver](https://semver.org).

## v0.4.3, 2026-10-30

- Added the `FinalUrl` field (`*string`, JSON `finalUrl`) to `FormattedPageDto0x43` in `jinareader`.
- The documented default of the `Base` field changed from 'initial' to 'final' in `CrawlerOptions0x28`, `CrawlerOptionsHeaderOnly0x3b` and the `Post` and `PostTargetUrl` JSON, form data and multipart body types. Only the doc comments changed, the Go types are the same.
- `jinareader` generates from Jina Reader API spec version 0.5.0+9e23d86.

## v0.4.2, 2026-10-28

- Updated Semgrep in the development image to 1.182.1 and refreshed its locked dependencies. certifi 2026.10.23 drops the emSign root CAs whose intermediate keys were compromised (GHSA-6q48-6c7w-9gw2).
- Dependabot no longer opens pull requests for `mcp` or the OpenTelemetry packages in `requirements-semgrep.txt`. Semgrep pins their versions itself, so bumping one of them on its own cannot install.
- Updated kin-openapi, which the code generator uses, from 0.142.0 to 0.149.0. The newer version fixes several panics on malformed specs reported in GitHub security advisories.
- The security scan merges its two gitleaks runs into one SARIF run. GitHub code scanning rejects a report with two runs from the same tool, so no scan results had been uploaded.

## v0.4.1, 2026-10-27

- The development image is based on Debian (`golang:1.26.10-trixie`) instead of Alpine. Alpine keeps only the latest build of each package, so the exact pins broke the image build whenever a package was updated. System packages now come from the snapshot.debian.org date the base image was built from, which keeps them fixed for a given base image digest. jq 1.8.2, shellcheck 0.11.0 and shfmt 3.13.1 keep their versions, and Semgrep runs in its own virtual environment.

## v0.4.0, 2026-10-24

- Added Reader API endpoints `GET /bot` and `GET /.well-known/http-message-signatures-directory` to `jinareader`: `Client.GetBot`, `Client.GetWellKnownHttpMessageSignaturesDirectory`, the `WithResponse` variants on `ClientWithResponses`, `NewGetBotRequest`, `NewGetWellKnownHttpMessageSignaturesDirectoryRequest`, `ParseGetBotResponse`, `ParseGetWellKnownHttpMessageSignaturesDirectoryResponse` and the response types `GetBotResponse` and `GetWellKnownHttpMessageSignaturesDirectoryResponse`.
- Added the `WebBotAuth` field to `GetParams`, `GetTargetUrlParams`, `PostParams`, `PostTargetUrlParams` and the `Post` and `PostTargetUrl` JSON, form data and multipart body types.
- `RobotsTxt` on those same params and body types changed from `*string` to a union type, for example `*PostJSONBody_RobotsTxt`. Code that sets `RobotsTxt` to a `*string` must change.
- `ClientInterface` and `ClientWithResponsesInterface` gained `GetBot`, `GetWellKnownHttpMessageSignaturesDirectory` and their `WithResponse` equivalents. Code that implements either interface, such as a mock, must add them.
- The position-based names oapi-codegen gives anonymous and union types in `jinareader` shifted again. Types such as `Cookie0x17`, `CrawlerOptions0x27` and `FormattedPageDto0x42` were renamed one slot up (for example `Cookie0x18`, `CrawlerOptions0x28`), and the numbered member aliases and `From`, `Merge` and `As` helpers on them were removed or renamed. Code that names these generated types directly must update it.
- The Reader spec rewrite in `scripts/fix-reader-spec.jq` now replaces a `oneOf` member that wraps a single `$ref` in `allOf` with the bare `$ref`. This stops oapi-codegen from emitting a per-member alias for each such property, such as `PostParamsWebBotAuth0`.
- `jinareader` generates from Jina Reader API spec version 0.5.0+2e3f8b4, up from 0.5.0+b2f8a45. The `jina` package and its spec version are unchanged.

## v0.3.0, 2026-10-22

- Added the Reader API's `/echo/{rest}` endpoint to `jinareader` for the `GET`, `POST`, `PUT`, `PATCH` and `DELETE` methods, and the matching `/echo` endpoint without the path parameter.
- Added client methods on `Client`: `GetEcho`, `PostEcho`, `PutEcho`, `PatchEcho`, `DeleteEcho`, `GetEchoRest`, `PostEchoRest`, `PutEchoRest`, `PatchEchoRest` and `DeleteEchoRest`, plus `WithBody` and `WithFormdataBody` variants for the `POST`, `PUT` and `PATCH` ones.
- Added the same methods with a `WithResponse` suffix on `ClientWithResponses`, for example `GetEchoRestWithResponse`, and the response types `GetEchoResponse`, `PostEchoResponse`, `PutEchoResponse`, `PatchEchoResponse`, `DeleteEchoResponse` and their `EchoRest` equivalents.
- Added request builders `NewGetEchoRequest`, `NewPostEchoRequest`, `NewPutEchoRestRequest` and the other `NewXxxEchoRequest` and `NewXxxEchoRestRequest` functions, and parsers such as `ParseGetEchoResponse` and `ParsePostEchoRestResponse`.
- Added request body types such as `PostEchoJSONBody`, `PostEchoFormdataBody`, `PostEchoMultipartBody` and `PostEchoRestJSONRequestBody`, with equivalents for `PUT` and `PATCH`.
- `ClientInterface` and `ClientWithResponsesInterface` gained the new methods. Code that implements either interface, such as a mock, must add them. Code that only calls `Client` or `ClientWithResponses` is unaffected.
- `jinareader` generates from Jina Reader API spec version 0.5.0+b2f8a45. The `jina` package and its spec version are unchanged.

## v0.2.2, 2026-10-21

- The module requires Go 1.26.10, up from 1.26.9, and the development image uses it. Go 1.26.10 fixes GO-2026-6809 in `crypto/x509`, which govulncheck reports for every package that makes HTTPS requests, so the v0.2.1 release pipeline failed its security scan and v0.2.1 was never published.
- Includes the v0.2.1 changes: the spec rewrite handles the Reader API's `/echo/{rest}` endpoint, and the issue-mirror pull is staggered.
- Generated from the same spec versions as v0.2.0.

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
