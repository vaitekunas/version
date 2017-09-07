package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

// help prints the usage of the cli
func help() {

	c := color.New(color.FgHiBlue)
	b := color.New(color.Bold)
	out := func(s string) string { return fmt.Sprintf(" %s  %s", c.Sprint("◈"), b.Sprint(s)) }

	fmt.Fprintf(os.Stderr, "\nUsage of %s:\n\n", b.Sprint(os.Args[0]))
	fmt.Fprintf(os.Stderr, "version [command] [arguments]\n\n")
	fmt.Fprintf(os.Stderr, "The commands are:\n\n")
	fmt.Fprintf(os.Stderr, fmt.Sprintf("\t%s\t%s", out("increase"), "increases the version by a major/minor/patch tick\n"))
	fmt.Fprintf(os.Stderr, "\n")
  fmt.Fprintf(os.Stderr, "Using \"version [--root=\"\"] [--all]\" lists available releases/versions\n")
	fmt.Fprintf(os.Stderr, "Use \"version help [command]\" for more information about a command\n\n")

}

// man shows cli manual
func man(cmd string) {
	c := color.New(color.FgHiBlue)
	b := color.New(color.Bold)
	out := func(s string) string { return fmt.Sprintf(" %s  %s", c.Sprint("◈"), b.Sprint(s)) }

	switch strings.ToLower(cmd) {

	case "increase":
		fmt.Fprintf(os.Stderr, "\nUsage of %s:\n\n", b.Sprint("version increase"))
		fmt.Fprintf(os.Stderr, "version increase [{--major, --minor, --patch}] [--special=\"\"] [--build=\"\"]\n\n")
		fmt.Fprintf(os.Stderr, "The arguments are:\n\n")
		fmt.Fprintf(os.Stderr, fmt.Sprintf("\t%s\t%s", out("--major"), "increase version by a major tick\n"))
		fmt.Fprintf(os.Stderr, fmt.Sprintf("\t%s\t%s", out("--minor"), "increase version by a minor tick\n"))
		fmt.Fprintf(os.Stderr, fmt.Sprintf("\t%s\t%s", out("--patch"), "increase version by a patch tick\n"))
    fmt.Fprintf(os.Stderr, fmt.Sprintf("\t%s\t%s", out("--special"), "specify pre-release version\n"))
    fmt.Fprintf(os.Stderr, fmt.Sprintf("\t%s\t%s", out("--build"), "add build-related metadata\n"))
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Only a single tick option (major/minor/patch) is allowed per increase\n")
    fmt.Fprintf(os.Stderr, "Setting special and build identifiers without tick updates will use the current version\n")
    fmt.Fprintf(os.Stderr, "Build metadata cannot be added to a release version, i.e. a pre-release version must be always specified\n")
		fmt.Fprintf(os.Stderr, "Command will fail when attempting to set a version that is smaller than the current version\n")
		fmt.Fprintf(os.Stderr, "Using \"version increase\" will bump the repository in pwd by a patch tick\n\n")

  case "":
    fmt.Fprintf(os.Stderr, "\nUsage of %s:\n\n", b.Sprint("version"))
		fmt.Fprintf(os.Stderr, "version [--root=\"\"] [--all]\n\n")
		fmt.Fprintf(os.Stderr, "The arguments are:\n\n")
    fmt.Fprintf(os.Stderr, fmt.Sprintf("\t%s\t%s", out("--root"), "root path from which to start listing repositories and their versions\n"))
		fmt.Fprintf(os.Stderr, fmt.Sprintf("\t%s\t%s", out("--all"), "list all versions\n"))
		fmt.Fprintf(os.Stderr, "\n")
    fmt.Fprintf(os.Stderr, "Using \"version\" will display the current version of the repository in pwd\n")
    fmt.Fprintf(os.Stderr, "Specifying a directory will recursively display the version(s) of the repositories contained there\n\n")

	default:
		fmt.Fprintf(os.Stderr, "\nUnknown command '%s'\n", b.Sprint(cmd))
		flag.Usage()
	}

}
