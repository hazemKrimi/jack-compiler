package main

import (
	"fmt"
	"os"
	"regexp"
)

func process(Path string) {
	fmt.Println(Path)
}

func main() {
	args := os.Args

	if len(args) != 2 {
		panic("You must provide a path for a Jack file or a directory that contains Jack files to compile!")
	}

	Path := args[1]

	if match, _ := regexp.MatchString(".+\\.jack", Path); match {
		process(Path)
	}
}
