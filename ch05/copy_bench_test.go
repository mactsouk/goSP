package filecopy

import (
	"io"
	"math/rand"
	"os"
	"testing"
	"time"
)

// --- File copy implementations ---

// Method 1: io.Copy
func copyFileIOCopy(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Sync()
}

// Method 2: manual buffer
func copyFileBuffer(src, dst string, bufSize int) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	buf := make([]byte, bufSize)
	for {
		n, errRead := in.Read(buf)
		if n > 0 {
			if _, errWrite := out.Write(buf[:n]); errWrite != nil {
				return errWrite
			}
		}
		if errRead == io.EOF {
			break
		}
		if errRead != nil {
			return errRead
		}
	}
	return out.Sync()
}

// Method 3: ReadFile/WriteFile
func copyFileReadWrite(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0644)
}

// --- Helpers ---

// makeTempFile creates a temp file of given size filled with random bytes.
func makeTempFile(size int64) (string, error) {
	f, err := os.CreateTemp("/tmp", "benchfile-*")
	if err != nil {
		return "", err
	}
	defer f.Close()

	rand.Seed(time.Now().UnixNano())
	buf := make([]byte, 8192)
	var written int64
	for written < size {
		rand.Read(buf)
		toWrite := int64(len(buf))
		if size-written < toWrite {
			toWrite = size - written
		}
		if _, err := f.Write(buf[:toWrite]); err != nil {
			return "", err
		}
		written += toWrite
	}
	return f.Name(), nil
}

// benchmark helper
func benchmarkCopy(b *testing.B, size int64, method func(string, string) error) {
	src, err := makeTempFile(size)
	if err != nil {
		b.Fatal(err)
	}
	defer os.Remove(src)

	dst := src + ".out"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := method(src, dst); err != nil {
			b.Fatal(err)
		}
		os.Remove(dst)
	}
}

// --- Benchmarks for different sizes ---
func BenchmarkIOCopy_1MB(b *testing.B) {
	benchmarkCopy(b, 1<<20, copyFileIOCopy)
}
func BenchmarkIOCopy_10MB(b *testing.B) {
	benchmarkCopy(b, 10<<20, copyFileIOCopy)
}
func BenchmarkBuffer_1MB(b *testing.B) {
	benchmarkCopy(b, 1<<20, func(s, d string) error {
		return copyFileBuffer(s, d, 4096)
	})
}
func BenchmarkBuffer_10MB(b *testing.B) {
	benchmarkCopy(b, 10<<20, func(s, d string) error {
		return copyFileBuffer(s, d, 4096)
	})
}
func BenchmarkReadWrite_1MB(b *testing.B) {
	benchmarkCopy(b, 1<<20, copyFileReadWrite)
}
func BenchmarkReadWrite_10MB(b *testing.B) {
	benchmarkCopy(b, 10<<20, copyFileReadWrite)
}
