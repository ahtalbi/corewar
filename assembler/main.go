package main

import (
	"fmt"
	"os"
	

	"corwar/globalvar"
)
func main() {

	if len(os.Args) < 2{


  fmt.Println("we nedd a least 2 player")

	}else if len(os.Args) > globalvar.MaxPlayers+1 {

		  fmt.Println("hhh just  4 player is the max size")

	}
	
}