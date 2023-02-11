package main // import "github.com/MartyHub/cac"

import (
	"os"

	"github.com/MartyHub/cac/internal"
)

func main() {
	if !internal.NewClient(internal.Parse(os.Args)).Run() {
		os.Exit(1)
	}
}
