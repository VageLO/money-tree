//go:build linux
// +build linux

package action

import (
	s "github.com/VageLO/money-tree/structs"
)

func getDrives(source *s.Source) []string {
	return []string{"/"}
}
