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
	fmt.Fprintf(os.Stderr, "version command [arguments]\n\n")
	fmt.Fprintf(os.Stderr, "The commands are:\n\n")
	fmt.Fprintf(os.Stderr, fmt.Sprintf("\t%s\t%s", out("increase"), "increases the version by a major/minor/patch tick\n"))
	fmt.Fprintf(os.Stderr, fmt.Sprintf("\t%s\t%s", out("list"), "lists versions of all repositories starting from some root path\n"))
	fmt.Fprintf(os.Stderr, "\n")
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
		fmt.Fprintf(os.Stderr, "version increase [--repo=\"\"] [--golang] [{--major, --minor, --patch}]\n\n")
		fmt.Fprintf(os.Stderr, "The arguments are:\n\n")
		fmt.Fprintf(os.Stderr, fmt.Sprintf("\t%s\t%s", out("--repo"), "path to the relevant repo, or a golang repo name\n"))
		fmt.Fprintf(os.Stderr, fmt.Sprintf("\t%s\t%s", out("--golang"), "this is a golang repo, i.e. the true path is $GOPATH/src/$repo\n"))
		fmt.Fprintf(os.Stderr, fmt.Sprintf("\t%s\t%s", out("--major"), "increase version by a major tick\n"))
		fmt.Fprintf(os.Stderr, fmt.Sprintf("\t%s\t%s", out("--minor"), "increase version by a minor tick\n"))
		fmt.Fprintf(os.Stderr, fmt.Sprintf("\t%s\t%s", out("--patch"), "increase version by a patch tick\n"))
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Only a single tick option (major/minor/patch) is allowed per increase\n")
		fmt.Fprintf(os.Stderr, "Command will fail if current version in .version has not been commited in the master branch\n")
		fmt.Fprintf(os.Stderr, "Using \"version increase\" will bump the repository in pwd by a patch tick\n\n")

	case "list":
		fmt.Fprintf(os.Stderr, "\nUsage of %s:\n\n", b.Sprint("version list"))
		fmt.Fprintf(os.Stderr, "version list [--root=\"\"] [--all]\n\n")
		fmt.Fprintf(os.Stderr, "The arguments are:\n\n")
		fmt.Fprintf(os.Stderr, fmt.Sprintf("\t%s\t%s", out("--root"), "root directory from which to start listing repository versions\n"))
		fmt.Fprintf(os.Stderr, fmt.Sprintf("\t%s\t%s", out("--all"), "list all versions of each repo\n"))
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, "Using \"version list\" will use pwd as the root directory\n\n")

	default:
		fmt.Fprintf(os.Stderr, "\nUnknown command '%s'\n", b.Sprint(cmd))
		flag.Usage()
	}

}
