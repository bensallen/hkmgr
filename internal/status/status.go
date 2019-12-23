package status

import (
	"fmt"
	"sort"

	"github.com/bensallen/hkmgr/internal/config"
)

//Current prints the current status of the VM(s).
func Current(cfg *config.Config, name string, debug bool) error {
	if name != "" {
		if vm, ok := cfg.VM[name]; ok {
			fmt.Printf("%s status is %s, PID: %d\n", name, vm.Status(), vm.PID)
		} else {
			return fmt.Errorf("%s not found in the configuration", name)
		}
	} else {
		names := make([]string, 0, len(cfg.VM))
		for name := range cfg.VM {
			names = append(names, name)
		}
		sort.Strings(names)
		for _, n := range names {
			vm := cfg.VM[n]
			fmt.Printf("%s status is %s, PID: %d\n", n, vm.Status(), vm.PID)
		}
	}
	return nil
}
