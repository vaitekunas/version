package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

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
	versions := &Versions{versions: []*Version{}}
	for _, line := range strings.Split(string(out), "\n") {

		// Remove quotes
		line = strings.TrimSpace(line)
		line = strings.Trim(line, `"`)

		// Verify correct output
		vparts := strings.Split(line, `\t`)
		if len(vparts) != 3 {
			continue
		}

		// Extract version
		v, err := ExtractVersion(vparts[0], vparts[1], vparts[2])
		if err != nil {
			continue
		}

		// Append completed version
		versions.Add(v)

	}

	// Sort with newest version being first
	sort.Sort(sort.Reverse(versions))

	return versions, nil

}

// ExtractVersion extracts the version from git output
func ExtractVersion(commitPart, timePart, tagPart string) (*Version, error) {

	// Parse UNIX timestamp
	tint, err := strconv.ParseInt(timePart, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("could not parse string to UNIX timestamp")
	}
	timestamp := time.Unix(tint, 0)

	// Initiate new version
	v := &Version{
		Commit: commitPart,
		Date:   timestamp,
	}

	// Match and parse version fields
	re := regexp.MustCompile(V_REGEX)
	match := re.FindStringSubmatch(tagPart)
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
				return nil, fmt.Errorf("error parsing major tick")
			}
			v.Major = vi

		case "minor":
			if erri != nil {
				return nil, fmt.Errorf("error parsing minor tick")
			}
			v.Minor = vi

		case "patch":
			if erri != nil {
				return nil, fmt.Errorf("error parsing patch tick")
			}
			v.Patch = vi

		case "special":
			v.Special = match[i]

		case "build":
			v.Build = match[i]

		}
	}

	return v, nil
}

// GetLastCommit returns current commit
func GetLastCommit(root string) (date time.Time, commit, author, message string, version *Version, err error) {

	// Change dir to repo root
	if err := os.Chdir(root); err != nil {
		return time.Now(), "", "", "", nil, fmt.Errorf("could not change path to '%s': %s", root, err.Error())
	}

	// Get last log
	cmd := exec.Command("git", "log", "-1", `--pretty="%H\t%at\t%an\t%d\t%s"`)
	out, err := cmd.Output()
	if err != nil {
		return time.Now(), "", "", "", nil, fmt.Errorf("could not get last commit: %s", err.Error())
	}

	// Cleanup
	outStr := strings.TrimSpace(string(out))
	outStr = strings.Trim(outStr, `"`)
	parts := strings.Split(outStr, `\t`)

	if len(parts) != 5 {
		return time.Now(), "", "", "", nil, fmt.Errorf("invalid git output")
	}

	// Extract version
	v, err := ExtractVersion(parts[0], parts[1], parts[3])
	if err != nil {
		return time.Now(), "", "", "", nil, fmt.Errorf("could not extract version")
	}

	return v.Date, parts[0], parts[2], parts[4], v, nil
}

// GetBranch returns current branch
func GetBranch(root string) (string, error) {

		// Change dir to repo root
		if err := os.Chdir(root); err != nil {
			return "", fmt.Errorf("could not change path to '%s': %s", root, err.Error())
		}

		// Get last log
		cmd := exec.Command("git", "branch","--all")
		out, err := cmd.Output()
		if err != nil {
			return "", fmt.Errorf("could not get branches: %s", err.Error())
		}

		// Find active branch
		for _, line := range strings.Split(string(out),"\n") {
			line = strings.TrimSpace(line)
			line = strings.Trim(line, `"`)
			if strings.HasPrefix(line, "*") {
				if len(line) < 3 {
					return "", fmt.Errorf("invalid branch name")
				}
				return line[2:], nil
			}
		}

		return "", fmt.Errorf("could not determine active branch")
}
