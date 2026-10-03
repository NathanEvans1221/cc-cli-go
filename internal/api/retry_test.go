package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
)

func TestDoRequest_RetriesRetryableStatusAndReusesTransport(t *testing.T) {
	var hits atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if hits.Add(1) < 3 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}))
	defer srv.Close()

	client := NewClient("secret")
	if _, ok := client.httpClient.Transport.(*http.Transport); !ok {
		t.Fatal("client does not keep a reusable transport")
	}
	client.baseURL = srv.URL
	resp, err := client.doRequest(context.Background(), http.MethodPost, "/messages", NewRequest("m", 1))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK || hits.Load() != 3 {
		t.Fatalf("status %d hits %d", resp.StatusCode, hits.Load())
	}
}
