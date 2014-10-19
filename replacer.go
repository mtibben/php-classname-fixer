package main

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
)

type phpClassReplacer struct {
	basepath              string
	sortedReplacementKeys []string
	files                 phpfiles
	replacer              *strings.Replacer
}

func sliceKeys(ss phpfiles) (s []string) {
	for k, _ := range ss {
		s = append(s, k)
	}
	return
}

func newPhpClassReplacer(basepath string, files phpfiles) *phpClassReplacer {

	replacer := &phpClassReplacer{
		basepath: basepath,
		files:    files,
	}
	sortedReplacementKeys := sliceKeys(files)
	sort.Sort(ByLength(sortedReplacementKeys))

	replacerArgs := []string{}
	for _, k := range sortedReplacementKeys {
		r := files[k]
		replacerArgs = append(replacerArgs,
			"use "+r.origClass.String(), "use "+r.newClass.String()[1:],
			`\`+r.origClass.String(), r.newClass.String(),
			r.origClass.String(), r.newClass.String())
	}

	replacer.replacer = strings.NewReplacer(replacerArgs...)

	return replacer
}

func (p *phpClassReplacer) updateNamespace(f *phpfile) {

	namespaceLine := "namespace " + f.newClass.namespace() + ";"

	if f.containsNamespace() {
		lines := strings.Split(f.Contents(), "\n")
		newparts := []string{}
		done := false
		for _, l := range lines {
			if !done && strings.HasPrefix(l, "namespace") {
				newparts = append(newparts, namespaceLine)
				done = true
			} else {
				newparts = append(newparts, l)
			}
		}

		f.contents = strings.Join(newparts, "\n")
	} else {
		parts := strings.SplitAfterN(f.Contents(), "<?php\n", 2)

		if len(parts) != 2 {
			panic("Less parts than expected")
		}

		f.contents = parts[0] +
			"\n" + namespaceLine + "\n" +
			parts[1]
	}
}

func (p *phpClassReplacer) replaceClasses(f *phpfile) {
	f.contents = p.replacer.Replace(f.Contents())
}

var re1 = regexp.MustCompile(`new\s+([A-Z][\w_]+)`)
var re2 = regexp.MustCompile(`([^\\a-zA-Z])([A-Z][\w_]+)::`)

func (p *phpClassReplacer) fixUnnamespacedClasses(f *phpfile) {
	useAsClasses := f.getUseAsClasses()
	uac := ""
	if len(useAsClasses) > 0 {
		uac = `(` + strings.Join(useAsClasses, "|") + `)`
	}

	f.contents = re1.ReplaceAllStringFunc(f.Contents(), func(s string) string {
		if len(uac) > 0 {
			re3 := regexp.MustCompile(`new\s+` + uac)
			if re3.MatchString(s) {
				return s
			}
		}

		return re1.ReplaceAllString(s, "new \\$1")
	})

	f.contents = re2.ReplaceAllStringFunc(f.Contents(), func(s string) string {
		if len(uac) > 0 {
			re3 := regexp.MustCompile(`([^\\a-zA-Z])` + uac + `::`)
			if re3.MatchString(s) {
				return s
			}
		}

		return re2.ReplaceAllString(s, "$1\\$2::")
	})

}

// replaceOtherClasses causes the classname to become namespaced
func (p *phpClassReplacer) fixNamespacedClassname(f *phpfile) {
	ns := "class \\" + f.newClass.namespace() + `\`
	f.contents = strings.Replace(f.Contents(), ns, "class ", -1)

	ns = "interface \\" + f.newClass.namespace() + `\`
	f.contents = strings.Replace(f.Contents(), ns, "interface ", -1)
}

func (p *phpClassReplacer) UpdateClassnames() {
	for _, f := range p.files {
		fmt.Print(".")
		p.updateNamespace(f)
		p.replaceClasses(f)
		p.fixUnnamespacedClasses(f)
		p.fixNamespacedClassname(f)
		f.Save()
	}
}
