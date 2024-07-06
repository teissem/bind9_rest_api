package bindfile

import (
	"os/exec"
	"strings"
)

func ReloadBind(command string) error {
	commandSplit := strings.Split(command, " ")
	mainCommand, arguments := commandSplit[0], commandSplit[1:]
	return exec.Command(mainCommand, arguments...).Run()
}
