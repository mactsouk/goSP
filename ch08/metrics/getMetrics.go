package main

import (
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/shirou/gopsutil/net"
)

// TCP metrics
var (
	tcpEstablished = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "tcp_established_connections",
		Help: "Number of TCP connections in ESTABLISHED state",
	})
	tcpListen = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "tcp_listen_connections",
		Help: "Number of TCP sockets in LISTEN state",
	})
	tcpOther = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "tcp_other_connections",
		Help: "Number of TCP connections in other states",
	})
)

// UDP metrics
var (
	udpListening = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "udp_listen_sockets",
		Help: "Number of UDP sockets currently listening",
	})
)

func init() {
	prometheus.MustRegister(tcpEstablished, tcpListen, tcpOther, udpListening)
}

// updateTCPMetrics collects TCP connection statistics
func updateTCPMetrics() {
	conns, err := net.Connections("tcp")
	if err != nil {
		log.Printf("Error reading TCP connections: %v", err)
		return
	}

	var est, listen, other float64
	for _, c := range conns {
		switch c.Status {
		case "ESTABLISHED":
			est++
		case "LISTEN":
			listen++
		default:
			other++
		}
	}

	tcpEstablished.Set(est)
	tcpListen.Set(listen)
	tcpOther.Set(other)
	log.Printf("TCP metrics updated: EST=%v LISTEN=%v OTHER=%v\n", est, listen, other)
}

// updateUDPMetrics collects UDP socket statistics
func updateUDPMetrics() {
	conns, err := net.Connections("udp")
	if err != nil {
		log.Printf("Error reading UDP connections: %v", err)
		return
	}

	// UDP doesn't have a "state", so we just count active sockets
	udpListening.Set(float64(len(conns)))
	log.Printf("UDP metrics updated: LISTEN=%v\n", len(conns))
}

func main() {
	// Periodically update metrics
	go func() {
		for {
			updateTCPMetrics()
			updateUDPMetrics()
			time.Sleep(5 * time.Second)
		}
	}()

	// Expose metrics for Prometheus
	addr := "0.0.0.0:9100"
	log.Printf("Metrics server listening on %s\n", addr)
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(addr, nil))
}

