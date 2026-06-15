package main

import (
	"bufio"
	"fmt"
	"os"
)


func gameinit(name string) bool{
	return name != ""
}

func end(){
	os.Exit(0)
}


func main(){
	
	// hmmmm
	x := "Hello, World!"

	//Scanner main
	intake := bufio.NewScanner(os.Stdin)

	// Game state start with text
	fmt.Printf("Well here is the first line I type for G, as per usual: %s", x)
	fmt.Println("What is your name, challenger?")

	pname := intake.Text()

	if gameinit(pname){
		end()
	}
}