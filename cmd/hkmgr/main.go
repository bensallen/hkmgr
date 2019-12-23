package main

import (
	"fmt"

	"github.com/bensallen/hkmgr/internal/root"
)

var version = "unknown"

func main() {
	root.Version = version
	if err := root.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
