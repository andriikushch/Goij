package main

import (
	"flag"
	"fmt"
	"github.com/j7mbo/goij/src/TypeRegistry"
	"path/filepath"
	"strings"
)

/* Provides the ability to handle multiple parameters for a single flag. */
type arrayFlags []string

/* Allows user to pass -exclude path1 -exclude path2 */
var excludeFlags arrayFlags

/* Only needed to fulfil flag interface. */
func (i *arrayFlags) String() string {
	return ""
}

/* Only needed to fulfil flag interface. */
func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)

	return nil
}

func main() {
	// @todo At least let the user know what's going on, no need for a progress bar

	reset, dir, file, exclude := declareCommandLineFlags()

	gen := TypeRegistry.NewAutoRegistryGenerator(TypeRegistry.AutoGeneratedRegistryWriter{})

	if *reset == true {
		gen.Reset(*file)

		return
	}

	/* Prevents error with directory having same name as file, for example. */
	if !strings.HasSuffix(*file, ".go") {
		*file += ".go"
	}

	absoluteDir, err := filepath.Abs(*dir)

	if err != nil {
		panic(fmt.Sprintf("Could not retrieve absolute filepath for dir: '%s', got error: '%s'", *dir, err.Error()))
	}

	absoluteFile, err := filepath.Abs(*file)

	if err != nil {
		panic(fmt.Sprintf("Could not retrieve absolute filepath for file: '%s', got error: '%s'", *file, err.Error()))
	}

	gen.Generate(absoluteFile, absoluteDir, exclude...)
}

func declareCommandLineFlags() (reset *bool, dir *string, file *string, exclude arrayFlags) {
	exclude = excludeFlags

	reset = flag.Bool("reset", false, "Reset the registry file")
	file = flag.String("o", "./Registry.go", "A relative or absolute filepath to write the registry to")
	dir = flag.String("dir", ".", "A relative or absolute directory to recurse and generate the type registry from")
	flag.Var(&exclude, "exclude", "Directories to exclude parsing for registry, such as vendor/")

	flag.Parse()

	return
}
