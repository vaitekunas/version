package main


import (
  "fmt"
)

// getVersion returns repository's current version
func getVersion(repo string) (string, error) {
  return "", fmt.Errorf("missing file '.version'")
}

// increase increases repository's semantic version
func increase(repo string, golang, major, minor, patch bool) error {
  if major && minor || major && patch || minor && patch {
    return fmt.Errorf("cannot incrase more than one level: choose major, minor or patch")
  }

  return nil
}

// list lists all version of all repositories starting with root path
func list(root string, all bool) {

}
