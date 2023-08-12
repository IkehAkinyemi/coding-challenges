package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	var c bool
	flag.BoolVar(&c, "c", false, "outputs the number of bytes in a file")
	flag.Parse()

	fp := flag.Arg(0)
	if fp == "" {
		log.Fatalf("add a file path to read")
	}

	file, err := os.ReadFile(fp)
	if err != nil {
		log.Fatalf("failed to read file: %+v", err)
	}

	stat, err := os.Stat(fp)
	if err != nil {
		log.Fatalf("failed to check file stats: %+v", err)
	}

	if c {
		count := countBytes(file)
		fmt.Printf("%d %s\n", count, stat.Name())
		return
	} else {
		log.Fatalf("specify a cwc flag action")
	}
}

func countBytes(file []byte) (count int) {
	return len(file)
}