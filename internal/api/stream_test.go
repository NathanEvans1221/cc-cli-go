package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestStream_ReadsSSE(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/messages" || r.Header.Get("x-api-key") != "secret" {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		w.Header().Set("content-type", "text/event-stream")
		_, _ = w.Write([]byte("not-data\n\ndata: {\"type\":\"message_stop\"}\n\ndata: [DONE]\n"))
	}))
	defer srv.Close()

	client := NewClient("secret")
	client.baseURL = srv.URL
	events, err := client.Stream(context.Background(), NewRequest("model", 8))
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for event := range events {
		got = append(got, event.Type)
	}
	if len(got) != 1 || got[0] != "message_stop" {
		t.Fatalf("events %v", got)
	}
}

func TestStream_UnexpectedStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "nope", http.StatusUnauthorized)
	}))
	defer srv.Close()

	client := NewClient("secret")
	client.baseURL = srv.URL
	if _, err := client.Stream(context.Background(), NewRequest("model", 8)); err == nil {
		t.Fatal("expected status error")
	}
}

func TestDoRequest_Failures(t *testing.T) {
	client := NewClient("secret")
	client.baseURL = "http://127.0.0.1:1"
	client.httpClient.Timeout = 200 * time.Millisecond
	if _, err := client.doRequest(context.Background(), http.MethodPost, "/messages", NewRequest("m", 1)); err == nil {
		t.Fatal("expected connection error")
	}
	if _, err := client.doRequest(context.Background(), http.MethodPost, "/messages", make(chan int)); err == nil {
		t.Fatal("expected marshal error")
	}
	if _, _, err := ParseDelta([]byte("{")); err == nil {
		t.Fatal("expected parse error")
	}
	if _, err := ParseContentBlock([]byte("{")); err == nil {
		t.Fatal("expected block parse error")
	}
	if _, err := ParseMessageStart([]byte("{")); err == nil {
		t.Fatal("expected message parse error")
	}
}
