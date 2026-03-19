package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

// Method 1: Using io.Copy (simple and efficient)
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

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Sync()
}

// Method 2: Using a manual buffer
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
		n, readErr := in.Read(buf)
		if n > 0 {
			if _, writeErr := out.Write(buf[:n]); writeErr != nil {
				return writeErr
			}
		}
		if readErr == io.EOF {
			break
		}
		if readErr != nil {
			return readErr
		}
	}

	return out.Sync()
}

// Method 3: Using os.ReadFile and os.WriteFile
func copyFileReadWrite(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0644)
}

// Method 4: Using io.CopyBuffer
func copyFileCopyBuffer(src, dst string, bufSize int) error {
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
	_, err = io.CopyBuffer(out, in, buf)
	if err != nil {
		return err
	}
	return out.Sync()
}

func main() {
	method := flag.String("method", "iocopy", "Copy method: iocopy | buffer | readwrite | copybuffer")
	bufSize := flag.Int("bufsize", 4096, "Buffer size (for buffer and copybuffer methods)")
	flag.Parse()

	if flag.NArg() != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <source> <dest>\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	src := flag.Arg(0)
	dst := flag.Arg(1)

	var err error
	switch *method {
	case "iocopy":
		err = copyFileIOCopy(src, dst)
	case "buffer":
		err = copyFileBuffer(src, dst, *bufSize)
	case "readwrite":
		err = copyFileReadWrite(src, dst)
	case "copybuffer":
		err = copyFileCopyBuffer(src, dst, *bufSize)
	default:
		fmt.Fprintf(os.Stderr, "Unknown method: %s\n", *method)
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error copying file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("File copied successfully using method '%s' → %s\n", *method, dst)
}
