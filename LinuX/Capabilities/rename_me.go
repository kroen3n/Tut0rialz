package main

import (
	"log"
	"os"
)

func main(){
	actualFile:= os.Args[1]
	newFile := os.Args[2]

	err := os.Rename(actualFile, newFile)

	if err != nil {
		log.Fatal(err)
	}
}
