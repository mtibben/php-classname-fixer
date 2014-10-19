package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

var basepath string
var write bool

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "php-classname-fixer [-w] [path]")
		flag.PrintDefaults()
	}
	flag.BoolVar(&write, "w", false, "write - saves changes to files")
	flag.Parse()

	basepath, _ = filepath.Abs(flag.Arg(0))
	basepath += "/"
	fmt.Printf("Searching %s\n", basepath)
}

func main() {
	finder := newPhpFinder(basepath)
	files := finder.Find()

	classCount := 0
	for k, v := range files {
		if v.PathDoesntMatchClassname() {
			fmt.Printf("   %s -> %s\n", k, v.newClass.String())
			classCount++
		}
	}

	if !write {
		fmt.Printf("Would have replaced %d classnames in %d files\n", classCount, len(files))
		fmt.Println("Run with the -w flag to actually write files")
		os.Exit(0)
	} else {
		fmt.Printf("Replacing %d classnames in %d files\n", classCount, len(files))

		replacer := newPhpClassReplacer(basepath, files)
		replacer.UpdateClassnames()
		fmt.Print("\n")
	}
}
