package main

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

// Proxy represents the HTTP proxy server.
type Proxy struct {
	Addr     string
	Listener net.Listener
	Logger   *log.Logger
	wg       sync.WaitGroup
	cancel   context.CancelFunc
	ctx      context.Context
}

// NewProxy creates a new Proxy instance.
func NewProxy(ctx context.Context, addr string) (*Proxy, error) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	logger := log.New(os.Stdout, "[PROXY] ", log.LstdFlags|log.Lshortfile)
	ctx, cancel := context.WithCancel(ctx)

	return &Proxy{
		Addr:     addr,
		Listener: listener,
		Logger:   logger,
		ctx:      ctx,
		cancel:   cancel,
	}, nil
}

// Start begins listening for incoming connections.
func (p *Proxy) Start() {
	p.Logger.Printf("Proxy server starting on %s", p.Addr)

	for {
		conn, err := p.Listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				p.Logger.Println("Listener closed — stopping accept loop.")
				break
			}
			p.Logger.Printf("Error accepting connection: %v", err)
			continue
		}

		p.wg.Add(1)
		go p.handleClient(conn)
	}
}

// Shutdown gracefully closes the listener and waits for active connections to finish.
func (p *Proxy) Shutdown() {
	p.Logger.Println("Shutting down proxy server...")
	p.cancel()             // signal all goroutines to stop
	_ = p.Listener.Close() // unblock Accept
	p.wg.Wait()            // wait for in-flight connections
	p.Logger.Println("Proxy server shut down gracefully.")
}

// handleClient processes an incoming client connection.
func (p *Proxy) handleClient(conn net.Conn) {
	defer p.wg.Done()
	defer conn.Close()

	p.Logger.Printf("Accepted connection from %s", conn.RemoteAddr())

	reader := bufio.NewReader(conn)
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	req, err := http.ReadRequest(reader)
	if err != nil {
		if err != io.EOF {
			p.Logger.Printf("Error reading request from %s: %v", conn.RemoteAddr(), err)
		}
		return
	}
	conn.SetReadDeadline(time.Time{}) // clear deadline

	if req.Method == http.MethodConnect {
		p.handleHTTPSConnect(conn, req)
	} else {
		p.handleHTTPRequest(conn, req)
	}
}

// handleHTTPRequest proxies regular HTTP requests (GET, POST, etc.)
func (p *Proxy) handleHTTPRequest(clientConn net.Conn, req *http.Request) {
	req.RequestURI = ""
	req.Header.Del("Proxy-Connection")
	req.Header.Del("Proxy-Authenticate")
	req.Header.Del("Proxy-Authorization")
	req.Header.Del("TE")

	if req.URL.Host != "" {
		req.Host = req.URL.Host
	}

	destConn, err := net.DialTimeout("tcp", req.URL.Host, 10*time.Second)
	if err != nil {
		p.Logger.Printf("Error connecting to destination %s: %v", req.URL.Host, err)
		p.writeErrorResponse(clientConn, http.StatusGatewayTimeout, "Gateway Timeout")
		return
	}
	defer destConn.Close()

	if err := req.Write(destConn); err != nil {
		p.Logger.Printf("Error writing request to destination %s: %v", req.URL.Host, err)
		return
	}

	resp, err := http.ReadResponse(bufio.NewReader(destConn), req)
	if err != nil {
		p.Logger.Printf("Error reading response from %s: %v", req.URL.Host, err)
		p.writeErrorResponse(clientConn, http.StatusBadGateway, "Bad Gateway")
		return
	}
	defer resp.Body.Close()

	resp.Header.Del("Proxy-Authenticate")
	resp.Header.Del("Proxy-Connection")

	if err := resp.Write(clientConn); err != nil {
		p.Logger.Printf("Error writing response to client %s: %v", clientConn.RemoteAddr(), err)
	}
	p.Logger.Printf("Proxied %s %s → %d", req.Method, req.URL, resp.StatusCode)
}

// handleHTTPSConnect handles the CONNECT method (used for HTTPS tunneling)
func (p *Proxy) handleHTTPSConnect(clientConn net.Conn, req *http.Request) {
	destHost := req.URL.Host
	if !strings.Contains(destHost, ":") {
		destHost += ":443"
	}

	p.Logger.Printf("Handling HTTPS CONNECT for %s", destHost)

	destConn, err := net.DialTimeout("tcp", destHost, 10*time.Second)
	if err != nil {
		p.Logger.Printf("Error connecting to destination %s: %v", destHost, err)
		p.writeErrorResponse(clientConn, http.StatusServiceUnavailable, "Service Unavailable")
		return
	}
	defer destConn.Close()

	_, err = clientConn.Write([]byte("HTTP/1.1 200 Connection Established\r\n\r\n"))
	if err != nil {
		p.Logger.Printf("Error acknowledging CONNECT to %s: %v", destHost, err)
		return
	}

	ctx, cancel := context.WithCancel(p.ctx)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		io.Copy(destConn, clientConn)
		cancel()
	}()

	go func() {
		defer wg.Done()
		io.Copy(clientConn, destConn)
		cancel()
	}()

	<-ctx.Done()
	wg.Wait()
	p.Logger.Printf("Closed CONNECT tunnel for %s", destHost)
}

// writeErrorResponse sends a minimal HTTP error response.
func (p *Proxy) writeErrorResponse(conn net.Conn, statusCode int, message string) {
	statusText := http.StatusText(statusCode)
	body := fmt.Sprintf("%s\n", message)
	resp := fmt.Sprintf(
		"HTTP/1.1 %d %s\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\r\n%s",
		statusCode, statusText, len(body), body,
	)
	conn.Write([]byte(resp))
}

func main() {
	var port string
	flag.StringVar(&port, "port", "8080", "Port to run the HTTP proxy on")
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	proxy, err := NewProxy(ctx, ":"+port)
	if err != nil {
		log.Fatalf("Failed to create proxy: %v", err)
	}

	go proxy.Start()

	<-ctx.Done()
	proxy.Shutdown()
}
