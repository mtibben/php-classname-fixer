package main

import (
	"path/filepath"
	"regexp"
	"strings"
)

type classname string

func (c classname) namespace() string {
	parts := strings.Split(string(c), `\`)
	ns := strings.Join(parts[1:len(parts)-1], `\`)

	if ns == "" {
		panic("Couldn't create new namespace")
	}

	return ns
}

func (c classname) class() string {
	parts := strings.Split(string(c), `\`)

	return parts[len(parts)-1]
}

func (c classname) String() string {
	return string(c)
}

type phpfile struct {
	basepath    string
	newbasepath string
	path        string
	contents    string
	origClass   classname
	class       classname
}

func newPhpFile(basepath, path string) *phpfile {
	return &phpfile{
		basepath: basepath,
		path:     path,
	}
}

func (f *phpfile) Contents() string {
	if f.contents == "" {
		f.contents = mustReadFile(f.path)
	}

	return f.contents
}

func (p *phpfile) namespacedPsrClassNameFromPath() string {
	path := p.path
	startPos := len(p.basepath) - 1
	endPos := len(path) - 4
	path = path[startPos:endPos]
	path = strings.Replace(path, "/", `\`, -1)

	return path
}

func (p *phpfile) containsNamespace() bool {
	return strings.Contains(p.Contents(), "\nnamespace ")
}

var classRe = regexp.MustCompile(`\n\s*((abstract)?\s*(final)?\s*class|interface)\s+(\S+)`)

func (p *phpfile) getClasses() (classnames []string) {
	matches := classRe.FindAllStringSubmatch(p.Contents(), 2)
	if len(matches) == 0 {
		return
	}

	for _, m := range matches {
		classnames = append(classnames, m[len(m)-1])
	}

	return classnames
}

func (p *phpfile) Save() {
	if p.newbasepath == "" {
		mustWriteFile(p.path, p.Contents())
	} else {
		mustDeleteFile(p.path)
		mustWriteFile(p.newPath(), p.Contents())
	}

}

func (p *phpfile) newPath() string {
	return filepath.Clean(filepath.Join(
		p.newbasepath,
		strings.Replace(p.class.String(), `\`, `/`, -1)+".php"))
}
