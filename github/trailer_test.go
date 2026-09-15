package github

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/google/go-github/v45/github"
	"github.com/stretchr/testify/require"
)

// newTestClient points a go-github client at a local test server.
func newTestClient(t *testing.T, handler http.Handler) *github.Client {
	t.Helper()
	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	baseURL, err := url.Parse(server.URL + "/")
	require.NoError(t, err)

	client := github.NewClient(nil)
	client.BaseURL = baseURL
	return client
}

// TestSetPRTrailerDetailsPatchesBodyOnly makes sure the PATCH carries no base
// branch. GitHub rejects a base field on a PR that is part of a stack, even
// when the value does not change.
func TestSetPRTrailerDetailsPatchesBodyOnly(t *testing.T) {
	var patched map[string]any

	handler := http.NewServeMux()
	handler.HandleFunc("/repos/braincorp/titanium/pulls/26457", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			_ = json.NewEncoder(w).Encode(map[string]any{
				"number":                26457,
				"title":                 "some title",
				"state":                 "open",
				"body":                  "existing body",
				"maintainer_can_modify": true,
				"base":                  map[string]any{"ref": "ian/parent-branch"},
			})
		case http.MethodPatch:
			raw, err := io.ReadAll(r.Body)
			require.NoError(t, err)
			require.NoError(t, json.Unmarshal(raw, &patched))
			_ = json.NewEncoder(w).Encode(map[string]any{"number": 26457})
		default:
			t.Errorf("unexpected method %s", r.Method)
		}
	})

	client := newTestClient(t, handler)

	err := setPRTrailerDetails(context.Background(), client, "braincorp/titanium", 26457, "pulumi output", "some details", "pulumi-preview")
	require.NoError(t, err)

	require.Contains(t, patched, "body")
	require.Contains(t, patched["body"], "existing body")
	require.Contains(t, patched["body"], `<details id="pulumi-preview">`)

	for _, field := range []string{"base", "title", "state", "maintainer_can_modify"} {
		require.NotContains(t, patched, field)
	}
}

// TestSetPRTrailerDetailsReplacesExistingTrailer makes sure a second run
// replaces the trailer instead of appending another one.
func TestSetPRTrailerDetailsReplacesExistingTrailer(t *testing.T) {
	var patched map[string]any

	handler := http.NewServeMux()
	handler.HandleFunc("/repos/braincorp/titanium/pulls/1", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			_ = json.NewEncoder(w).Encode(map[string]any{
				"number": 1,
				"body":   "existing body\n<details id=\"pulumi-preview\"><summary>old summary</summary>\n\nold details\n\n</details>",
			})
		case http.MethodPatch:
			raw, err := io.ReadAll(r.Body)
			require.NoError(t, err)
			require.NoError(t, json.Unmarshal(raw, &patched))
			_ = json.NewEncoder(w).Encode(map[string]any{"number": 1})
		}
	})

	client := newTestClient(t, handler)

	err := setPRTrailerDetails(context.Background(), client, "braincorp/titanium", 1, "new summary", "new details", "pulumi-preview")
	require.NoError(t, err)

	body, ok := patched["body"].(string)
	require.True(t, ok)
	require.NotContains(t, body, "old details")
	require.Contains(t, body, "new details")
	require.Equal(t, 1, countSubstring(body, `<details id="pulumi-preview">`))
}

func countSubstring(s string, sub string) int {
	count := 0
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			count++
		}
	}
	return count
}
