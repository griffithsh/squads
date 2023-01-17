package embedded

import (
	"embed"
	"fmt"
	"path"
	"strings"
)

//go:embed content/*
var content embed.FS

// Filenames embedded in the squads binary.
func Filenames() ([]string, error) {
	// dirHandler handlers directories recursively.
	var dirHandler func(string) ([]string, error)
	dirHandler = func(dirName string) ([]string, error) {
		dirs, err := content.ReadDir(dirName)
		result := make([]string, 0, len(dirs))
		if err != nil {
			return nil, fmt.Errorf("read %s: %v", dirName, err)
		}
		for _, dir := range dirs {
			if dir.IsDir() {
				files, err := dirHandler(path.Join(dirName, dir.Name()))
				if err != nil {
					return nil, err
				}
				result = append(result, files...)
				continue
			}

			result = append(result, path.Join(dirName, dir.Name()))
		}
		return result, nil
	}

	files, err := dirHandler("content")
	if err != nil {
		return nil, fmt.Errorf("dirHandler: %v", err)
	}

	result := make([]string, len(files))
	for i, f := range files {
		result[i] = strings.TrimPrefix(f, "content/")
	}
	return result, nil
}

// Get the contents of an embedded file.
func Get(filename string) ([]byte, error) {
	return content.ReadFile(path.Join("content", filename))
}
