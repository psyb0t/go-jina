// Package jina is a Go client for the Jina AI APIs (embeddings, reranking,
// models, chat completions and batch jobs), generated with oapi-codegen from
// the OpenAPI document Jina AI publishes. The Reader API is in the reader
// subpackage.
package jina

//go:generate bash scripts/fetch-spec.sh https://api.jina.ai/openapi.json openapi.json
//go:generate bash scripts/fix-spec.sh fix-jina-spec.jq openapi.json openapi.fixed.json
//go:generate go tool oapi-codegen --config=config.yaml api/openapi.fixed.json
