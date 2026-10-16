// Package jinareader is a Go client for the Jina AI Reader API, generated
// with oapi-codegen from the OpenAPI document r.jina.ai publishes.
package jinareader

//go:generate bash ../scripts/fetch-spec.sh https://r.jina.ai/openapi.json reader.json
//go:generate bash ../scripts/fix-spec.sh fix-reader-spec.jq reader.json reader.fixed.json
//go:generate go tool oapi-codegen --config=config.yaml ../api/reader.fixed.json
//go:generate go tool oapi-codegen --config=config.models.yaml ../api/reader.fixed.json
