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
		printErr("Git not found. Cannot proceed, sorry")
		os.Exit(1)
	}
	flag.Usage = help
}

func main() {

	// Subcommands
	incCmd := flag.NewFlagSet("increase", flag.ExitOnError)

	// Version increase flags
	majorPtr := incCmd.Bool("major", false, "increase major version")
	minorPtr := incCmd.Bool("minor", false, "increase minor version")
	patchPtr := incCmd.Bool("patch", false, "increase patch version")
	specialPtr := incCmd.String("special", "", "set pre-release version ")
	buildPtr := incCmd.String("build", "", "set build metadata")

	// Version list flags
	listRootPtr := flag.String("root", "", "root path where listing version should start")
	listallPtr := flag.Bool("all", false, "show all versions")

	// Parse subcommand flags
	if len(os.Args) > 1 {
		switch strings.ToLower(os.Args[1]) {

		case "increase":
			incCmd.Parse(os.Args[2:])

		case "help":
			if len(os.Args) >= 3 {
				man(os.Args[2])
			} else {				
        man("")
			}
			os.Exit(0)

    case "--all", "--root":

    default:
      flag.Usage()
      os.Exit(1)
		}

		// Increase version
		if incCmd.Parsed() {
			if err := Increase(*majorPtr, *minorPtr, *patchPtr, *specialPtr, *buildPtr); err != nil {
				printErr("FAILED: %s", err.Error())
			}
			os.Exit(0)
		}
	}

	// Parse global
	flag.Parse()

	// List versions
	if err := List(strings.TrimRight(*listRootPtr, "/"), *listallPtr); err != nil {
		printErr("FAILED: %s", err.Error())
	}

}
