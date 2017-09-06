package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
  "sort"
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
  v.Special = strings.Replace(v.Special,"-",".",-1)
  w.Special = strings.Replace(w.Special,"-",".",-1)

	// Pre-release versions have a lower precedence than the associated normal version.
	if v.Special == "" && w.Special != "" {
		return true
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
		}else if partsv[i] < partsw[i] {
      return false
    }
	}

	// Should not be reached
	return false

}

// getLastRepoVersion returns the highest version from commited tags
func GetVersions(dir string) (*Versions, error) {

	// Change dir to repo root
	if err := os.Chdir(dir); err != nil {
		return nil, fmt.Errorf("could not change path to '%s': %s", dir, err.Error())
	}

	// Get all tags
	cmd := exec.Command("git", "log", "--tags", `--pretty="%h\t%at\t%D"`)
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("could not list versions: %s", err.Error())
	}

	// Find newest version
	re := regexp.MustCompile(V_REGEX)
	versions := &Versions{versions: []*Version{}}

MainLoop:
	for _, line := range strings.Split(string(out), "\n") {

    // Remove quotes
    line = strings.Replace(line,`"`,"",-1)

		// Verify correct output
		vparts := strings.Split(line, `\t`)
		if len(vparts) != 3 {
			continue
		}

		// Parse UNIX timestamp
		tint, err := strconv.ParseInt(vparts[1], 10, 64)
		if err != nil {
			panic(err)
		}
		timestamp := time.Unix(tint, 0)

		// Initiate new version
		v := &Version{
			Commit: vparts[0],
			Date:   timestamp,
		}

		// Match and parse version fields
		match := re.FindStringSubmatch(vparts[2])
		for i, name := range re.SubexpNames() {

			if i == 0 || len(match) < i+1 {
				continue
			}

			// Attempt field conversion to int
			vi, erri := strconv.Atoi(match[i])

			// Fill the version struct
			switch name {

			case "major":
				if erri != nil {
					continue MainLoop
				}
				v.Major = vi

			case "minor":
				if erri != nil {
					continue MainLoop
				}
				v.Minor = vi

			case "patch":
				if erri != nil {
					continue MainLoop
				}
				v.Patch = vi

			case "special":
				v.Special = match[i]

			case "build":
				v.Build = match[i]

			}

		}

		// Append completed version
		versions.Add(v)

	}

  // Sort with newest version being first
  sort.Sort(sort.Reverse(versions))

	return versions, nil

}

// increase increases repository's semantic version
func increase(major, minor, patch bool, special, build string) error {
	if major && minor || major && patch || minor && patch {
		return fmt.Errorf("cannot incrase more than one level: choose major, minor or patch")
	}

	return nil
}

// list lists all version of all repositories starting with root path
func list(root string, all bool) {

}