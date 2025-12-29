package scan

import (
	"io/fs"
	"path/filepath"
)

type Result struct {
	Path string
}

func FindFolders(root, name string) ([]Result, error) {
	var results []Result

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		if !d.IsDir() {
			return nil
		}

		switch d.Name() {
		case ".git", ".idea", ".vscode":
			return filepath.SkipDir
		}

		if d.Name() == name {
			results = append(results, Result{Path: path})
			return filepath.SkipDir
		}

		return nil
	})

	return results, err
}

func FindFoldersGlob(root string, pattern string) ([]Result, error) {
	var results []Result

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		if !d.IsDir() {
			return nil
		}

		switch d.Name() {
		case ".git", ".idea", ".vscode":
			return filepath.SkipDir
		}

		matched, err := filepath.Match(pattern, d.Name())
		if err != nil {
			return err
		}

		if matched {
			results = append(results, Result{Path: path})
			return filepath.SkipDir
		}

		return nil
	})

	return results, err
}
