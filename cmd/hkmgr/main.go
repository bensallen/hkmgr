package main

import (
	"fmt"

	"github.com/bensallen/hkmgr/internal/root"
)

func main() {
	if err := root.Run(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
