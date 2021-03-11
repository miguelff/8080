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

	debug := flag.String("d", "all", "debug opcode execution. Examples: '-d all' '-d \"C9 CD\"'")
	flag.Parse()

	c := emu.Load(rom)

	for err == nil {
		err = c.Step(emu.MakeDebugFilter(*debug))
	}

	fmt.Fprintf(os.Stderr, "%+v\n", err)
}
