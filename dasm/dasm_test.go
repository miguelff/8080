package dasm

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/miguelff/8080/encoding"
)

func squish(asm string) string {
	chunks := strings.Split(asm, "\n")
	for i := range chunks {
		chunks[i] = strings.TrimSpace(chunks[i])
	}
	return strings.Join(chunks, "\n") + "\n"
}

func TestDisassemble(t *testing.T) {
	for _, tC := range []struct {
		desc string
		code string
		want string
	}{
		{
			"Space invaders head",
			"00 00 00 c3 d4 18 00 00 f5 c5 d5 e5 c3 8c 00",
			`NOP
				NOP
				NOP
				JMP $18D4
				NOP
				NOP
				PUSH PSW
				PUSH B
				PUSH D
				PUSH H
				JMP $008C`,
		},
	} {
		t.Run(tC.desc, func(t *testing.T) {
			w := strings.Builder{}
			err := Disassemble(bytes.NewReader(encoding.HexToBin(tC.code)), &w)
			if err != nil {
				t.Errorf("unexpected error when dissassembling binary: %v", err)
			}

			want := squish(tC.want)
			if got := w.String(); got != want {
				t.Errorf("got:\n%s\nwant:\n%s", got, want)
			}
		})
	}
}

func TestDisassemble_EndToEnd(t *testing.T) {
	bin, err := os.Open("../invaders/invaders.h")
	if err != nil {
		t.Fatal("cannot read object file")
	}
	asm, err := ioutil.ReadFile("../invaders/invaders.h.asm")
	if err != nil {
		t.Fatal("cannot read input assembly file")
	}

	w := strings.Builder{}
	err = Disassemble(bin, &w)
	if err != nil {
		t.Errorf("unexpected error when dissassembling binary: %v", err)
	}

	if got, want := w.String(), string(asm); got != want {
		t.Errorf("got %s \n\n want %s", got, want)
	}
}
