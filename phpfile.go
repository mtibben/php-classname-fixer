package main

import (
	"fmt"
	"regexp"
	"strings"
)

type phpfiles map[string]*phpfile

type classname string

func (c classname) namespace() string {
	parts := strings.Split(string(c), `\`)
	ns := strings.Join(parts[0:len(parts)-1], `\`)

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
	contents  phpFileFragments
	origClass classname
	newClass  classname
}

func newPhpFile(path string) *phpfile {
	return &phpfile{
		path: path,
	}
}

var fragmentsRe = regexp.MustCompile("(?ms)(\\<\\?php(.*?))((\nnamespace[^\n]+)(.*?))?(\n\\s*(class|trait|interface|abstract class|final class)\\s+)(\\S+)(.*?\\{)(.+)")

func (f *phpfile) SetContents(s string) {

	matches := fragmentsRe.FindStringSubmatch(s)

	if len(matches) < 9 {
		fmt.Println(f.path)
		fmt.Println(s)
		panic("Unexpected length")
	}

	f.contents = phpFileFragments{
		preNs:     matches[1],
		ns:        matches[4],
		postNs:    matches[5],
		preClass:  matches[6],
		classname: matches[8],
		postClass: matches[9],
		therest:   matches[10],
	}

	if f.contents.ns == "" {
		f.contents.preNs = "<?php\n"
		f.contents.postNs = matches[2] + matches[5]
	}
}

func (f *phpfile) Contents() string {
	if f.contents.String() == "" {
		f.SetContents(mustReadFile(f.path))
	}

	return f.contents.String()
}

func (p *phpfile) expectedClassNameFromPath() string {
	path := p.path
	startPos := len(basepath) - 1
	endPos := len(path) - 4
	path = path[startPos:endPos]
	path = strings.Replace(path, "/", `\`, -1)
	path = strings.Trim(path, "\\")

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

var useasRe1 = regexp.MustCompile(`\nuse\s+\S+\s+as\s+(\S+)\s*\;`)
var useasRe2 = regexp.MustCompile(`\nuse\s+\S+\\(\S+)\s*\;`)

func (p *phpfile) getUseAsClasses() (useAsClasses []string) {
	matches1 := useasRe1.FindAllStringSubmatch(p.contents.postNs, -1)
	for _, m := range matches1 {
		useAsClasses = append(useAsClasses, m[1])
	}
	matches2 := useasRe2.FindAllStringSubmatch(p.contents.postNs, -1)
	for _, m := range matches2 {
		useAsClasses = append(useAsClasses, m[1])
	}

	return useAsClasses
}

func (p *phpfile) PathDoesntMatchClassname() bool {
	return p.newClass != p.origClass
}

func (p *phpfile) Save() {
	mustWriteFile(p.path, p.Contents())
}

type phpFileFragments struct {
	preNs     string
	ns        string
	postNs    string
	preClass  string
	classname string
	postClass string
	therest   string
}

func (f *phpFileFragments) String() string {
	return f.preNs + f.ns + f.postNs + f.preClass + f.classname + f.postClass + f.therest
}
