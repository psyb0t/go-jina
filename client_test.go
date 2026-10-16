package jina_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/psyb0t/go-jina"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testAPIKey = "test-key"

func newTestClient(
	t *testing.T,
	handler http.HandlerFunc,
) *jina.ClientWithResponses {
	t.Helper()

	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)

	auth := jina.WithRequestEditorFn(
		func(_ context.Context, req *http.Request) error {
			req.Header.Set("Authorization", "Bearer "+testAPIKey)

			return nil
		},
	)

	client, err := jina.NewClientWithResponses(srv.URL, auth)
	require.NoError(t, err)

	return client
}

func assertBearer(t *testing.T, r *http.Request) {
	t.Helper()

	assert.Equal(t, "Bearer "+testAPIKey, r.Header.Get("Authorization"))
}

func TestEmbeddings(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/v1/embeddings", r.URL.Path)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assertBearer(t, r)

		var body map[string]any

		assert.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, "jina-embeddings-v3", body["model"])
		assert.Equal(t, "retrieval.query", body["task"])
		assert.Equal(t, "hello", body["input"])

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(
			`{"model":"jina-embeddings-v3","object":"list",` +
				`"data":[{"object":"embedding","index":0,` +
				`"embedding":[0.5,-0.25]}],` +
				`"usage":{"total_tokens":3}}`,
		))
	})

	var input jina.EmbeddingsV3Request_Input

	require.NoError(t, input.FromEmbeddingsV3RequestInput0("hello"))

	task := jina.EmbeddingsV3RequestTaskRetrievalQuery

	var body jina.EmbeddingRequest

	require.NoError(t, body.FromEmbeddingsV3Request(jina.EmbeddingsV3Request{
		Model: jina.EmbeddingsV3RequestModelJinaEmbeddingsV3,
		Input: input,
		Task:  &task,
	}))

	res, err := client.EmbeddingsV1EmbeddingsPostWithResponse(
		context.Background(),
		body,
	)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode())
	require.NotNil(t, res.JSON200)

	assert.Equal(t, "jina-embeddings-v3", res.JSON200.Model)

	data, err := res.JSON200.Data.AsEmbeddingResponseData0()
	require.NoError(t, err)
	require.Len(t, data, 1)

	vec, err := data[0].Embedding.AsSingleEmbeddingDataEmbedding1()
	require.NoError(t, err)
	assert.Equal(t, []float64{0.5, -0.25}, vec)

	usage, err := res.JSON200.Usage.AsBaseUsage()
	require.NoError(t, err)
	assert.Equal(t, 3, usage.TotalTokens)
}

func TestRerank(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/v1/rerank", r.URL.Path)
		assertBearer(t, r)

		var body jina.TextRerankerRequest

		assert.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		assert.Equal(t, "what is a vector", body.Query)
		assert.Len(t, body.Documents, 2)

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(
			`{"model":"jina-reranker-v2-base-multilingual","object":"list",` +
				`"results":[{"index":1,"relevance_score":0.9},` +
				`{"index":0,"relevance_score":0.1}],` +
				`"usage":{"total_tokens":12}}`,
		))
	})

	docs := make([]jina.TextRerankerRequest_Documents_Item, 2)
	require.NoError(t, docs[0].FromTextRerankerRequestDocuments0("doc a"))
	require.NoError(t, docs[1].FromTextRerankerRequestDocuments0("doc b"))

	var body jina.RerankingRequest

	require.NoError(t, body.FromTextRerankerRequest(jina.TextRerankerRequest{
		Model:     jina.TextRerankerRequestModelJinaRerankerV2BaseMultilingual,
		Query:     "what is a vector",
		Documents: docs,
	}))

	res, err := client.RerankV1RerankPostWithResponse(
		context.Background(),
		body,
	)
	require.NoError(t, err)
	require.NotNil(t, res.JSON200)
	require.Len(t, res.JSON200.Results, 2)

	assert.Equal(t, 1, res.JSON200.Results[0].Index)
	assert.InDelta(t, 0.9, res.JSON200.Results[0].RelevanceScore, 0)
	assert.Equal(t, 12, res.JSON200.Usage.TotalTokens)
}

func TestGetModel(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/v1/models/jina-embeddings-v3", r.URL.Path)
		assertBearer(t, r)

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(
			`{"id":"jina-embeddings-v3","context_length":8192}`,
		))
	})

	res, err := client.GetModelV1ModelsModelIdGetWithResponse(
		context.Background(),
		"jina-embeddings-v3",
	)
	require.NoError(t, err)
	require.NotNil(t, res.JSON200)

	assert.Equal(t, "jina-embeddings-v3", res.JSON200.Id)
	assert.Equal(t, 8192, res.JSON200.ContextLength)
}
