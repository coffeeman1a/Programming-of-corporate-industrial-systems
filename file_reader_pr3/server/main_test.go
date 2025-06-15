package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestIsSpace checks the isSpace helper with various runes.
func TestIsSpace(t *testing.T) {
	s := []struct {
		r    rune
		want bool
	}{
		{' ', true},
		{'\n', true},
		{'\t', true},
		{'a', false},
		{'一', false},
	}
	for _, tt := range s {
		if got := isSpace(tt.r); got != tt.want {
			t.Errorf("isSpace(%q) = %v; want %v", tt.r, got, tt.want)
		}
	}
}

// helper to write temp file and return its path
func writeTempFile(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", "testfile-*.txt")
	if err != nil {
		t.Fatalf("creating temp file: %v", err)
	}
	defer f.Close()
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("writing temp file: %v", err)
	}
	return f.Name()
}

func TestReadFile_Empty(t *testing.T) {
	path := writeTempFile(t, "")
	defer os.Remove(path)

	res, err := readFile(path)
	if err != nil {
		t.Fatalf("readFile error: %v", err)
	}
	if res.WordCount != 0 {
		t.Errorf("empty file WordCount = %d; want 0", res.WordCount)
	}
	if res.CharCount != 0 {
		t.Errorf("empty file CharCount = %d; want 0", res.CharCount)
	}
}

func TestReadFile_NormalContent(t *testing.T) {
	content := "Hello, world!\nThis is a test.  "
	path := writeTempFile(t, content)
	defer os.Remove(path)

	res, err := readFile(path)
	if err != nil {
		t.Fatalf("readFile error: %v", err)
	}
	// "Hello," and "world!" count as words
	wantWords := uint64(6) // Hello, world!, This, is, a, test.
	if res.WordCount != wantWords {
		t.Errorf("WordCount = %d; want %d", res.WordCount, wantWords)
	}
	// count of runes equals len([]rune(content))
	runes := []rune(content)
	if res.CharCount != uint64(len(runes)) {
		t.Errorf("CharCount = %d; want %d", res.CharCount, len(runes))
	}
}

func TestReadFile_Unicode(t *testing.T) {
	content := "Привет мир"
	path := writeTempFile(t, content)
	defer os.Remove(path)

	res, err := readFile(path)
	if err != nil {
		t.Fatalf("readFile error: %v", err)
	}
	if res.WordCount != 2 {
		t.Errorf("WordCount = %d; want 2", res.WordCount)
	}
	if res.CharCount != uint64(len([]rune(content))) {
		t.Errorf("CharCount = %d; want %d", res.CharCount, len([]rune(content)))
	}
}

func TestUploadHandler_Success(t *testing.T) {
	// Prepare a small text file
	content := "Go testing sample"
	tempPath := writeTempFile(t, content)
	defer os.Remove(tempPath)

	// Create multipart form data
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	field, err := writer.CreateFormFile("uploadfile", filepath.Base(tempPath))
	if err != nil {
		t.Fatalf("CreateFormFile error: %v", err)
	}
	file, err := os.Open(tempPath)
	if err != nil {
		t.Fatalf("Open temp file: %v", err)
	}
	io.Copy(field, file)
	file.Close()
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/upload", &buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()

	uploadHandler(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Status = %d; want %d", resp.StatusCode, http.StatusOK)
	}

	// Decode response
	var analysis FileAnalysis
	if err := json.NewDecoder(resp.Body).Decode(&analysis); err != nil {
		t.Fatalf("decoding response: %v", err)
	}
	if analysis.WordCount != 3 {
		t.Errorf("WordCount = %d; want 3", analysis.WordCount)
	}
}

func TestUploadHandler_MissingFile(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/upload", nil)
	w := httptest.NewRecorder()
	uploadHandler(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Status = %d; want %d", resp.StatusCode, http.StatusBadRequest)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	if !strings.Contains(string(body), "file error") {
		t.Errorf("Body = %q; want file error message", body)
	}
}
