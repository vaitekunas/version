package main

import (
	"fmt"
  "os"
	"github.com/fatih/color"
  "github.com/vaitekunas/lentele"
  "strings"
)

// print displays a message
func print(in string, a ...interface{}) {
	if len(a) > 0 {
		in = fmt.Sprintf(in, a...)
	}

	b := color.New(color.FgHiBlue)
	fmt.Printf(" %s  %s\n", b.Sprint("◈"), in)
}

func getRepoName(dir string) string {
  root := ""
  repo := ""

  if strings.Contains(dir,"github.com") {
    idx := strings.Index(dir,"github.com")
    root = dir[:idx]
    repo = dir[idx:]
  }else {
    idx := strings.LastIndex(strings.TrimRight(dir,"/"),"/")+1
    root = dir[:idx]
    repo = dir[idx:]
  }

  return fmt.Sprintf("%s%s",root,color.New(color.Bold).Sprint(repo))
}

// printVersionTable displays version data in a table
func printVersionTable(repos []string, repoVersions map[string]*Versions, last bool) {

  bold := func(v interface{}) interface{}{
    return color.New(color.Bold).Sprint(v)
  }

  blue := func(v interface{}) interface{}{
    return color.New(color.FgHiBlue).Add(color.Bold).Sprint(v)
  }

  repoName := func(v interface{}) interface{}{
    vs, ok := v.(string)
    if !ok {
      return v
    }
    repo := strings.TrimSpace(vs)
    return strings.Replace(vs,repo,getRepoName(repo),1)
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
    header.Modify(bold,"Repository", "Date", "Commit", "Version")
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
      row.Insert(alignedRepo, version.Date.Format("2006-01-02 15:04"), version.Commit, alignedVersion).Modify(repoName,"Repository")
      if i == 0 && !last {
        row.Modify(blue,"Version")
      }
    }
  }

  table.AddFootnote("Version order is based on the semantic versioning specification (http://semver.org/)")

  table.Render(os.Stdout,false, true, false, lentele.LoadTemplate("classic"))
  fmt.Printf("\n")


}
