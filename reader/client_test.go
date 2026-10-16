package jinareader_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	jinareader "github.com/psyb0t/go-jina/reader"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testAPIKey    = "test-key"
	testTargetURL = "https://example.com/page"
	jsonMediaType = "application/json"
)

func newTestClient(
	t *testing.T,
	handler http.HandlerFunc,
) *jinareader.ClientWithResponses {
	t.Helper()

	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)

	auth := jinareader.WithRequestEditorFn(
		func(_ context.Context, req *http.Request) error {
			req.Header.Set("Authorization", "Bearer "+testAPIKey)

			return nil
		},
	)

	client, err := jinareader.NewClientWithResponses(srv.URL, auth)
	require.NoError(t, err)

	return client
}

func writePage(w http.ResponseWriter) {
	w.Header().Set("Content-Type", jsonMediaType)
	_, _ = w.Write([]byte(
		`{"code":200,"status":20000,` +
			`"data":{"title":"Example","url":"https://example.com/page",` +
			`"content":"# Example"}}`,
	))
}

func TestGetTargetURL(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/"+testTargetURL, r.URL.Path)
		assert.Equal(t, "Bearer "+testAPIKey, r.Header.Get("Authorization"))
		assert.Equal(t, jsonMediaType, r.Header.Get("Accept"))
		assert.Equal(t, "markdown", r.Header.Get("X-Respond-With"))

		writePage(w)
	})

	accept := jsonMediaType
	respondWith := "markdown"

	res, err := client.GetTargetUrlWithResponse(
		context.Background(),
		testTargetURL,
		&jinareader.GetTargetUrlParams{
			Accept:       &accept,
			XRespondWith: &respondWith,
		},
	)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode())
	require.NotNil(t, res.JSON200)
	require.NotNil(t, res.JSON200.Data)

	assert.Equal(t, "# Example", *res.JSON200.Data.Content)
	assert.Equal(t, testTargetURL, res.JSON200.Data.Url)
}

func TestPostURLInBody(t *testing.T) {
	t.Parallel()

	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/", r.URL.Path)
		assert.Equal(t, "Bearer "+testAPIKey, r.Header.Get("Authorization"))
		assert.Equal(t, jsonMediaType, r.Header.Get("Content-Type"))

		var body jinareader.PostJSONBody

		assert.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		require.NotNil(t, body.Url)
		assert.Equal(t, testTargetURL, *body.Url)

		writePage(w)
	})

	accept := jsonMediaType
	target := testTargetURL

	res, err := client.PostWithResponse(
		context.Background(),
		&jinareader.PostParams{Accept: &accept},
		jinareader.PostJSONRequestBody{Url: &target},
	)
	require.NoError(t, err)
	require.NotNil(t, res.JSON200)
	require.NotNil(t, res.JSON200.Data)

	assert.Equal(t, "Example", *res.JSON200.Data.Title)
}
