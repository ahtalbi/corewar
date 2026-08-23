package main

import (
	"fmt"
	"os"
	"corwar/globalvar"
)

func main() {

	Args := os.Args
	if len(Args) == globalvar.Zeroplayer{

	 fmt.Println("not player  to assmebler")
	 return
	}
	Exist := IsExistist(Args)
	if len(Exist) ==globalvar.Zeroplayer {

	fmt.Println("not valid  player  to assmebler")
	 return
	}




	

}
