package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ./asm file.s [file2.s ...]")
		return
	}
	ok := true
	for _, file := range os.Args[1:] {
		if err := Assemble(file); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			ok = false
		}
	}
	if !ok {
		os.Exit(1)
	}
}
