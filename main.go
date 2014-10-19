package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

var basepath string
var dryrun bool

func init() {
	flag.BoolVar(&dryrun, "n", false, "dry-run - shows what classes would have changed")
	flag.Parse()

	basepath, _ = filepath.Abs(flag.Arg(0))
	basepath += "/"
	fmt.Printf("Searching %s\n", basepath)
}

func main() {
	finder := newPhpFinder(basepath)
	files := finder.find()

	classCount := 0
	for k, v := range files {
		if v.PathDoesntMatchClass() {
			fmt.Printf("   %s -> %s\n", k, v.newClass.String())
			classCount++
		}
	}
	if dryrun {
		os.Exit(0)
	}

	fmt.Printf("Replacing %d classnames in %d files\n", classCount, len(files))

	replacer := newPhpClassReplacer(basepath, files)
	replacer.UpdateClassnames()
	fmt.Print("\n")
}
