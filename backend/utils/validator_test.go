package utils

import "testing"

func TestAllowedFileTypes(t *testing.T) {

	if !IsAllowedFileType("resume.pdf") {
		t.Errorf("pdf should be allowed")
	}

	if !IsAllowedFileType("image.jpg") {
		t.Errorf("jpg should be allowed")
	}

	if IsAllowedFileType("virus.exe") {
		t.Errorf("exe should not be allowed")
	}
}
