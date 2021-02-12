package main

import (
	_ "embed"
	"fmt"
	"os"

	"github.com/miguelff/8080/emu"
)

func main() {
	//go:embed "invaders.rom"
	var rom []byte
	var err error

	c := emu.Computer{}
	c.Load(rom)

	for err == nil {
		err = c.Step()
	}

	fmt.Fprintf(os.Stderr, "%+v\n", err)
}
