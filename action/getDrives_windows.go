//go:build windows
// +build windows

package action

import (
	m "github.com/VageLO/money-tree/modal"
	s "github.com/VageLO/money-tree/structs"
	"strings"
	"unicode/utf16"

	"golang.org/x/sys/windows"
)

func getDrives(source *s.Source) []string {
	defer m.ErrorModal(source.Pages, source.Modal)

	n, err := windows.GetLogicalDriveStrings(0, nil)
	check(err)

	a := make([]uint16, n)
	windows.GetLogicalDriveStrings(n, &a[0])
	s := string(utf16.Decode(a))
	return strings.Split(strings.TrimRight(s, "\x00"), "\x00")
}
