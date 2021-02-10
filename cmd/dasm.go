package main

import (
	"os"

	"github.com/miguelff/8080/dasm"
)

func main() {
	dasm.Disassemble(os.Stdin, os.Stdout)
}
