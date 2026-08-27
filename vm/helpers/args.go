package helpers

import (
	"corwar/globalvar"
	"fmt"
	"os"
)

func Pathsolver(args []string) []string {
	binarypath := globalvar.Binarypath
	result := []string{}
	errorpath := ""
	for i := 0; i < len(args); i++ {
	
		_, err := os.Stat(binarypath + args[i])
		if err != nil {
		
			errorpath += args[i]+" "
		}

		result = append(result, binarypath+args[i])
	}

	if len(errorpath) != 0 {

		fmt.Println("cannot find binary file for those names "+errorpath)

	}

	return result

}
