package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	var c bool
	var l bool
	var w bool
	var m bool

	flag.BoolVar(&c, "c", false, "outputs the number of bytes in a file")
	flag.BoolVar(&l, "l", false, "outputs the number of lines in a file")
	flag.BoolVar(&w, "w", false, "outputs the number of words in a file")
	flag.BoolVar(&m, "m", false, "outputs the number of multi(bytes) in a file")
	flag.Parse()

	if !c && !l && !w && !m {
		c, l, w = true, true, true
	}

	var file []byte
	var fileName string
	var err error

	if fp := flag.Arg(0); fp != "" {
		file, err = os.ReadFile(fp)
		if err != nil {
			reportErr(fmt.Sprintf("failed to read file: %+v", err))
		}

		fps := strings.Split(fp, "/")
		fileName = fps[len(fps) - 1]
	} else {
		stat, err := os.Stdin.Stat()
		if err != nil {
			reportErr(fmt.Sprintf("failed to read file: %+v", err))
		}

		if (stat.Mode() & os.ModeNamedPipe) != 0 {
			file, err = io.ReadAll(os.Stdin)
			if err != nil {
				reportErr(fmt.Sprintf("failed to read file: %+v", err))
			}
		}
	}

	var outputs string

	if l {
		scanner := bufio.NewScanner(bytes.NewBuffer(file))
		lines := 0
		for scanner.Scan() {
			lines++
		}
		outputs += fmt.Sprintf("    %d", lines)
	}

	if w {
		outputs += fmt.Sprintf("   %d", len(bytes.Fields(file)))
	}

	if c || m {
		if m && supportMultibyte() {
			outputs += fmt.Sprintf("  %d", len(bytes.Runes(file)))
		} else {
			outputs += fmt.Sprintf("  %d", len(file))
		}
	}

	fmt.Printf("%s %s\n", outputs, fileName)
}



func getLocale() string {
	if lcAll := os.Getenv("LC_ALL"); lcAll != "" {
		return lcAll
	}
	if lcCType := os.Getenv("LC_CTYPE"); lcCType != "" {
		return lcCType
	}

	return os.Getenv("LANG")
}

func supportMultibyte() bool {
	locale := getLocale()
	return strings.Contains(strings.ToLower(locale), "utf-8") || strings.Contains(getLocale(), "utf8")
}

func reportErr(msg string) {
	fmt.Fprintf(os.Stderr, msg)
	os.Exit(1)
}
