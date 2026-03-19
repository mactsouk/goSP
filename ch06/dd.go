package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"
)

func humanBytes(bps float64) string {
	const unit = 1024
	if bps < unit {
		return fmt.Sprintf("%.1f B/s", bps)
	}
	div, exp := int64(unit), 0
	for n := bps / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB/s",
		bps/float64(div), "KMGTPE"[exp])
}

func reportProgress(start time.Time, blocksCopied, totalBytes int64) {
	elapsed := time.Since(start).Seconds()
	speed := humanBytes(float64(totalBytes) / elapsed)
	fmt.Fprintf(os.Stderr,
		"Progress: %d blocks copied (%d bytes, %s)\n",
		blocksCopied, totalBytes, speed)
}

func main() {
	fmt.Printf("Process ID: %d\n", os.Getpid())

	var (
		ifile string
		ofile string
		bs    int
		count int
		seek  int
	)

	flag.StringVar(&ifile, "if", "", "Input file (default stdin)")
	flag.StringVar(&ofile, "of", "", "Output file (default stdout)")
	flag.IntVar(&bs, "bs", 512, "Block size in bytes")
	flag.IntVar(&count, "count", 0, "Number of blocks to copy (0 means all)")
	flag.IntVar(&seek, "seek", 0, "Number of blocks to skip on output")
	flag.Parse()

	// Open input file (or stdin)
	var in *os.File
	var err error
	if ifile == "" {
		in = os.Stdin
	} else {
		in, err = os.Open(ifile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening input file: %v\n", err)
			os.Exit(1)
		}
		defer in.Close()
	}

	// Open output file (or stdout)
	var out *os.File
	if ofile == "" {
		out = os.Stdout
	} else {
		if _, err := os.Stat(ofile); err == nil {
			fmt.Fprintf(os.Stderr, "Output file %s already exists, not overwriting.\n", ofile)
			os.Exit(1)
		}
		out, err = os.OpenFile(ofile, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening output file: %v\n", err)
			os.Exit(1)
		}
		defer out.Close()
	}

	// Apply seek on output (if requested)
	if seek > 0 {
		offset := int64(seek) * int64(bs)
		_, err := out.Seek(offset, io.SeekStart)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Seek error: %v\n", err)
			os.Exit(1)
		}
	}

	buf := make([]byte, bs)
	var blocksCopied int64
	var totalBytes int64

	start := time.Now()
	done := make(chan struct{})

	// --- Signal handling ---
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGUSR1)

	go func() {
		for {
			sig := <-sigChan
			switch sig {
			case os.Interrupt, syscall.SIGTERM:
				reportProgress(start,
					atomic.LoadInt64(&blocksCopied),
					atomic.LoadInt64(&totalBytes))
				fmt.Fprintln(os.Stderr, "Interrupted, exiting.")
				os.Exit(1)
			case syscall.SIGUSR1:
				reportProgress(start,
					atomic.LoadInt64(&blocksCopied),
					atomic.LoadInt64(&totalBytes))
			}
		}
	}()

	// --- Copy loop ---
	for {
		if count > 0 && blocksCopied >= int64(count) {
			break
		}

		n, readErr := io.ReadFull(in, buf)
		if n > 0 {
			_, writeErr := out.Write(buf[:n])
			if writeErr != nil {
				fmt.Fprintf(os.Stderr, "Write error: %v\n", writeErr)
				os.Exit(1)
			}
			atomic.AddInt64(&blocksCopied, 1)
			atomic.AddInt64(&totalBytes, int64(n))
		}

		if readErr == io.EOF || readErr == io.ErrUnexpectedEOF {
			break
		}
		if readErr != nil {
			fmt.Fprintf(os.Stderr, "Read error: %v\n", readErr)
			os.Exit(1)
		}
	}

	close(done)
	reportProgress(start,
		atomic.LoadInt64(&blocksCopied),
		atomic.LoadInt64(&totalBytes))
	fmt.Fprintln(os.Stderr, "Done.")
}
