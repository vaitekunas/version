package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
)

const (
	V_REGEX = `v(?P<major>\d+).(?P<minor>\d+).(?P<patch>\d+)(-(?P<special>[a-z0-9\.-]+)(\+(?P<build>[a-z0-9\.-]+))?)?`
)

// Versions implements the sort.Interface
type Versions struct {
	versions []*Version
}

// Add adds a new version to the slice of versions
func (v *Versions) Add(vnew *Version) {
	v.versions = append(v.versions, vnew)
}

// Swap implements sort.Interface.Len
func (v *Versions) Len() int {
	return len(v.versions)
}

// Swap implements sort.Interface.Less
func (v *Versions) Less(i, j int) bool {
	return Larger(v.versions[j], v.versions[i])
}

// Swap implements sort.Interface.Swap
func (v *Versions) Swap(i, j int) {
	temp := v.versions[i]
	v.versions[i] = v.versions[j]
	v.versions[j] = temp
}

// Version holds all the relevant details on a semantic version
type Version struct {
	Date                time.Time
	Commit              string
	Major, Minor, Patch int
	Special             string
	Build               string
}

// String outputs a string version
func (v *Version) String() string {
	str := fmt.Sprintf("v%d.%d.%d", v.Major, v.Minor, v.Patch)
	if v.Special != "" && v.Build != "" {
		str = fmt.Sprintf("%s-%s+%s", str, v.Special, v.Build)
	} else if v.Special != "" {
		str = fmt.Sprintf("%s-%s", str, v.Special)
	}

	return str
}

// Larger compares version v to version w and returns true if v is larger.
// Uses comparison rules described in http://semver.org/
func Larger(v, w *Version) bool {

	// Compare release version
	if v.Major > w.Major {
		return true
	}
	if v.Major < w.Major {
		return false
	}
	if v.Major == w.Major && v.Minor > w.Minor {
		return true
	}
	if v.Major == w.Major && v.Minor < w.Minor {
		return false
	}
	if v.Major == w.Major && v.Minor == w.Minor && v.Patch > w.Patch {
		return true
	}
	if v.Major == w.Major && v.Minor == w.Minor && v.Patch < w.Patch {
		return false
	}

	// Replace hyphens
	v.Special = strings.Replace(v.Special, "-", ".", -1)
	w.Special = strings.Replace(w.Special, "-", ".", -1)

	// Pre-release versions have a lower precedence than the associated normal version.
	if v.Special == "" && w.Special != "" {
		return true
	} else if v.Special != "" && w.Special == "" {
		return false
	}

	// Two versions that differ only in the build metadata, have the same precedence.
	// Deviation from semver rules: commit date decides precedence
	if v.Special == w.Special {
		if v.Date.Unix() > w.Date.Unix() {
			return true
		}
		return false
	}

	// Identifiers of the special tick
	partsv := strings.Split(v.Special, ".")
	partsw := strings.Split(w.Special, ".")
	splen := len(partsv)
	if len(partsw) > splen {
		splen = len(partsw)
	}

	// Compare all special tick parts
	for i := 0; i <= splen-1; i++ {

		// A larger set of pre-release fields has a higher precedence than
		// a smaller set, if all of the preceding identifiers are equal
		if i > len(partsw)-1 {
			return true
		}
		if i > len(partsv)-1 {
			return false
		}

		// Numeric identifiers have lower precedence than non-numeric identifiers.
		vint, errv := strconv.Atoi(partsv[i])
		wint, errw := strconv.Atoi(partsw[i])
		if errv != nil && errw == nil {
			return true
		} else if errv == nil && errw != nil {
			return false
		}

		// Compare integers numerically
		if errv == nil && errw == nil {
			if vint > wint {
				return true
			}
			return false
		}

		// Compare strings lexicographically
		if partsv[i] > partsw[i] {
			return true
		} else if partsv[i] < partsw[i] {
			return false
		}
	}

	// Should not be reached
	return false

}

// Increase increases repository's semantic version
func Increase(major, minor, patch bool, special, build string) error {

	// Validate increment
	if major && minor || major && patch || minor && patch {
		return fmt.Errorf("cannot increase more than one level: choose major, minor or patch")
	}

	// Default increase is a patch tick
	if !major && !minor && !patch && special == "" {
		patch = true
	}

	// Get pwd
	root, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("could not determine current directory: %s", err.Error())
	}

	// Determine current version
	versions, err := GetVersions(root)
	if err != nil {
		return fmt.Errorf("could not determine version: %s", err.Error())
	}
	current := versions.versions[0]

	newVersion := &Version{
		Major:   current.Major,
		Minor:   current.Minor,
		Patch:   current.Patch,
		Special: special,
		Build:   build,
	}
	if major {
		newVersion.Major++
		newVersion.Minor = 0
		newVersion.Patch = 0
	}
	if minor {
		newVersion.Minor++
		newVersion.Patch = 0
	}
	if patch {
		newVersion.Patch++
	}

	// Validate
	if !Larger(newVersion, current) {
		printErr("cannot apply increase: proposed version (%s) is lower than the current version (%s)", newVersion.String(), current.String())
		os.Exit(1)
	}

	// Get last commit
	ctime, commit, author, message, version, err := GetLastCommit(root)
	if err != nil {
		return fmt.Errorf("could not get last commit: %s", err.Error())
	}
	if version.String() != "v0.0.0" {
		return fmt.Errorf("current commit already has a version: %s", version.String())
	}

	// Get branch
	branch, err := GetBranch(root)
	if err != nil {
		return fmt.Errorf("could not get active branch name: %s", err.Error())
	}

	// Formatting functions
	// TODO: put all of this in utils.go and unify outputs
	bold := color.New(color.Bold).Sprint
	abort := color.New(color.FgHiRed).Add(color.Bold).Sprint
	success := color.New(color.FgHiGreen).Add(color.Bold).Sprint
	bullet := func() string { return color.New(color.FgHiBlue).Sprint("◈") }
	out := func(s string, a ...interface{}) {
		if len(a) > 0 {
			s = fmt.Sprintf(s, a...)
		}
		fmt.Printf("\t %s  %s\n", bullet(), s)
	}

	// Print information
	fmt.Println("")

	fmt.Println("Repository:")
	out(getRepoName(root))
	fmt.Println("")

	fmt.Println("Commit to be tagged as the new version:")
	out("Branch:\t%s", bold(branch))
	out("Message:\t%s", bold(message))
	out("Hash:\t%s", bold(commit))
	out("Date:\t%s", bold(ctime.Format("2006-01-02 15:04:06")))
	out("Author:\t%s", bold(author))
	fmt.Println("")

	fmt.Println("Version increment:")
	if current.String() != "v0.0.0" {
		out("Current version: %s", bold(current.String()))
	} else {
		out("Current version: %s", bold("none"))
	}
	out("Proposed version after increase: %s", bold(newVersion.String()))

	fmt.Println("")
	fmt.Printf("%s", bold("Tag new version? [Y/n] (default: n): "))
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	if text != "Y\n" {
		fmt.Println(abort("\nVersion update aborted\n"))
		return nil
	}

	// Apply tag
	if err := exec.Command("git", "tag", "-a", newVersion.String(), "-m", fmt.Sprintf(`"Version %s"`, newVersion.String()), commit).Run(); err != nil {
		return fmt.Errorf("could not apply tag: %s", err.Error())
	}

	fmt.Println(success("\nVersion updated\n"))

	return nil
}

// List lists all version of all repositories starting with root path
func List(root string, all bool) error {

	if root == "" {
		dir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("could not change directory: %s", err.Error())
		}
		root = dir
	} else if f, err := os.Stat(root); (err != nil && os.IsNotExist(err)) || (err == nil && !f.IsDir()) {
		return fmt.Errorf("provided root path is not a directory")
	}

	repos := []string{}
	repoVersions := map[string]*Versions{}

	var scan func(string)
	scan = func(dir string) {
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			return
		}

		for _, file := range files {
			name := file.Name()
			if file.IsDir() {
				if name == ".git" {
					v, errv := GetVersions(dir)
					if errv != nil {
						continue
					}
					repos = append(repos, dir)
					repoVersions[dir] = v
				} else if name[:1] != "." {
					scan(fmt.Sprintf("%s/%s", dir, name))
				}
			}
		}
	}

	scan(root)

	sort.Strings(repos)
	printVersionTable(repos, repoVersions, !all)

	return nil
}
