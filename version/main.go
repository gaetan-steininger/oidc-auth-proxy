package version

import "fmt"

var (
	version  string
	revision string
)

func PrintVersion(json bool) {
	if json {
		fmt.Printf("{\"version\": \"%s\", \"revision\": \"%s\"}\n", version, revision)
	} else {
		fmt.Printf("Version : %s\nRevision : %s\n", version, revision)
	}
}
