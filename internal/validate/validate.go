package validate

import (
	"fmt"

	"github.com/bensallen/hkmgr/internal/config"
	"github.com/kr/pretty"
)

func Run(cfg *config.Config) error {

	fmt.Printf("%# v\n", pretty.Formatter(cfg))

	return nil

}
