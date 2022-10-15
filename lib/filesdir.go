package lib

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
)

// GetFilesByExt search recursively into the given root path, to seek the files that
// match the specified file extension
//
// - root [string] ~ Root path to search in on
//
// - ext [string] ~ File extension to look for
func GetFilesByExt(root, ext string) []string {
	// 👆🏽 two vars, one type declaration, sample

	var files []string
	filepath.WalkDir(root, func(path string, d fs.DirEntry, e error) error {

		if e != nil { return e }

		if filepath.Ext(d.Name()) == ext {
			files = append(files, path)
		}

		return nil
	})

	return files
}

// GetFilesByName search recursively into the given root path, to seek the files that
// match the specified file filename
//
// - root [string] ~ Root path to search in on
//
// - fileNeme [string] ~ File name to look for
func GetFilesByName(root, fileNeme string) []string {
	// 👆🏽 two vars, one type declaration, sample

	var files []string
	filepath.WalkDir(root, func(path string, d fs.DirEntry, e error) error {

		if e != nil { return e }

		if d.Name() == fileNeme {
			files = append(files, path)
		}

		return nil
	})

	return files
}

func FileExists(filename string) (bool, error) {
	_, err := os.Stat(filename)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}