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
	replacer0             *strings.Replacer
	replacer1             *strings.Replacer
	replacer2             *strings.Replacer
	replacer3             *strings.Replacer
	origClasses           []string
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

	replacer0Args := []string{}
	replacer1Args := []string{}
	replacer2Args := []string{}
	replacer3Args := []string{}
	for _, k := range sortedReplacementKeys {
		r := files[k]
		if r.PathDoesntMatchClassname() {
			replacer0Args = append(replacer0Args,
				` \`+r.origClass.String(), ` \`+r.newClass.String(),
			)

			replacer1Args = append(replacer1Args,
				r.origClass.String(), r.newClass.String(),
			)

			replacer2Args = append(replacer2Args,
				`(`+r.origClass.String()+` $`, `(\`+r.newClass.String()+` $`,
				` `+r.origClass.String()+` $`, ` \`+r.newClass.String()+` $`,
				`'`+r.origClass.String()+`'`, `'\`+r.newClass.String()+`'`,
				`\`+r.origClass.String(), `\`+r.newClass.String(),
			)

			replacer3Args = append(replacer3Args,
				`\`+r.origClass.String(), `\`+r.newClass.String(),
				r.origClass.String(), `\`+r.newClass.String(),
			)
		}
	}

	replacer.replacer0 = strings.NewReplacer(replacer0Args...)
	replacer.replacer1 = strings.NewReplacer(replacer1Args...)
	replacer.replacer2 = strings.NewReplacer(replacer2Args...)
	replacer.replacer3 = strings.NewReplacer(replacer3Args...)

	return replacer
}

func (p *phpClassReplacer) getReplacements() (map[string]string, []string) {
	replacements := map[string]string{}

	for _, v := range p.files {
		if v.PathDoesntMatchClassname() {
			c := strings.Trim(v.origClass.String(), "\\")
			d := strings.Trim(v.newClass.String(), "\\")
			replacements[c] = d
		}
	}
	sortedReplacementKeys := sliceKeys2(replacements)
	sort.Sort(ByLength(sortedReplacementKeys))

	replacementStr := ""
	for _, k := range sortedReplacementKeys {
		replacementStr += fmt.Sprintf("%s %s\n", k, replacements[k])
	}

	return replacements, sortedReplacementKeys
}

func (p *phpClassReplacer) namespaceLine(f *phpfile) string {
	namespaceLine := "\nnamespace " + f.newClass.namespace() + ";"

	return namespaceLine
}

var reFindClassesWithNew = regexp.MustCompile(`new\s+([A-Z][\w_]+)`)
var reFindClassesWithStaticCall = regexp.MustCompile(`([^\w_\\])([A-Z][\w_\\]+)::`)
var reFindClassesWIthInstanceOf = regexp.MustCompile(`instanceof\s+([A-Z][\w_\\]+)`)
var reFindClassesInFunctionSigs = regexp.MustCompile(`([\(,]\s*)([A-Z][\w_\\]+) \$`)

func (p *phpClassReplacer) fixTheRest(f *phpfile) {
	ignoreClasses := f.getUseAsClasses()

	uac := ""
	if len(ignoreClasses) > 0 {
		uac = `(` + strings.Join(ignoreClasses, "|") + `)`
	}

	f.contents.therest = reFindClassesWithNew.ReplaceAllStringFunc(f.contents.therest, func(s string) string {
		if len(uac) > 0 {
			re4 := regexp.MustCompile(`new\s+` + uac)
			if re4.MatchString(s) {
				return s
			}
		}

		return reFindClassesWithNew.ReplaceAllString(s, "new \\$1")
	})

	f.contents.therest = reFindClassesWithStaticCall.ReplaceAllStringFunc(f.contents.therest, func(s string) string {
		if len(uac) > 0 {
			re4 := regexp.MustCompile(`([^\\a-zA-Z])` + uac + `::`)
			if re4.MatchString(s) {
				return s
			}
		}

		return reFindClassesWithStaticCall.ReplaceAllString(s, "$1\\$2::")
	})

	f.contents.therest = reFindClassesWIthInstanceOf.ReplaceAllStringFunc(f.contents.therest, func(s string) string {
		if len(uac) > 0 {
			re4 := regexp.MustCompile(`instanceof\s+` + uac)
			if re4.MatchString(s) {
				return s
			}
		}

		return reFindClassesWIthInstanceOf.ReplaceAllString(s, "instanceof \\$1")
	})

	f.contents.therest = reFindClassesInFunctionSigs.ReplaceAllStringFunc(f.contents.therest, func(s string) string {
		if len(uac) > 0 {
			re4 := regexp.MustCompile(`[\( ]` + uac + ` \$`)
			if re4.MatchString(s) {
				return s
			}
		}

		return reFindClassesInFunctionSigs.ReplaceAllString(s, "$1\\$2 $")
	})
}

func inArray(s string, ss []string) bool {
	for _, a := range ss {
		if a == s {
			return true
		}
	}

	return false
}

var reGetExtendClasses = regexp.MustCompile(`(?ms)([\w\\]+)`)

func (p *phpClassReplacer) fixPostClass(f *phpfile) {
	ignoreClasses := f.getUseAsClasses()
	f.contents.postClass = p.replacer3.Replace(f.contents.postClass)
	f.contents.postClass = reGetExtendClasses.ReplaceAllStringFunc(f.contents.postClass, func(s string) string {
		if s == "extends" || s == "implements" || s[0] == '\\' {
			return s
		}

		if inArray(s, ignoreClasses) {
			return s
		}

		return `\` + s
	})
}

var useasRe3 = regexp.MustCompile(`\nuse\s+(\S+)(\s+as\s+(\S+))?\s*\;`)

type useAsLine struct {
	useClass classname
	asClass  classname
}

func (u *useAsLine) String() string {
	if u.asClass == "" || u.useClass.class() == u.asClass.class() {
		return fmt.Sprintf("\nuse %s;", u.useClass)
	} else {
		return fmt.Sprintf("\nuse %s as %s;", u.useClass, u.asClass.class())
	}
}

func (p *phpClassReplacer) updatePostNs(t string, f *phpfile) string {
	return useasRe3.ReplaceAllStringFunc(t, func(s string) string {

		m := useasRe3.FindStringSubmatch(s)

		useClass := p.replacer1.Replace(m[1])
		asClass := m[1]
		if m[3] != "" {
			asClass = m[3]
		}

		l := useAsLine{classname(useClass), classname(asClass)}

		return l.String()
	})
}

var whitespaceInStaticCall = regexp.MustCompile(`(\S)\s+::(\S)`)

func somePsr2Fixes(s string) string {
	return whitespaceInStaticCall.ReplaceAllString(s, "$1::$2")
}

func (p *phpClassReplacer) UpdateClassnames() {
	for _, f := range p.files {
		fmt.Print(".")

		alreadyIsNamespaced := f.containsNamespace()

		f.contents.therest = somePsr2Fixes(f.contents.therest)
		f.contents.ns = p.namespaceLine(f)
		f.contents.postNs = p.updatePostNs(f.contents.postNs, f)

		if !alreadyIsNamespaced {
			p.fixPostClass(f)
			p.fixTheRest(f)
		}

		f.contents.classname = f.newClass.class()
		f.contents.postClass = p.replacer0.Replace(f.contents.postClass)
		f.contents.therest = p.replacer2.Replace(f.contents.therest)

		f.Save()
	}
}
