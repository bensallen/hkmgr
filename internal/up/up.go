package up

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/bensallen/hkmgr/internal/config"
)

var hyperkitPath = "hyperkit"

func Run(cfg *config.Config, debug bool, dryRun bool) error {

	for name, netTypes := range cfg.Network {
		fmt.Printf("Configuring Network: %#v\n", name)
		if !dryRun {
			net := netTypes.NetType()
			if err := net.Up(); err != nil {
				return err
			}
		}
	}
	for _, vm := range cfg.VM {
		fmt.Printf("Booting VM: %s\n", vm.UUID)
		cmdArgs := vm.Cli()

		if debug {
			fmt.Printf("cmd: %s %s\n", hyperkitPath, strings.Join(cmdArgs, " "))
		}

		if !dryRun {
			cmd := exec.Command(hyperkitPath, cmdArgs...)
			err := cmd.Start()
			if err != nil {
				fmt.Printf("%v\n", err)
				continue
			}
			fmt.Printf("pid: %d\n", cmd.Process.Pid)
		}
	}
	return nil
}
