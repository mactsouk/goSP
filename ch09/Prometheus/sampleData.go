// sampleData.go
//
// Populates Pushgateway with realistic CPU, memory, and HTTP metrics.
// Supports historical backfill for N days.
// Usage:
//   go run sampleData.go --push-url http://localhost:9091/metrics/job/testjob --backfill-days 7

package main

import (
	"bytes"
	"flag"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"time"
)

var pushURL string
var backfillDays int

func init() {
	flag.StringVar(&pushURL, "push-url", "http://localhost:9091/metrics/job/testjob", "Pushgateway URL")
	flag.IntVar(&backfillDays, "backfill-days", 7, "Number of days to backfill")
	rand.Seed(time.Now().UnixNano())
}

// Metric definitions
type cpuMetric struct {
	cpu   string
	mode  string
	value float64
}
type memMetric struct {
	mtype string
	value float64
}
type httpMetric struct {
	handler string
	method  string
	status  string
	value   float64
}

func main() {
	flag.Parse()
	fmt.Printf("Starting realistic simulation for %d days...\n", backfillDays)

	// Initialize metrics
	cpuMetrics := []cpuMetric{
		{"0", "idle", 0.2},
		{"0", "system", 0.05},
		{"1", "idle", 0.18},
		{"1", "user", 0.22},
	}
	memMetrics := []memMetric{
		{"used", 1.2e9},
		{"free", 9.8e8},
		{"cached", 5.5e8},
	}
	httpMetrics := []httpMetric{
		{"/api/v1/query", "GET", "200", 100},
		{"/api/v1/query", "GET", "500", 3},
	}

	// Backfill historical data
	backfillMetrics(cpuMetrics, memMetrics, httpMetrics)

	fmt.Println("Backfill complete. Exiting.")
}

// backfillMetrics simulates historical metrics over N days in discrete steps
func backfillMetrics(cpu []cpuMetric, mem []memMetric, httpM []httpMetric) {
	step := 5 * time.Minute
	totalSteps := backfillDays * 24 * int(time.Hour/step)
	count := 0

	start := time.Now().Add(-time.Duration(backfillDays*24) * time.Hour)
	for ts := start; count < totalSteps; ts = ts.Add(step) {
		body := generateMetrics(cpu, mem, httpM, ts)
		if err := pushMetrics(body); err != nil {
			fmt.Println("Push error:", err)
		}
		updateMetrics(cpu, mem, httpM, ts)

		count++
		if count%(totalSteps/20) == 0 { // every 5% progress
			fmt.Printf("Backfill progress: %3d%% (%d/%d steps)\n", count*100/totalSteps, count, totalSteps)
		}
	}
}

// generateMetrics builds Prometheus-compatible metrics (without timestamps)
func generateMetrics(cpu []cpuMetric, mem []memMetric, httpM []httpMetric, ts time.Time) []byte {
	var buf bytes.Buffer

	for _, m := range cpu {
		buf.WriteString(fmt.Sprintf(
			"node_cpu_seconds_total{cpu=\"%s\",mode=\"%s\",instance=\"localhost:9100\",job=\"node\"} %f\n",
			m.cpu, m.mode, m.value))
	}
	for _, m := range mem {
		buf.WriteString(fmt.Sprintf(
			"node_memory_bytes{type=\"%s\",instance=\"localhost:9100\",job=\"node\"} %f\n",
			m.mtype, m.value))
	}
	for _, m := range httpM {
		buf.WriteString(fmt.Sprintf(
			"http_requests_total{handler=\"%s\",method=\"%s\",status=\"%s\",instance=\"localhost:9090\",job=\"prometheus\"} %f\n",
			m.handler, m.method, m.status, m.value))
	}
	return buf.Bytes()
}

// pushMetrics sends metrics to Pushgateway
func pushMetrics(body []byte) error {
	req, err := http.NewRequest("POST", pushURL, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "text/plain")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode/100 != 2 {
		return fmt.Errorf("push failed, status %d", resp.StatusCode)
	}
	return nil
}

// updateMetrics applies realistic variations
func updateMetrics(cpu []cpuMetric, mem []memMetric, httpM []httpMetric, ts time.Time) {
	secOfDay := float64(ts.Hour()*3600 + ts.Minute()*60 + ts.Second())

	// CPU: day/night cycles + noise
	for i := range cpu {
		base := 0.05
		if cpu[i].mode == "idle" {
			base = 0.2
		}
		cpu[i].value = base + 0.05*math.Sin(2*math.Pi*secOfDay/86400) + rand.Float64()*0.01
		if cpu[i].value < 0 {
			cpu[i].value = 0
		}
	}

	// Memory: slow sinusoidal variation
	for i := range mem {
		mem[i].value *= 1 + 0.0001*(math.Sin(2*math.Pi*secOfDay/86400)+rand.Float64()*0.01)
		if mem[i].value < 0 {
			mem[i].value = 0
		}
	}

	// HTTP: random bursts
	for i := range httpM {
		base := 50.0
		if httpM[i].status == "200" {
			base = 100
		} else if httpM[i].status == "500" {
			base = 2
		}
		spike := 0.0
		if rand.Float64() < 0.05 {
			spike = rand.Float64() * 20
		}
		httpM[i].value = base + spike + rand.Float64()*5
	}
}
