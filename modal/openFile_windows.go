//go:build windows
// +build windows

package modal

import (
	"fmt"
	s "github.com/VageLO/money-tree/structs"
	"os/exec"
	"syscall"
)

func OpenFiles(filePath string, source *s.Source) {
	defer ErrorModal(source.Pages, source.Modal)

	cmd := exec.Command("cmd")
	cmd.SysProcAttr = &syscall.SysProcAttr{CmdLine: fmt.Sprintf(`/c start "" "%s"`, filePath)}
	check(cmd.Run())
}
