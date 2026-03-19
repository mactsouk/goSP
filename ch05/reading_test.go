package main

import (
	"bufio"
	"io"
	"os"
	"testing"
)

// Unbuffered reading: 1-byte reads
func unbufferedRead(file *os.File) (int, error) {
	buf := make([]byte, 1)
	total := 0
	for {
		n, err := file.Read(buf)
		if n > 0 {
			total += n
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return total, err
		}
	}
	return total, nil
}

// Buffered reading: uses bufio.Reader
func bufferedRead(file *os.File) (int, error) {
	reader := bufio.NewReader(file)
	total := 0
	for {
		_, err := reader.ReadByte()
		if err == io.EOF {
			break
		}
		if err != nil {
			return total, err
		}
		total++
	}
	return total, nil
}

func BenchmarkUnbufferedRead(b *testing.B) {
	filename := "/tmp/largefile.txt"

	// Ensure file exists (create once if not present)
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		f, err := os.Create(filename)
		if err != nil {
			b.Fatalf("cannot create file: %v", err)
		}
		data := make([]byte, 10<<20) // 10 MB
		for i := range data {
			data[i] = byte('A' + (i % 26))
		}
		_, err = f.Write(data)
		f.Close()
		if err != nil {
			b.Fatalf("cannot write file: %v", err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f, err := os.Open(filename)
		if err != nil {
			b.Fatalf("cannot open file: %v", err)
		}
		_, err = unbufferedRead(f)
		f.Close()
		if err != nil {
			b.Fatalf("unbuffered read failed: %v", err)
		}
	}
}

func BenchmarkBufferedRead(b *testing.B) {
	filename := "/tmp/largefile.txt"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f, err := os.Open(filename)
		if err != nil {
			b.Fatalf("cannot open file: %v", err)
		}
		_, err = bufferedRead(f)
		f.Close()
		if err != nil {
			b.Fatalf("buffered read failed: %v", err)
		}
	}
}
