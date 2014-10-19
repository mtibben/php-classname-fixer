package main

import (
	"flag"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
)

type phpfiles map[string]*phpfile

var basepath string
var locations []location
var dryrun bool

type location struct {
	path            string
	findClasses     bool
	replaceClasses  bool
	prefixNamespace string
	newBasePath     string
}

func init() {
	usr, _ := user.Current()
	home := usr.HomeDir

	basepath = home + "/Projects/99designs/contests/"
	locations = []location{
		// {"bundles/", true, true, "\\NinetyNine\\ContestsBundles", "src/"},
		// {"controllers/", true, true, "\\NinetyNine\\ContestsSpfControllers", "src/"},
		// {"classes/", true, true, "\\NinetyNine", "src/"},
		// {"tests/unit/", true, true, "\\NinetyNine\\ContestsTests\\Unit", "src/"},
		// {"tests/helpers/", true, true, "\\NinetyNine\\ContestsTests\\Helpers", "src/"},
		// {"tests/system/", true, true, "\\NinetyNine\\ContestsTests\\System", "src/"},
		{"src/", true, true, "", ""},
	}

	flag.BoolVar(&dryrun, "n", false, "dry-run - shows what classes would have changed")
	flag.Parse()
}

// problems:
// \Help\LoggedInHelpUrlBuilder -> \Help\Logged\InHelpUrlBuilder
// Contests_GoogleTagManagerPixelFactory::create();
// \NinetyNine\Contests\Google\TagManagerPixelFactory::create();
func main() {
	files := phpfiles{}
	for _, l := range locations {
		if l.findClasses {
			fullpath := basepath + l.path
			finder := newPhpFinder(fullpath)
			filepath.Walk(fullpath, finder.findPhpFiles)

			for k, v := range finder.files {
				if l.findClasses {
					if l.newBasePath != "" {
						v.newbasepath = basepath + l.newBasePath
					}
					v.class = classname(l.prefixNamespace + v.class.String())
				}
				files[k] = v
			}
		}
	}

	if dryrun {
		fmt.Println("Found classnames:")
		for k, v := range files {
			fmt.Printf("   %s -> %s\n", k, v.class.String())
		}

		os.Exit(0)
	}

	fmt.Println("Replacing classnames")

	replacer := newPhpClassReplacer(basepath, files)
	replacer.UpdateClassnames()

}
