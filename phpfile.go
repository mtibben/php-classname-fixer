package main

import (
	"regexp"
	"strings"
)

type phpfiles map[string]*phpfile

type classname string

func (c classname) namespace() string {
	parts := strings.Split(string(c), `\`)
	ns := strings.Join(parts[1:len(parts)-1], `\`)

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
	path      string
	contents  string
	origClass classname
	newClass  classname
}

func newPhpFile(path string) *phpfile {
	return &phpfile{
		path: path,
	}
}

func (f *phpfile) Contents() string {
	if f.contents == "" {
		f.contents = mustReadFile(f.path)
	}

	return f.contents
}

func (p *phpfile) expectedClassNameFromPath() string {
	path := p.path
	startPos := len(basepath) - 1
	endPos := len(path) - 4
	path = path[startPos:endPos]
	path = strings.Replace(path, "/", `\`, -1)

	return path
}

func (p *phpfile) containsNamespace() bool {
	return strings.Contains(p.Contents(), "\nnamespace ")
}

var namespaceRe = regexp.MustCompile(`\nnamespace\s+(\S+)\s*;`)

func (p *phpfile) getNamespace() string {
	matches := namespaceRe.FindStringSubmatch(p.Contents())

	return matches[1]
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

var useasRe = regexp.MustCompile(`\nuse\s+\S+\s+as\s+(\S+)\s*;`)

func (p *phpfile) getUseAsClasses() (useAsClasses []string) {
	matches := useasRe.FindAllStringSubmatch(p.Contents(), 1)
	if len(matches) == 0 {
		return
	}

	for _, m := range matches {
		useAsClasses = append(useAsClasses, m[len(m)-1])
	}

	return useAsClasses
}

func (p *phpfile) PathDoesntMatchClassname() bool {
	return p.newClass != p.origClass
}

func (p *phpfile) Save() {
	mustWriteFile(p.path, p.Contents())
}
