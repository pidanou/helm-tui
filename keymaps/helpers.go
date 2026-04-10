package keymaps

import (
	"slices"
	"strings"
)

func Contains(keys []string, target string) bool {
	return slices.Contains(keys, target)
}

func formatHelp(keymaps []string) string {
	return strings.Join(keymaps, "/")
}
