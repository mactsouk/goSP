package main

import (
	"bufio"
	"os"
	"testing"
)

func benchmarkBufferedWriteWithSize(b *testing.B, size int, bufSize int) {
	data := make([]byte, size)
	for i := 0; i < size; i++ {
		data[i] = byte('A' + (i % 26))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f, err := os.CreateTemp("", "buffered-var-*.txt")
		if err != nil {
			b.Fatal(err)
		}

		writer := bufio.NewWriterSize(f, bufSize)

		// write byte by byte but through the buffer
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

// Benchmarks with different buffer sizes for 1MB data
func BenchmarkBufferedWrite_1MB_512(b *testing.B) {
	benchmarkBufferedWriteWithSize(b, 1<<20, 512)
}
func BenchmarkBufferedWrite_1MB_4KB(b *testing.B) {
	benchmarkBufferedWriteWithSize(b, 1<<20, 4*1024)
}
func BenchmarkBufferedWrite_1MB_32KB(b *testing.B) {
	benchmarkBufferedWriteWithSize(b, 1<<20, 32*1024)
}
func BenchmarkBufferedWrite_1MB_256KB(b *testing.B) {
	benchmarkBufferedWriteWithSize(b, 1<<20, 256*1024)
}

// Benchmarks with different buffer sizes for 10MB data
func BenchmarkBufferedWrite_10MB_512(b *testing.B) {
	benchmarkBufferedWriteWithSize(b, 10<<20, 512)
}
func BenchmarkBufferedWrite_10MB_4KB(b *testing.B) {
	benchmarkBufferedWriteWithSize(b, 10<<20, 4*1024)
}
func BenchmarkBufferedWrite_10MB_32KB(b *testing.B) {
	benchmarkBufferedWriteWithSize(b, 10<<20, 32*1024)
}
func BenchmarkBufferedWrite_10MB_256KB(b *testing.B) {
	benchmarkBufferedWriteWithSize(b, 10<<20, 256*1024)
}
