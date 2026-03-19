// readProm.go
//
// Prometheus CLI tool: lists metrics or fetches all values for a metric+optional labels
// Usage:
//   List all series:
//     go run readProm.go --url http://localhost:9090 --list
//
//   Fetch values for a metric:
//     go run readProm.go --url http://localhost:9090 --metric node_cpu_seconds_total --labels cpu=0,mode=idle

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type InstantResponse struct {
	Status string `json:"status"`
	Data   struct {
		Result []struct {
			Value []interface{} `json:"value"`
		} `json:"result"`
	} `json:"data"`
	Error     string `json:"error"`
	ErrorType string `json:"errorType"`
}

type RangeResponse struct {
	Status string `json:"status"`
	Data   struct {
		Result []struct {
			Values [][]interface{} `json:"values"`
		} `json:"result"`
	} `json:"data"`
	Error     string `json:"error"`
	ErrorType string `json:"errorType"`
}

func main() {
	baseURL := flag.String("url", "http://localhost:9090", "Prometheus base URL")
	list := flag.Bool("list", false, "List all metrics")
	metric := flag.String("metric", "", "Metric name to fetch")
	labels := flag.String("labels", "", "Comma-separated labels to filter, e.g., cpu=0,mode=idle")
	startStr := flag.String("start", "", "Start time (RFC3339) or empty for 7 days ago")
	endStr := flag.String("end", "", "End time (RFC3339) or empty for now")
	step := flag.String("step", "60s", "Step duration, e.g., 15s, 1m")
	flag.Parse()

	if *list {
		listMetrics(*baseURL)
		return
	}

	if *metric == "" {
		fmt.Fprintln(os.Stderr, "Error: must specify either --list or --metric")
		os.Exit(1)
	}

	labelSelector := parseLabels(*labels)
	start, end := parseTimeRange(*startStr, *endStr)
	fetchMetricRange(*baseURL, *metric, labelSelector, start, end, *step)
}

// parseLabels converts "key1=val1,key2=val2" into Prometheus selector "{key1="val1",key2="val2"}"
func parseLabels(s string) string {
	if s == "" {
		return ""
	}
	parts := strings.Split(s, ",")
	var sel []string
	for _, p := range parts {
		kv := strings.SplitN(p, "=", 2)
		if len(kv) != 2 {
			continue
		}
		key := strings.TrimSpace(kv[0])
		val := strings.TrimSpace(kv[1])
		sel = append(sel, fmt.Sprintf(`%s="%s"`, key, val))
	}
	return "{" + strings.Join(sel, ",") + "}"
}

// listMetrics fetches all metric names from Prometheus
func listMetrics(baseURL string) {
	api := fmt.Sprintf("%s/api/v1/label/__name__/values", baseURL)
	resp, err := http.Get(api)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error querying metric names: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result struct {
		Status string   `json:"status"`
		Data   []string `json:"data"`
	}
	if err := json.Unmarshal(body, &result); err != nil || result.Status != "success" {
		fmt.Fprintln(os.Stderr, "Invalid Prometheus response.")
		os.Exit(1)
	}

	sort.Strings(result.Data)
	fmt.Println("Metrics:")
	for _, name := range result.Data {
		fmt.Println(" -", name)
	}
}

// parseTimeRange parses optional start/end timestamps
func parseTimeRange(startStr, endStr string) (time.Time, time.Time) {
	var end time.Time
	if endStr == "" {
		end = time.Now()
	} else {
		t, err := time.Parse(time.RFC3339, endStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid end time: %v\n", err)
			os.Exit(1)
		}
		end = t
	}

	var start time.Time
	if startStr == "" {
		start = end.Add(-7 * 24 * time.Hour)
	} else {
		t, err := time.Parse(time.RFC3339, startStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid start time: %v\n", err)
			os.Exit(1)
		}
		start = t
	}
	return start, end
}

// fetchMetricRange fetches metric values in the specified range
func fetchMetricRange(baseURL, metric, labels string, start, end time.Time, step string) {
	query := metric + labels
	// Use Unix timestamps instead of RFC3339 for better compatibility
	startTS := fmt.Sprintf("%d", start.Unix())
	endTS := fmt.Sprintf("%d", end.Unix())
	api := fmt.Sprintf("%s/api/v1/query_range?query=%s&start=%s&end=%s&step=%s",
		baseURL, url.QueryEscape(query), startTS, endTS, step)

	resp, err := http.Get(api)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching range: %v\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var rr RangeResponse
	if err := json.Unmarshal(body, &rr); err != nil {
		fmt.Fprintf(os.Stderr, "Invalid JSON response: %v\nResponse body: %s\n", err, string(body))
		os.Exit(1)
	}

	if rr.Status != "success" {
		fmt.Fprintf(os.Stderr, "Prometheus error: %s (%s)\n", rr.Error, rr.ErrorType)
		fmt.Fprintf(os.Stderr, "Query was: %s\n", query)
		os.Exit(1)
	}

	if len(rr.Data.Result) == 0 {
		fmt.Println("No samples found.")
		return
	}

	for _, series := range rr.Data.Result {
		for _, pair := range series.Values {
			if len(pair) != 2 {
				continue
			}
			ts, val := parsePair(pair)
			fmt.Printf("%s  %s\n", ts.Format(time.RFC3339), val)
		}
	}
}

// parsePair converts a [timestamp, value] pair to proper types
func parsePair(pair []interface{}) (time.Time, string) {
	tsFloat := 0.0
	switch v := pair[0].(type) {
	case float64:
		tsFloat = v
	case string:
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			tsFloat = 0
		} else {
			tsFloat = f
		}
	}
	sec := int64(tsFloat)
	nsec := int64((tsFloat - float64(sec)) * 1e9)
	ts := time.Unix(sec, nsec)
	valStr := fmt.Sprintf("%v", pair[1])
	return ts, valStr
}
