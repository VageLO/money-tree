// +build linux

package modal

import (
	s "github.com/VageLO/money-tree/structs"
	"os/exec"
)

func OpenFiles(filePath string, source *s.Source) {
	defer ErrorModal(source.Pages, source.Modal)

    err := exec.Command("xdg-open", filePath).Start()
    check(err)
}
