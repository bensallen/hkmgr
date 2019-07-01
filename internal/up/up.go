package up

import (
	"fmt"

	"github.com/bensallen/hkmgr/internal/config"
	"github.com/kr/pretty"
)

func Run(cfg *config.Config) error {

	for _, vm := range cfg.VM {

		//fmt.Printf("%# v\n", pretty.Formatter(vm))
		fmt.Printf("%# v\n", pretty.Formatter(vm.Cli()))

	}

	return nil

}
