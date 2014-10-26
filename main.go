package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var basepath string
var outputFile string
var write bool

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "php-classname-fixer [-w] [path]")
		flag.PrintDefaults()
	}
	flag.BoolVar(&write, "w", false, "write - saves changes to files")
	flag.StringVar(&outputFile, "o", "", "output class replacements to file")
	flag.Parse()

	basepath, _ = filepath.Abs(flag.Arg(0))
	basepath += "/"
	fmt.Printf("Searching %s\n", basepath)
}

func main() {
	finder := newPhpFinder(basepath)
	files := finder.Find()

	classCount := 0
	replacements := map[string]string{}
	for _, v := range files {
		if v.PathDoesntMatchClassname() {
			c := strings.Trim(v.origClass.String(), "\\")
			d := strings.Trim(v.newClass.String(), "\\")
			fmt.Printf("   %s -> %s\n", c, d)
			replacements[c] = d
			classCount++
		}
	}
	sortedReplacementKeys := sliceKeys2(replacements)
	sort.Sort(ByLength(sortedReplacementKeys))

	replacementStr := ""
	for _, k := range sortedReplacementKeys {
		replacementStr += fmt.Sprintf("%s %s\n", k, replacements[k])
	}

	if outputFile != "" {
		mustWriteFile(outputFile, replacementStr)
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

func sliceKeys2(ss map[string]string) (s []string) {
	for k, _ := range ss {
		s = append(s, k)
	}
	return
}
