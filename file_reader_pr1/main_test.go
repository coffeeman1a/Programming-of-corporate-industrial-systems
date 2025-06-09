package main

import (
	"os"
	"testing"
)

func TestReadFile(t *testing.T) {
	content := "Hello world!\nЭто тест."
	err := os.WriteFile("test.txt", []byte(content), 0644)
	if err != nil {
		t.Fatalf("Ошибка при создании файла: %v", err)
	}
	defer os.Remove("test.txt") // удалить после теста

	analysis, err := readFile("test.txt")
	if err != nil {
		t.Fatalf("Ошибка при чтении файла: %v", err)
	}

	if analysis.wordCount != 4 {
		t.Errorf("Ожидали 4 слова, а получили %d", analysis.wordCount)
	}
}

func TestReadFile_NotExists(t *testing.T) {
	_, err := readFile("definitely_no_such_file.txt")
	if err == nil {
		t.Errorf("Ожидали ошибку при чтении несуществующего файла, а получили nil")
	}
}

func TestReadFile_EmptyFile(t *testing.T) {
	filename := "empty_test.txt"
	err := os.WriteFile(filename, []byte(""), 0644)
	if err != nil {
		t.Fatalf("Не удалось создать пустой файл: %v", err)
	}
	defer os.Remove(filename)

	analysis, err := readFile(filename)
	if err != nil {
		t.Fatalf("Не ожидали ошибку при чтении пустого файла: %v", err)
	}
	if analysis.wordCount != 0 {
		t.Errorf("Для пустого файла ожидали 0 слов, получили %d", analysis.wordCount)
	}
	if analysis.charCount != 0 {
		t.Errorf("Для пустого файла ожидали 0 символов, получили %d", analysis.charCount)
	}
}
