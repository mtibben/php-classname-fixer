package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

func mustReadFile(path string) string {
	c, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	return string(c)
}

func mustWriteFile(path, contents string) {
	// fmt.Println("Writing", path)
	err := os.MkdirAll(filepath.Dir(path), 0755)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(path, []byte(contents), 0644)
	if err != nil {
		panic(err)
	}
}

func mustDeleteFile(path string) {
	err := os.Remove(path)
	if err != nil {
		panic(err)
	}
}
