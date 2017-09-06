package main

import (
	"fmt"
  "os"
	"github.com/fatih/color"
  "github.com/vaitekunas/lentele"
)

// print displays a message
func print(in string, a ...interface{}) {
	if len(a) > 0 {
		in = fmt.Sprintf(in, a...)
	}

	b := color.New(color.FgHiBlue)
	fmt.Printf(" %s  %s\n", b.Sprint("◈"), in)
}

// printVersionTable displays version data in a table
func printVersionTable(repos []string, repoVersions map[string]*Versions, last bool) {

  current := func(v interface{}) interface{}{
    return color.New(color.Bold).Sprint(v)
  }

  blue := func(v interface{}) interface{}{
    return color.New(color.FgHiBlue).Sprint(v)
  }

  table := lentele.New("Repository", "Date", "Commit", "Version")
  if !last {
    if len(repos) > 1 {
      table.AddTitle("All versions")
    }else{
      table.AddTitle("All versions per repository")
    }
    table.AddTitle("(ordered from the highest to the lowest)")
  }else{
    if len(repos) > 1 {
      table.AddTitle("Current versions per repository")
    }else{
      table.AddTitle("Current version")
    }
  }

  if header, err := table.GetRowByName("header"); err == nil {
    header.Modify(current,"Repository", "Date", "Commit", "Version")
  }

  // Repository path format
  longestRepo := 0
  for _, repo := range repos {
      if len(repo) > longestRepo {
        longestRepo = len(repo)
      }
  }

  longestVersion := 0
  for _,versions := range repoVersions {
    for _, version := range versions.versions {
      if vlen := len(version.String()); vlen > longestVersion {
        longestVersion = vlen
      }
    }
  }

  formatRepo := fmt.Sprintf("%%-%ds",longestRepo)
  formatVersion := fmt.Sprintf("%%-%ds",longestVersion)

  for _, repo := range repos {
    versions, ok := repoVersions[repo]
    if !ok {
      continue
    }
    if last {
      versions.versions = versions.versions[:1]
    }

    for i, version := range versions.versions {
      row := table.AddRow("")
      alignedRepo := fmt.Sprintf(formatRepo,repo)
      alignedVersion := fmt.Sprintf(formatVersion,version.String())
      if version.String() == "v0.0.0" {
        alignedVersion = "N/A"
      }
      row.Insert(alignedRepo, version.Date.Format("2006-01-02 15:04"), version.Commit, alignedVersion)
      if i == 0 && !last {
        row.Modify(blue,"Repository","Date", "Commit", "Version")
      }
    }
  }

  table.AddFootnote("Version order is based on the semantic versioning specification (http://semver.org/)")

  table.Render(os.Stdout,false, true, true, lentele.LoadTemplate("modern"))
  fmt.Printf("\n\n")


}
