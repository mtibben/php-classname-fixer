package main

import (
	"fmt"
	"os"
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

func (p *phpFinder) extractClassFromFile(path string) {
	f := newPhpFile(p.basepath, path)
	if f.containsNamespace() {
		f.origClass = classname(f.namespacedPsrClassNameFromPath())
		f.class = classname(f.origClass)

	} else {

		classnames := f.getClasses()
		if len(classnames) == 0 {
			fmt.Println("Warning: couldn't find class in ", path)
			return
		}
		if len(classnames) > 1 {
			fmt.Println("Warning: more than 1 class in ", path)
		}

		f.origClass = classname(classnames[0])
		f.class = classname(f.namespacedPsrClassNameFromPath())

		classparts := strings.Split(f.class.String()[1:], `\`)
		expectedClass := strings.Join(classparts, `_`)

		if expectedClass != f.origClass.String() {
			fmt.Println("Warning: unexpected classname", f.origClass, "expected", expectedClass)
		}
	}

	p.files[f.origClass.String()] = f
}

func (p *phpFinder) findPhpFiles(path string, info os.FileInfo, err error) error {
	if !strings.HasSuffix(path, ".php") {
		return nil
	}

	p.extractClassFromFile(path)

	return nil
}
