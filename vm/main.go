package main

import (
	"corwar/vm/helpers"
	"corwar/vm/loader"
	"fmt"
	"log"
	"os"
)

func main() {
	args := os.Args[1:]
	if len(args) < 2{
		log.Fatal("Please provdide at least 2 players")
	}

	paths := helpers.Pathsolver(args)
	fmt.Println(paths)
	if len(paths)<2{
		log.Fatal("Please provdide at least 2 valid  players")

	}
    for _,path := range paths{
		loader.Load(path)
	}


}
