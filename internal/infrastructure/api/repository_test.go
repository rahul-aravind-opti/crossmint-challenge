package api_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/crossmint/megaverse-challenge/internal/infrastructure/api"
	pkgretry "github.com/crossmint/megaverse-challenge/pkg/retry"
)

func newTestClient(baseURL string) *api.Client {
	return api.NewClient(api.ClientConfig{
		BaseURL:     baseURL,
		CandidateID: "test-id",
		Timeout:     time.Second,
		RetryConfig: pkgretry.Config{
			MaxAttempts:  1,
			InitialDelay: time.Millisecond,
			MaxDelay:     time.Millisecond,
			Multiplier:   1.0,
		},
		RequestsPerSecond: 100,
	})
}

func TestGetCurrentMapParsesContent(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/map/test-id", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(`{"map":{"content":[[null,{"type":0}],[{"type":1,"color":"red"},{"type":2,"direction":"left"}]]}}`))
		require.NoError(t, err)
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	t.Cleanup(server.Close)

	repo := api.NewRepository(newTestClient(server.URL))

	megaverse, err := repo.GetCurrentMap(context.Background())
	require.NoError(t, err)
	require.Equal(t, 2, megaverse.Height)
	require.Equal(t, 2, megaverse.Width)

	obj, err := megaverse.GetObject(0, 1)
	require.NoError(t, err)
	require.NotNil(t, obj)
	require.Equal(t, "POLYANET", obj.GetType())

	obj, err = megaverse.GetObject(1, 0)
	require.NoError(t, err)
	require.Equal(t, "SOLOON", obj.GetType())

	obj, err = megaverse.GetObject(1, 1)
	require.NoError(t, err)
	require.Equal(t, "COMETH", obj.GetType())
}

func TestGetCurrentMapHandlesEmptyContent(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(`{"map":{"content":[]}}`))
		require.NoError(t, err)
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	t.Cleanup(server.Close)

	repo := api.NewRepository(newTestClient(server.URL))

	megaverse, err := repo.GetCurrentMap(context.Background())
	require.NoError(t, err)
	require.Equal(t, 0, megaverse.Height)
	require.Equal(t, 0, megaverse.Width)
}
