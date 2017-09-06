package main

import (
	"flag"
	"os"
	"strings"
)

func init() {
  flag.Usage = help
}

func main() {

	// Subcommands
	incCmd := flag.NewFlagSet("increase", flag.ExitOnError)
	listCmd := flag.NewFlagSet("list", flag.ExitOnError)

	// Version increase flags
	majorPtr := incCmd.Bool("major", false, "increase major version")
	minorPtr := incCmd.Bool("minor", false, "increase minor version")
	patchPtr := incCmd.Bool("patch", false, "increase patch version")
	goPtr := incCmd.Bool("golang", false, "golang repository")
	repoPtr := incCmd.String("repo", "", "path to repository (optional)")

	// Version list flags
	listRootPtr := listCmd.String("root", os.Getenv("GOPATH"), "root path where listing version should start")
  allPtr := incCmd.Bool("all", false, "show all versions")

  // Plain version
	if len(os.Args) == 1 {
		version, err := getVersion("")
		if err != nil {
			print("could not determine version: %s", err.Error())
			os.Exit(1)
		}
		print("current version: %s", version)
		os.Exit(0)
	}

	// Parse subcommand flags
	switch strings.ToLower(os.Args[1]) {

	case "increase":
		incCmd.Parse(os.Args[2:])

	case "list":
    listCmd.Parse(os.Args[2:])

	case "help":
		if len(os.Args) >= 3 {
			man(os.Args[2])
		} else {
			flag.Usage()
		}

	default: // Assuming "version path/to/repo" was used
		if f, err := os.Stat(os.Args[1]); os.IsNotExist(err) || (f != nil && !f.IsDir()) {
			flag.Usage()
			os.Exit(1)
		}

		version, err := getVersion(os.Args[1])
		if err != nil {
			print("could not determine version: %s", err.Error())
			os.Exit(1)
		}
		print("current version: %s", version)
		os.Exit(0)
	}

	// Increase version
	if incCmd.Parsed() {
		if err := increase(*repoPtr, *goPtr, *majorPtr, *minorPtr, *patchPtr); err != nil {
			print("could not increase version: %s", err.Error())
		}
	}

	// List versions
	if listCmd.Parsed() {
		list(*listRootPtr, *allPtr)
	}

}
