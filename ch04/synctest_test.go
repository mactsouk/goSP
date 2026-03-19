package synctest_test

import (
	"bufio"
	"context"
	"io"
	"net"
	"net/http"
	"strings"
	"testing"
	"testing/synctest"
	"time"
)

// Example 1: Testing context.AfterFunc
// This demonstrates how synctest eliminates timing issues in concurrent tests
func TestAfterFunc(t *testing.T) {
	synctest.Run(func() {
		ctx, cancel := context.WithCancel(context.Background())
		funcCalled := false

		context.AfterFunc(ctx, func() {
			funcCalled = true
		})

		// Wait for all goroutines in the bubble to reach a blocking state
		synctest.Wait()

		// Assert the function hasn't been called before cancellation
		if funcCalled {
			t.Fatalf("AfterFunc function called before context is canceled")
		}

		cancel()

		// Wait again for the AfterFunc goroutine to execute
		synctest.Wait()

		// Assert the function was called after cancellation
		if !funcCalled {
			t.Fatalf("AfterFunc function not called after context is canceled")
		}
	})
}

// Example 2: Testing time-based functionality
// Shows how synctest uses fake time to make time-based tests fast and reliable
func TestWithTimeout(t *testing.T) {
	synctest.Run(func() {
		const timeout = 5 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		// Sleep for almost the entire timeout duration
		time.Sleep(timeout - time.Nanosecond)
		synctest.Wait()

		// Context should not be expired yet
		if err := ctx.Err(); err != nil {
			t.Fatalf("before timeout, ctx.Err() = %v; want nil", err)
		}

		// Sleep the remaining nanosecond to trigger timeout
		time.Sleep(time.Nanosecond)
		synctest.Wait()

		// Context should now be expired
		if err := ctx.Err(); err != context.DeadlineExceeded {
			t.Fatalf("after timeout, ctx.Err() = %v; want DeadlineExceeded", err)
		}
	})
}

// Example 3: Testing HTTP 100-Continue behavior
// Demonstrates testing network protocols with in-memory connections
func TestHTTP100Continue(t *testing.T) {
	synctest.Run(func() {
		// Create in-memory network connection using net.Pipe
		srvConn, cliConn := net.Pipe()
		defer srvConn.Close()
		defer cliConn.Close()

		// Configure HTTP transport to use our pipe connection
		tr := &http.Transport{
			DialContext: func(ctx context.Context, network, address string) (net.Conn, error) {
				return cliConn, nil
			},
			// Enable "Expect: 100-continue" handling
			ExpectContinueTimeout: 5 * time.Second,
		}

		body := "request body content"

		// Send HTTP request in a separate goroutine
		go func() {
			req, _ := http.NewRequest("PUT", "http://test.example/", strings.NewReader(body))
			req.Header.Set("Expect", "100-continue")
			resp, err := tr.RoundTrip(req)
			if err != nil {
				t.Errorf("RoundTrip: unexpected error %v", err)
			} else {
				resp.Body.Close()
			}
		}()

		// Read the request headers from server side
		req, err := http.ReadRequest(bufio.NewReader(srvConn))
		if err != nil {
			t.Fatalf("ReadRequest: %v", err)
		}

		// Start reading the body in another goroutine
		var gotBody strings.Builder
		go io.Copy(&gotBody, req.Body)

		// Wait for all goroutines to block
		synctest.Wait()

		// Verify client hasn't sent body yet (waiting for 100 Continue)
		if got := gotBody.String(); got != "" {
			t.Fatalf("before sending 100 Continue, unexpectedly read body: %q", got)
		}

		// Send 100 Continue response
		srvConn.Write([]byte("HTTP/1.1 100 Continue\r\n\r\n"))
		synctest.Wait()

		// Verify client now sends the body
		if got := gotBody.String(); got != body {
			t.Fatalf("after sending 100 Continue, read body %q, want %q", got, body)
		}

		// Complete the request with 200 OK
		srvConn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	})
}

// Example 4: Testing channel-based concurrent operations
func TestChannelCommunication(t *testing.T) {
	synctest.Run(func() {
		ch := make(chan string, 1)
		results := make(chan string)

		// Producer goroutine
		go func() {
			ch <- "hello"
			close(ch)
		}()

		// Consumer goroutine
		go func() {
			var collected []string
			for msg := range ch {
				collected = append(collected, msg)
			}
			results <- strings.Join(collected, ",")
		}()

		// Wait for all goroutines to complete their work
		synctest.Wait()

		// Verify the result
		select {
		case result := <-results:
			if result != "hello" {
				t.Fatalf("expected 'hello', got %q", result)
			}
		default:
			t.Fatal("no result received")
		}
	})
}

// Example 5: Demonstrating bubble isolation
func TestBubbleIsolation(t *testing.T) {
	// This channel is created outside the bubble
	externalCh := make(chan int)

	synctest.Run(func() {
		// This channel is created inside the bubble
		internalCh := make(chan string)

		go func() {
			// Operations on internal channels are durably blocking
			internalCh <- "bubble message"
		}()

		go func() {
			// Operations on external channels are NOT durably blocking
			// This would prevent proper bubble detection if it blocked
			select {
			case <-externalCh:
				// This won't happen in our test
			case msg := <-internalCh:
				if msg != "bubble message" {
					t.Errorf("expected 'bubble message', got %q", msg)
				}
			}
		}()

		synctest.Wait() // This waits for the channel communication to complete
	})
}
