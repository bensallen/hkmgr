package up

import (
	"fmt"
	"os"

	"github.com/bensallen/hkmgr/internal/config"
)

//Run ...
func Run(cfg *config.Config, debug bool, dryRun bool) error {
	for _, netTypes := range cfg.Network {
		net := netTypes.NetType()
		if err := net.Discover(); err != nil {
			return err
		}
	}

	for _, vm := range cfg.VM {
		fmt.Printf("Booting VM: %s\n", vm.UUID)

		if err := os.MkdirAll(vm.RunDir, os.ModePerm); err != nil {
			fmt.Printf("%v\n", err)
			continue
		}

		if !dryRun {
			err := vm.Up()
			if err != nil {
				fmt.Printf("%v\n", err)
				continue
			}
			for _, vmnet := range vm.Network {
				if vmnet.Device != "" {
					if net, ok := cfg.Network[vmnet.MemberOf]; ok {
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
