package embedded

import (
	"testing"
)

func TestFilenamesAndGet(t *testing.T) {
	files, err := Filenames()
	if err != nil {
		t.Fatalf("Filenames did not return files: %v", err)
	}
	for _, file := range files {
		b, err := Get(file)
		if err != nil {
			t.Errorf("get %q: %v", file, err)
		}

		if len(b) == 0 {
			t.Errorf("%q: zero byte file", file)
		}
	}
}
