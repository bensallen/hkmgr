package up

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/bensallen/hkmgr/internal/config"
)

var hyperkitPath = "hyperkit"

func Run(cfg *config.Config, debug bool) error {
	for _, vm := range cfg.VM {
		fmt.Printf("Booting VM: %s\n", vm.UUID)
		cmdArgs := vm.Cli()

		if debug {
			fmt.Printf("cmd: %s %s\n", hyperkitPath, strings.Join(cmdArgs, " "))
		}

		cmd := exec.Command(hyperkitPath, cmdArgs...)
		err := cmd.Start()
		if err != nil {
			fmt.Printf("%v\n", err)
			continue
		}
		fmt.Printf("pid: %d\n", cmd.Process.Pid)
	}
	return nil
}
