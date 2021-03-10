package main

import (
	_ "embed"
	"flag"
	"fmt"
	"os"

	"github.com/miguelff/8080/emu"
)

func main() {
	//go:embed "invaders.rom"
	var rom []byte
	var err error

	debug := flag.Bool("d", false, "debug execution")
	flag.Parse()

	c := emu.Load(rom)

	for err == nil {
		err = c.Step(*debug)
	}

	fmt.Fprintf(os.Stderr, "%+v\n", err)
}
