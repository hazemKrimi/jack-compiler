package main

import (
	"bufio"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/hazemKrimi/jack-compiler/internal/parser"
	"github.com/hazemKrimi/jack-compiler/internal/tokenizer"
)

func process(inputPath string) error {
	file, err := os.Open(inputPath)
	reader := bufio.NewReader(file)

	if err != nil {
		return err
	}

	tokens := make([]tokenizer.Token, 0, 1000)

	if err := tokenizer.ExtractTokens(&tokens, reader); err != nil {
		return err
	}

	parsed, err := parser.ParseTokens(tokens)
	
	if err != nil {
		return err
	}

	outputPath := strings.Replace(inputPath, ".jack", ".xml", 1)

	if err := os.WriteFile(outputPath, []byte(parsed), 0644); err != nil {
		return err
	}

	return nil
}

func walker(path string, entry fs.DirEntry, err error) error {
	if err != nil {
		return err
	}

	if match, _ := regexp.MatchString(".+\\.jack$", path); match {
		return process(path)
	}

	return nil
}

func main() {
	args := os.Args

	if len(args) != 2 {
		panic("You must provide a path for a Jack file or a directory that contains Jack files to compile!")
	}

	info, err := os.Stat(args[1])

	if err != nil {
		panic(err)
	}

	if info.IsDir() {
		err := filepath.WalkDir(args[1], walker)

		if err != nil {
			panic(err)
		}

		return
	}

	if match, _ := regexp.MatchString(".+\\.jack$", args[1]); !match {
		panic("You must provide a path for a Jack file or a directory that contains Jack files to compile!")
	}

	if err := process(args[1]); err != nil {
		panic(err)
	}
}
