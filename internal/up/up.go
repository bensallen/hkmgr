package up

import (
	"fmt"
	"os"

	"github.com/bensallen/hkmgr/internal/config"
)

//Run ...
func Run(cfg *config.Config, vmName string, debug bool, dryRun bool) error {
	for _, netTypes := range cfg.Network {
		net := netTypes.NetType()
		if err := net.Discover(); err != nil {
			return err
		}
	}

	if vmName != "" {
		if vm, ok := cfg.VM[vmName]; ok {
			if err := upVM(vm, cfg.Network, dryRun); err != nil {
				return fmt.Errorf("error bringing vm: %s up, %v", vm.UUID, err)
			}
		} else {
			return fmt.Errorf("VM %s not found in the configuration", vmName)
		}
	} else {
		for _, vm := range cfg.VM {
			if err := upVM(vm, cfg.Network, dryRun); err != nil {
				fmt.Printf("Error bringing vm: %s up, %v\n", vm.UUID, err)
			}
		}
	}

	for name, netTypes := range cfg.Network {
		fmt.Printf("Configuring Network: %#v\n", name)
		if !dryRun {
			net := netTypes.NetType()
			if err := net.Up(); err != nil {
				return err
			}
		}
	}

	return nil
}

func upVM(vm config.VMConfig, netCfg config.Network, dryRun bool) error {
	fmt.Printf("Booting VM: %s\n", vm.UUID)

	if err := os.MkdirAll(vm.RunDir, os.ModePerm); err != nil {
		return err
	}

	if !dryRun {
		if err := vm.Up(); err != nil {
			return err
		}

		for _, vmnet := range vm.Network {
			if vmnet.Device != "" {
				if net, ok := netCfg[vmnet.MemberOf]; ok {
					if tap := net.Tap; tap != nil {
						if tap.BridgeDev != nil {
							fmt.Printf("adding member %s to network %s for vm %s\n", vmnet.Device, vmnet.MemberOf, vm.UUID)
							tap.BridgeDev.Members = append(tap.BridgeDev.Members, vmnet.Device)
						}
					}
				} else {
					fmt.Printf("warning: could not find configured network %s for vm %s\n", vmnet.MemberOf, vm.UUID)
				}
			}
		}
	}
	return nil
}
