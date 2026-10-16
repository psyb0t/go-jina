package jina_test

import (
	"context"
	"net/http"
	"os"

	"github.com/psyb0t/go-jina"
)

// This is the README usage example. It talks to the live API, so it has no
// Output comment and only has to compile.
//
//nolint:testableexamples // needs a real API key, compile check only
func Example() {
	auth := jina.WithRequestEditorFn(
		func(_ context.Context, req *http.Request) error {
			req.Header.Set("Authorization", "Bearer "+os.Getenv("JINA_API_KEY"))

			return nil
		},
	)

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

	res, err := client.EmbeddingsV1EmbeddingsPostWithResponse(
		context.Background(),
		body,
	)
	if err != nil {
		panic(err)
	}

	_ = res.JSON200
}
