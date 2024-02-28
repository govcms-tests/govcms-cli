package data

import (
	"strings"
)

type Resource int

const (
	DISTRIBUTION Resource = iota
	SAAS
	PAAS
	LAGOON
	TESTS
	SCAFFOLD_TOOLING
)

type Installation struct {
	Name string
	Path string
	Type string
}

var (
	resourceMap = map[string]Resource{
		"distribution":     DISTRIBUTION,
		"saas":             SAAS,
		"paas":             PAAS,
		"lagoon":           LAGOON,
		"tests":            TESTS,
		"scaffold":         SCAFFOLD_TOOLING,
		"scaffold-tooling": SCAFFOLD_TOOLING,
	}
)

func StringToResource(str string) (Resource, bool) {
	resource, ok := resourceMap[strings.ToLower(str)]
	return resource, ok
}
