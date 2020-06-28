package main

import (
	"github.com/ktugan/loci"
)

func main() {
	config := loci.LoadConfig()

	err := loci.PrepConfig(&config)
	if err != nil {
		panic(err)
	}

	loci.Loci(config)
}
