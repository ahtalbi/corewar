package helpers

import (
	"fmt"
	"io"
	"log"
	"os"
)

func Readfile(path string) {

	magicbytezie := 4

	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Failed to open file: %s", err)
	}
	defer file.Close()

	magicBytes := make([]byte, magicbytezie)

	_, err = io.ReadFull(file, magicBytes)
	if err != nil {
		log.Fatalf("Failed to read magic bytes: %s", err)
	}

	fmt.Printf("Magic bytes (Hex): %x\n", magicBytes)
}
