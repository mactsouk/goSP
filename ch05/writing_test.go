package main

import (
	"bufio"
	"os"
	"testing"
)

func benchmarkUnbufferedWrite(b *testing.B, size int) {
	data := make([]byte, size)
	for i := 0; i < size; i++ {
		data[i] = byte('A' + (i % 26)) // simple content
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f, err := os.CreateTemp("", "unbuffered-*.txt")
		if err != nil {
			b.Fatal(err)
		}

		// unbuffered: write one byte at a time
		for j := 0; j < len(data); j++ {
			_, err := f.Write([]byte{data[j]})
			if err != nil {
				b.Fatal(err)
			}
		}
		f.Close()
		os.Remove(f.Name())
	}
}

func benchmarkBufferedWrite(b *testing.B, size int) {
	data := make([]byte, size)
	for i := 0; i < size; i++ {
		data[i] = byte('A' + (i % 26))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f, err := os.CreateTemp("", "buffered-*.txt")
		if err != nil {
			b.Fatal(err)
		}
		writer := bufio.NewWriter(f)

		// buffered: write through bufio.Writer
		for j := 0; j < len(data); j++ {
			err := writer.WriteByte(data[j])
			if err != nil {
				b.Fatal(err)
			}
		}

		writer.Flush()
		f.Close()
		os.Remove(f.Name())
	}
}

// Benchmarks for 1MB
func BenchmarkUnbufferedWrite_1MB(b *testing.B) {
	benchmarkUnbufferedWrite(b, 1<<20)
}
func BenchmarkBufferedWrite_1MB(b *testing.B) {
	benchmarkBufferedWrite(b, 1<<20)
}

// Benchmarks for 10MB
func BenchmarkUnbufferedWrite_10MB(b *testing.B) {
	benchmarkUnbufferedWrite(b, 10<<20)
}
func BenchmarkBufferedWrite_10MB(b *testing.B) {
	benchmarkBufferedWrite(b, 10<<20)
}
