package main

import (
	"aletheiaware.com/flamego/assembler"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

var (
	output  = flag.String("o", "", "Output file")
	address = flag.String("a", "", "Address file")
	verbose = flag.Bool("v", false, "Log additional information")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options] input\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
	args := flag.Args()

	if len(args) == 0 {
		fmt.Println("Missing input file")
		flag.Usage()
		return
	}

	a := assembler.NewAssembler()

	for _, i := range args {
		f, err := os.Open(i)
		if err != nil {
			log.Fatal(err)
		}
		if _, err := a.ReadFrom(f); err != nil {
			log.Fatal(err)
		}
	}

	writer := io.Discard
	if *output != "" {
		f, err := os.Create(*output)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		writer = f
	}

	count, err := a.WriteTo(writer)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Wrote", count, "bytes")

	if *address != "" {
		f, err := os.Create(*address)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		for _, a := range a.Addressables() {
			f.WriteString(fmt.Sprintf("0x%016x", a.AbsoluteAddress()))
			if s, ok := a.(fmt.Stringer); ok {
				f.WriteString(": ")
				f.WriteString(s.String())
			}
			f.WriteString("\n")
		}
	}

	if *verbose {
		for n, c := range a.Constants() {
			log.Println(n, c)
		}
		for _, l := range a.Labels() {
			log.Println(l, "address", l.AbsoluteAddress())
		}
	}
}
