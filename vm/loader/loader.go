package loader


import (

	"corwar/vm/helpers"
)
type Player struct{
Name  string
wasf string
code []byte

}


func Load(path string)(Player,error){
helpers.Readfile(path)

return Player{},nil
}



