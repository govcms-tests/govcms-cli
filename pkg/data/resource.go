package data

import (
	"strings"
)

type Resource int

const (
	DISTRIBUTION Resource = iota
	SAAS
	PAAS
)

type Installation struct {
	Name     string
	Path     string
	Resource Resource
}

var (
	resourceMap = map[string]Resource{
		"distribution": DISTRIBUTION,
		"saas":         SAAS,
		"paas":         PAAS,
	}
)

func StringToResource(str string) (Resource, bool) {
	resource, ok := resourceMap[strings.ToLower(str)]
	return resource, ok
}
