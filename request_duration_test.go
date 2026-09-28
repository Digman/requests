package requests

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRequestDurationTracksSendAndEnd(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		if flusher, ok := w.(http.Flusher); ok {
			flusher.Flush()
		}
		time.Sleep(20 * time.Millisecond)
		_, _ = io.WriteString(w, "ok")
	}))
	defer server.Close()

	r := NewClient("test").NewRequest().Get(server.URL).Send()
	if r.Duration != 0 {
		t.Fatalf("Send() duration=%v, want zero before End()", r.Duration)
	}

	_, body, err := r.End()
	if err != nil {
		t.Fatalf("End() error=%v", err)
	}
	if body != "ok" {
		t.Fatalf("body=%q, want ok", body)
	}
	if r.Duration <= 0 {
		t.Fatalf("complete duration=%v, want positive", r.Duration)
	}
}
