package main

import (
	"fmt"

)

func main(){
	ar := []int{}
	for i :=0 ; i<100 ; i++ {
	 	ar = append(ar,i)
	}
	fmt.Println(ar)
}

