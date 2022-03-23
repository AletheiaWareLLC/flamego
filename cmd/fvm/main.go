package main

import (
	"aletheiaware.com/flamego/vm"
	"io"
	"log"
	"os"
)

func main() {
	machine := vm.NewMachine()

	if len(os.Args) > 1 {
		// Copy binary into memory
		f, err := os.Open(os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		d, err := io.ReadAll(f)
		if err != nil {
			log.Fatal(err)
		}
		machine.Memory.Set(0, d)
	}

	// Signal the first context of the first core
	machine.Processor.Signal(0)

	// Run until processor halts
	var cycle int
	for ; !machine.Processor.HasHalted(); cycle++ {
		machine.Clock()
	}
	log.Println("Cycles:", cycle)
}
