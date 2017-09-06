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
func printVersionTable(versions *Versions, last bool) {

  if versions.Len() == 0 {
    fmt.Println("No versions")
    return
  }

  if last {
    versions.versions = versions.versions[:1]
  }

  current := func(v interface{}) interface{}{
    return color.New(color.Bold).Sprint(v)
  }

  blue := func(v interface{}) interface{}{
    return color.New(color.FgHiBlue).Sprint(v)
  }

  table := lentele.New("Date", "Commit", "Version")
  if !last {
    table.AddTitle("Complete list of repository versions")
    table.AddTitle("(ordered from the highest to the lowest)")
  }else{
    table.AddTitle("Current repository version")
  }

  if header, err := table.GetRowByName("header"); err == nil {
    header.Modify(current,"Date", "Commit", "Version")
  }

  for i, version := range versions.versions {
    row := table.AddRow("")
    row.Insert(version.Date.Format("2006-01-02 15:04"), version.Commit, version.String())
    if i == 0 {
      row.Modify(blue,"Date", "Commit", "Version")
    }else{

    }
  }

  table.AddFootnote("Version order is based on the semantic versioning specification (http://semver.org/)")

  table.Render(os.Stdout,false, true, true, lentele.LoadTemplate("modern"))
  fmt.Printf("\n\n")

}
