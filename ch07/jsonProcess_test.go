package main

import (
	"bytes"
	"encoding/json"
	"io"
	"testing"
)

type Record struct {
	ID      int     `json:"id"`
	Name    string  `json:"name"`
	Email   string  `json:"email"`
	Balance float64 `json:"balance"`
}

func generateRecords(n int) []Record {
	records := make([]Record, n)
	for i := 0; i < n; i++ {
		records[i] = Record{
			ID:      i + 1,
			Name:    "User" + string(rune('A'+(i%26))),
			Email:   "user@example.com",
			Balance: float64(i) * 10.5,
		}
	}
	return records
}

func BenchmarkSingleRecord(b *testing.B) {
	records := generateRecords(1000) // same dataset size as stream version
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, rec := range records {
			data, _ := json.Marshal(rec)
			var decoded Record
			_ = json.Unmarshal(data, &decoded)
		}
	}
}

func BenchmarkStream(b *testing.B) {
	records := generateRecords(1000)

	// Pre-encode all records into a JSON stream once
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	for _, r := range records {
		_ = enc.Encode(r)
	}
	streamData := buf.Bytes()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dec := json.NewDecoder(bytes.NewReader(streamData))
		var rec Record
		for {
			if err := dec.Decode(&rec); err == io.EOF {
				break
			} else if err != nil {
				b.Fatal(err)
			}
		}
	}
}
