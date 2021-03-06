package main

import (
	"aletheiaware.com/flamego"
	"aletheiaware.com/flamego/vm"
	"flag"
	"fmt"
	"log"
	"os"
)

var (
	memory  = flag.String("m", "", "The file to load into memory")
	storage = flag.String("s", "", "The file to load into storage")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options]\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	machine := vm.NewMachine()

	if *memory != "" {
		// Copy file into memory
		f, err := os.Open(*memory)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		machine.Memory.Load(f)
	}

	if *storage != "" {
		s := vm.NewFileStorage(machine.Memory, flamego.DeviceControlBlockAddress)
		if err := s.Open(*storage); err != nil {
			log.Fatal(err)
		}
		machine.Processor.AddDevice(s)
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
