package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
)

type cwc struct {
	fileName string
	words int
	lines int
	characters int
}

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

	filePaths := flag.Args()
	results := []cwc{}

	if len(filePaths) == 0 {
		res, err := processStdin(m)
		if err != nil {
			reportErr(err.Error())
		}
		results = append(results, res)
	} else {
		for _, fp := range filePaths {
			res, err := processFile(fp, m)
			if err != nil {
				reportErr(err.Error())
			}
			results = append(results, res)
		}
	}

	displayResult(results, c, w, l, m)
}

func cwcReader(file io.Reader, mFlag bool) (cwc, error) {
	bufReader := bufio.NewReader(file)
	res := cwc{}
	isWord := false

	for {
		r, size, err := bufReader.ReadRune()
		if err != nil  && err != io.EOF {
			return cwc{}, err
		}

		if err == io.EOF {
			break
		}

		if mFlag && supportMultibyte() {
			res.characters++
		} else {
			res.characters += size
		}

		if r == '\n' {
			res.lines++
		}

		if unicode.IsSpace(r) {
			isWord = false
		} else if !isWord {
			res.words++
			isWord = true
		}
	}

	return res, nil
}

func processStdin(mFlag bool) (cwc, error) {
	res, err := cwcReader(os.Stdin, mFlag)
	if err != nil {
		return cwc{}, err
	}
	
	res.fileName = fileName(os.Stdin)

	return res, nil
}

func processFile(fpath string, mFlag bool) (cwc, error) {
	file, err := os.Open(fpath)
	if err != nil {
		return cwc{}, err
	}
	defer file.Close()

	res, err := cwcReader(file, mFlag)
	if err != nil {
		return cwc{}, err
	}

	res.fileName = fileName(file)

	return res, nil
}

func displayResult(results []cwc, cFlag, wFlag, lFlag, mFlag bool) {
	for _, res := range results {
		var outputs string
		if lFlag {
			outputs += fmt.Sprintf("    %d", res.lines)
		}
		if wFlag {
			outputs += fmt.Sprintf("   %d", res.words)
		}
		if cFlag || mFlag {
			outputs += fmt.Sprintf("  %d", res.characters)
		}

		fmt.Printf("%s %s\n", outputs, res.fileName)
	}
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

func fileName(file *os.File) string {
	info, err := file.Stat()
	if err != nil {
		reportErr(err.Error())
	}

	return info.Name()
}