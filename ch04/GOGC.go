package main

import (
	"fmt"
	"os"
	"runtime"
	"time"
)

func main() {
	fmt.Println("GOGC:", os.Getenv("GOGC"))
	var data [][]byte

	for i := 0; i < 15; i++ {
		chunk := make([]byte, 50<<20) // Allocate 50 MB chunks
		data = append(data, chunk)

		if (i+1)%3 == 0 && len(data) > 2 {
			data = data[len(data)-2:]
		}

		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		fmt.Printf("After iter %d: HeapAlloc = %d MB, NumGC = %d\n",
			i+1, m.HeapAlloc/1024/1024, m.NumGC)
		time.Sleep(500 * time.Millisecond) // give GC time to run

	}
}

