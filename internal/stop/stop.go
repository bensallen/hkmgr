package stop

import (
	"fmt"

	"github.com/bensallen/hkmgr/internal/config"
)

// Run stops all VMs or the specific VMs passed as "name".
func Run(cfg *config.Config, name string, signal string) error {
	if name != "" {
		if vm, ok := cfg.VM[name]; ok {
			err := vm.Down(signal)
			if err != nil {
				fmt.Printf("Stopping VM failed, %v\n", err)
			}
		} else {
			return fmt.Errorf("%s not found in the configuration", name)
		}
	} else {
		for name, vm := range cfg.VM {
			err := vm.Down(signal)
			if err != nil {
				fmt.Printf("Stopping VM %s failed, %v\n", name, err)
			}
		}
	}
	return nil
}
