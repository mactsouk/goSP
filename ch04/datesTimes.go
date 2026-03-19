package main

import (
	"fmt"
	"time"
)

func main() {
	// 1. Get the current local time
	now := time.Now()
	fmt.Println("Current local time:", now)

	// 2. Format the current time in different layouts
	fmt.Println("Default format:", now.Format(time.RFC3339))
	fmt.Println("Custom format:", now.Format("Monday, 02-Jan-2006 15:04:05 MST"))

	// 3. Parse a time string into a time.Time value
	layout := "2006-01-02 15:04:05"
	input := "2025-09-03 14:30:00"
	parsedTime, err := time.Parse(layout, input)
	if err != nil {
		fmt.Println("Error parsing time:", err)
		return
	}
	fmt.Println("Parsed time:", parsedTime)

	// 4. Add and subtract durations
	twoHoursLater := now.Add(2 * time.Hour)
	thirtyMinutesEarlier := now.Add(-30 * time.Minute)
	fmt.Println("Two hours later:", twoHoursLater)
	fmt.Println("Thirty minutes earlier:", thirtyMinutesEarlier)

	// 5. Compare times
	if parsedTime.After(now) {
		fmt.Println("Parsed time is in the future.")
	} else {
		fmt.Println("Parsed time is in the past.")
	}

	// 6. Work with time zones
	// Convert local time to UTC
	utc := now.UTC()
	fmt.Println("UTC time:", utc)

	// Load a specific time zone
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		fmt.Println("Error loading location:", err)
		return
	}
	nyTime := now.In(loc)
	fmt.Println("Time in New York:", nyTime)

	// 7. Calculate duration between two times
	duration := parsedTime.Sub(now)
	fmt.Printf("Duration between now and parsed time: %v\n", duration)
}
