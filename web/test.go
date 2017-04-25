package main

import (
	"path/filepath"
	"fmt"
)

func main(){
	 files,err := filepath.Glob(".")
	if err != nil{
		fmt.Println(err)
	}
	for _,f := range files{
		fmt.Println(f)
	}
}

