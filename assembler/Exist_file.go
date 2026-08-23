package main


import ("os"
"fmt")


func IsExistist(paths []string)[]string {

	paths  = []string{}


	for _,player :=range paths {
		_,err := os.Stat(player)
		if err!=nil {

			if os.IsNotExist(err) {
            fmt.Println(player,"not found dont play with me")
        }
		}
		paths = append(paths, player)




	}
 return  paths

}
