package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"
)

type FileAnalysis struct {
	fileName  string
	wordCount uint64
	charCount uint64
}

type FileResult struct {
	analysis *FileAnalysis
	err      error
}

func main() {
	results := make(chan FileResult)
	files := []string{"boo.txt", "foo1.txt", "poo.txt"}
	for _, file := range files {
		go func(filename string) {
			a, e := readFile(filename)
			results <- FileResult{analysis: a, err: e}
		}(file)
	}

	var wordCount uint64 = 0
	var charCount uint64 = 0
	for i := 0; i < len(files); i++ {
		res := <-results
		if res.err != nil {
			fmt.Printf("%d. Ошибка при обработке файла %s: %v\n", i+1, files[i], res.err.Error())
		} else {
			fmt.Printf("%d. %s: %d слов, %d символов\n", i+1, res.analysis.fileName, res.analysis.wordCount, res.analysis.charCount)
			wordCount += res.analysis.wordCount
			charCount += res.analysis.charCount
		}
	}

	fmt.Printf("Итог: %d слов, %d символов", wordCount, charCount)

}

func readFile(p string) (*FileAnalysis, error) {
	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}

	res := FileAnalysis{
		fileName:  p,
		wordCount: 0,
		charCount: 0,
	}

	b := make([]byte, 4096) // 4KB buffer
	var tail string

	// reading file
	for {
		n, err := f.Read(b)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		if n == 0 {
			break
		}

		chunk := tail + string(b[:n])
		words := strings.Fields(chunk)
		lastRune, _ := utf8.DecodeLastRuneInString(chunk)
		if !isSpace(lastRune) {
			tail = words[len(words)-1]
			words = words[:len(words)-1]
		} else {
			tail = ""
		}

		res.wordCount += uint64(len(words))
		res.charCount += uint64(utf8.RuneCountInString(chunk))
	}

	if tail != "" {
		res.wordCount++ // don't forget about the last word!!!!
	}

	return &res, nil
}

func isSpace(r rune) bool {
	return unicode.IsSpace(r)
}
