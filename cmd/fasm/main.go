package main

import (
	"aletheiaware.com/flamego/assembler"
	"flag"
	"log"
	"os"
)

var output = flag.String("o", "", "Output file")

func main() {
	flag.Parse()

	a := assembler.NewAssembler()

	for _, i := range flag.Args() {
		f, err := os.Open(i)
		if err != nil {
			log.Fatal(err)
		}
		if _, err := a.ReadFrom(f); err != nil {
			log.Fatal(err)
		}
	}

	writer := os.Stdout
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
}
