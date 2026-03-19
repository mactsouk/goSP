package main

import (
    "fmt"
    "runtime/pprof"
)

func main() {
    for _, p := range pprof.Profiles() {
        fmt.Println(p.Name())
    }
}

