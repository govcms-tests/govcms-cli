package config

import "sort"

// Define a type to represent your configuration
type ReposList map[string]string

// Initialize the govcms repositories list with the initial values
var GovCMSReposList = ReposList{
	"distribution":     "govCMS/GovCMS",
	"saas":             "govCMS/scaffold",
	"paas":             "govCMS/scaffold",
	"lagoon":           "govCMS/lagoon",
	"scaffold-tooling": "govCMS/scaffold-tooling",
	"tests":            "govcms-tests/tests",
}

func init() {
	// Sort the keys alphabetically
	keys := make([]string, 0, len(GovCMSReposList))
	for key := range GovCMSReposList {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// Rebuild the map in sorted order
	sortedRepos := make(ReposList)
	for _, key := range keys {
		sortedRepos[key] = GovCMSReposList[key]
	}
	GovCMSReposList = sortedRepos
}
