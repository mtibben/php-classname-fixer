package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type phpFinder struct {
	basepath string
	files    phpfiles
}

func newPhpFinder(basepath string) *phpFinder {
	return &phpFinder{
		basepath: basepath,
		files:    phpfiles{},
	}
}

func (p *phpFinder) extractMetadataFromFile(path string) {
	f := newPhpFile(path)
	classnames := f.getClasses()
	if len(classnames) == 0 {
		fmt.Fprintln(os.Stderr, "Warning: couldn't find class in ", path)
		return
	}
	if len(classnames) > 1 {
		fmt.Fprintln(os.Stderr, "Warning: more than 1 class in ", path)
	}
	if f.containsNamespace() {
		f.origClass = classname("\\" + f.getNamespace() + "\\" + classnames[0])
		f.newClass = classname(f.expectedClassNameFromPath())

	} else {
		f.origClass = classname(classnames[0])
		f.newClass = classname(f.expectedClassNameFromPath())
	}

	p.files[f.origClass.String()] = f
}

func (p *phpFinder) findPhpFiles(path string, info os.FileInfo, err error) error {
	if !strings.HasSuffix(path, ".php") {
		return nil
	}

	p.extractMetadataFromFile(path)

	return nil
}

func (p *phpFinder) Find() phpfiles {
	filepath.Walk(p.basepath, p.findPhpFiles)

	return p.files
}
