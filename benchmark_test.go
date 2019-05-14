package main

import (
	"bytes"
	"net/http"
	"testing"
)

func BenchmarkTokeniseEndpoint(b *testing.B) {
	buffer := bytes.NewBuffer([]byte(`{"Body":"123"}`))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = http.Post("http://0.0.0.0:8080/v1/tokenise", "application/json", buffer)
	}
}

func BenchmarkHealthEndpoint(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = http.Get("http://0.0.0.0:8080/health")
	}
}
