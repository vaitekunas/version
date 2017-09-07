package main

import (
	"flag"
	"os"
	"os/exec"
	"strings"
)

func init() {
	_, err := exec.LookPath("git")
	if err != nil {
		print("Git not found. Cannot proceed, sorry")
		os.Exit(1)
	}
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
	specialPtr := incCmd.String("special","", "set pre-release version ")
	buildPtr := incCmd.String("build","", "set build metadata")

	// Version list flags
	listRootPtr := listCmd.String("root", "", "root path where listing version should start")
  listallPtr := listCmd.Bool("all", false, "show all versions")

	// General flags
	allPtr := flag.Bool("all", false, "show all versions")

	// Plain version
	if len(os.Args) == 1 {
    root, err := os.Getwd()
    if err != nil {
      print("could not determine version: %s", err.Error())
      os.Exit(1)
    }
		versions, err := GetVersions(root)
		if err != nil {
			print("could not determine version: %s", err.Error())
			os.Exit(1)
		}

    printVersionTable([]string{root}, map[string]*Versions{root: versions}, !*allPtr)
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
    flag.CommandLine.Parse(os.Args[2:])

		if f, err := os.Stat(os.Args[1]); os.IsNotExist(err) || (f != nil && !f.IsDir()) {
			flag.Usage()
			os.Exit(1)
		}

		versions, err := GetVersions(os.Args[1])
		if err != nil {
			print("could not determine version: %s", err.Error())
			os.Exit(1)
		}
    printVersionTable([]string{os.Args[1]}, map[string]*Versions{os.Args[1]: versions}, !*allPtr)
    os.Exit(0)
	}

	// Increase version
	if incCmd.Parsed() {
		if err := increase(*majorPtr, *minorPtr, *patchPtr, *specialPtr, *buildPtr); err != nil {
			print("could not increase version: %s", err.Error())
		}
	}

	// List versions
	if listCmd.Parsed() {
		list(strings.TrimRight(*listRootPtr,"/"), *listallPtr)
	}

}
