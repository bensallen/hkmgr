package main

import (
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/integrii/flaggy"
	"github.com/kr/pretty"
)

var version = "unknown"

var upSubcommand *flaggy.Subcommand
var downSubcommand *flaggy.Subcommand
var destroySubcommand *flaggy.Subcommand
var sshSubcommand *flaggy.Subcommand
var consoleSubcommand *flaggy.Subcommand

func init() {
	flaggy.SetName("hkmgr")
	flaggy.SetDescription("VM manager for hyperkit")

	flaggy.DefaultParser.AdditionalHelpPrepend = "http://github.com/bensallen/hkmgr"

	upSubcommand = flaggy.NewSubcommand("up")
	upSubcommand.Description = "Start VMs"

	downSubcommand = flaggy.NewSubcommand("down")
	downSubcommand.Description = "Stop VMs"

	destroySubcommand = flaggy.NewSubcommand("destroy")
	destroySubcommand.Description = "Destroy VMs"

	sshSubcommand = flaggy.NewSubcommand("ssh")
	sshSubcommand.Description = "SSH to VM"

	consoleSubcommand = flaggy.NewSubcommand("console")
	consoleSubcommand.Description = "Open Console of VM"

	flaggy.AttachSubcommand(upSubcommand, 1)
	flaggy.AttachSubcommand(downSubcommand, 1)
	flaggy.AttachSubcommand(destroySubcommand, 1)
	flaggy.AttachSubcommand(sshSubcommand, 1)
	flaggy.AttachSubcommand(consoleSubcommand, 1)

	// set the version and parse all inputs into variables
	flaggy.SetVersion(version)
	flaggy.Parse()
}

func main() {

	var config Config
	if _, err := toml.DecodeFile("example_cfg.toml", &config); err != nil {
		// handle error
	}

	fmt.Printf("%# v\n", pretty.Formatter(config))
}
