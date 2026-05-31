package utils

import (
	"os"
	"testing"
)

func TestGenerateSHA256(t *testing.T) {

	testFile := "test.txt"

	err := os.WriteFile(
		testFile,
		[]byte("hello distribyte"),
		0644,
	)

	if err != nil {
		t.Fatal(err)
	}

	defer os.Remove(testFile)

	hash1, err := GenerateSHA256(testFile)

	if err != nil {
		t.Fatal(err)
	}

	hash2, err := GenerateSHA256(testFile)

	if err != nil {
		t.Fatal(err)
	}

	if hash1 != hash2 {
		t.Errorf("hashes should match")
	}
}
