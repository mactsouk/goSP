package main

import (
	"encoding/json"
	"testing"
)

var cfg = ServerConfig{
	Name: "api",
	Port: 8080,
	Tags: []string{"go", "prod"},
	Limits: map[string]uint32{
		"requests_per_min": 1200,
		"conn":             100,
	},
	TLS: &TLSConfig{
		Enabled: true,
		CertPEM: "-----BEGIN CERTIFICATE-----...",
	},
	Metadata: map[string]any{
		"region":   "eu-central-1",
		"revision": 42,
		"debug":    false,
	},
}

// Benchmark: Standard library JSON encoder (compact)
func BenchmarkStdJSONCompact(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = json.Marshal(cfg)
	}
}

// Benchmark: Custom reflection-based encoder (compact)
func BenchmarkCustomJSONCompact(b *testing.B) {
	enc := &Encoder{} // no indentation, compact mode
	for i := 0; i < b.N; i++ {
		_ = enc.Marshal(cfg)
	}
}

// Benchmark: Standard library JSON encoder (pretty)
func BenchmarkStdJSONPretty(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = json.MarshalIndent(cfg, "", "  ")
	}
}

// Benchmark: Custom reflection-based encoder (pretty)
func BenchmarkCustomJSONPretty(b *testing.B) {
	enc := &Encoder{Indent: "  ", SortMapKeys: true}
	for i := 0; i < b.N; i++ {
		_ = enc.Marshal(cfg)
	}
}
