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

	// fmt.Println(strings.Join(sortedReplacementKeys, "\n"))
	// os.Exit(1)

	replacerArgs := []string{}
	for _, k := range sortedReplacementKeys {
		r := files[k]
		replacerArgs = append(replacerArgs,
			"use "+r.origClass.String(), "use "+r.class.String()[1:],
			`\`+r.origClass.String(), r.class.String(),
			r.origClass.String(), r.class.String())
		// fmt.Println(r.origClass, " -> ", r.class)
	}

	// os.Exit(1)

	replacer.replacer = strings.NewReplacer(replacerArgs...)

	return replacer
}

func (p *phpClassReplacer) updateNamespace(f *phpfile) {

	namespaceLine := "namespace " + f.class.namespace() + ";"

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

var re1 = regexp.MustCompile(`new ([A-Z])`)
var re2 = regexp.MustCompile(`([^\\a-zA-Z])([A-Z][\w_]+)::`)

func (p *phpClassReplacer) fixUnnamespacedClasses(f *phpfile) {
	f.contents = re1.ReplaceAllString(f.Contents(), "new \\$1")
	f.contents = re2.ReplaceAllString(f.Contents(), "$1\\$2::")
}

// replaceOtherClasses causes the classname to become namespaced
func (p *phpClassReplacer) fixNamespacedClassname(f *phpfile) {
	ns := "class \\" + f.class.namespace() + `\`
	f.contents = strings.Replace(f.Contents(), ns, "class ", -1)

	ns = "interface \\" + f.class.namespace() + `\`
	f.contents = strings.Replace(f.Contents(), ns, "interface ", -1)
}

// remove the base namespace from classes
// func (p *phpClassReplacer) replaceRedundantNs(f *phpfile) {
// 	ns := `\` + f.newNamespace() + `\`
// 	f.contents = strings.Replace(f.contents, ns, "", -1)
// }

func (p *phpClassReplacer) UpdateClassnames() {
	for _, f := range p.files {
		fmt.Println("Updating", f.path)
		p.updateNamespace(f)
		p.replaceClasses(f)
		p.fixUnnamespacedClasses(f)
		p.fixNamespacedClassname(f)
		f.Save()
	}
}
