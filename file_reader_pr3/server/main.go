package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"
)

const uploadsDir = "uploads/"

type FileAnalysis struct {
	FileName  string `json:"file_name"`
	WordCount uint64 `json:"word_count"`
	CharCount uint64 `json:"char_count"`
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	// file size limitation 10 MB
	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("uploadfile")
	if err != nil {
		http.Error(w, "file error", http.StatusBadRequest)
		return
	}
	defer file.Close()

	path := uploadsDir + handler.Filename

	f, err := os.Create(path)
	if err != nil {
		http.Error(w, "save error", http.StatusInternalServerError)
		return
	}
	defer f.Close()

	io.Copy(f, file)

	res, err := readFile(path)
	if err != nil {
		http.Error(w, "read error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(res)
}

func readFile(p string) (*FileAnalysis, error) {
	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}

	res := FileAnalysis{
		FileName:  p,
		WordCount: 0,
		CharCount: 0,
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

		res.WordCount += uint64(len(words))
		res.CharCount += uint64(utf8.RuneCountInString(chunk))
	}

	if tail != "" {
		res.WordCount++ // don't forget about the last word!!!!
	}

	return &res, nil
}

func main() {
	http.HandleFunc("/upload", uploadHandler)
	http.ListenAndServe(":8080", nil)
}

func isSpace(r rune) bool {
	return unicode.IsSpace(r)
}
