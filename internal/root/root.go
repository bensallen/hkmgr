package root

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/bensallen/hkmgr/internal/config"
	"github.com/bensallen/hkmgr/internal/console"
	"github.com/bensallen/hkmgr/internal/destroy"
	"github.com/bensallen/hkmgr/internal/ssh"
	start "github.com/bensallen/hkmgr/internal/start"
	"github.com/bensallen/hkmgr/internal/status"
	stop "github.com/bensallen/hkmgr/internal/stop"
	"github.com/integrii/flaggy"
	"github.com/kr/pretty"
)

var version = "unknown"

func Run() error {
	var startSubcommand *flaggy.Subcommand
	var stopSubcommand *flaggy.Subcommand
	var destroySubcommand *flaggy.Subcommand
	var validateSubcommand *flaggy.Subcommand
	var sshSubcommand *flaggy.Subcommand
	var statusSubcommand *flaggy.Subcommand
	var consoleSubcommand *flaggy.Subcommand

	configPath := "hkmgr.toml"
	var debug bool
	var dryRun bool
	var vmName string

	flaggy.SetName("hkmgr")
	flaggy.SetDescription("VM manager for hyperkit")

	flaggy.DefaultParser.AdditionalHelpPrepend = "http://github.com/bensallen/hkmgr"

	flaggy.String(&configPath, "c", "config", "Path to configuration TOML file")
	flaggy.Bool(&debug, "d", "debug", "Enable debug output")
	flaggy.Bool(&dryRun, "n", "dry-run", "Don't execute any commands that affect change, just show what will be run")

	startSubcommand = flaggy.NewSubcommand("start")
	startSubcommand.Description = "Start VMs"
	startSubcommand.AddPositionalValue(&vmName, "name", 1, false, "Specify a VM, otherwise all VMs will be run")

	stopSubcommand = flaggy.NewSubcommand("stop")
	stopSubcommand.Description = "Stop VMs"
	var stopSignal string
	stopSubcommand.String(&stopSignal, "s", "signal", "Signal to send to VM")
	stopSubcommand.AddPositionalValue(&vmName, "name", 1, false, "Specify a VM, otherwise all VMs will be stopped")

	destroySubcommand = flaggy.NewSubcommand("destroy")
	destroySubcommand.Description = "Destroy VMs"
	destroySubcommand.AddPositionalValue(&vmName, "name", 1, false, "Specify a VM, otherwise all VMs will be destroyed!")

	validateSubcommand = flaggy.NewSubcommand("validate")
	validateSubcommand.Description = "Validate configuration"
	validateSubcommand.AddPositionalValue(&vmName, "name", 1, false, "Specify a VM")

	statusSubcommand = flaggy.NewSubcommand("status")
	statusSubcommand.Description = "Display status of VMs"
	statusSubcommand.AddPositionalValue(&vmName, "name", 1, false, "Specify a VM")

	sshSubcommand = flaggy.NewSubcommand("ssh")
	sshSubcommand.Description = "SSH to VM"
	sshSubcommand.AddPositionalValue(&vmName, "name", 1, true, "Specify a VM")

	consoleSubcommand = flaggy.NewSubcommand("console")
	consoleSubcommand.Description = "Open Console of VM"
	consoleSubcommand.AddPositionalValue(&vmName, "name", 1, true, "Specify a VM")

	flaggy.AttachSubcommand(startSubcommand, 1)
	flaggy.AttachSubcommand(stopSubcommand, 1)
	//flaggy.AttachSubcommand(destroySubcommand, 1)
	//flaggy.AttachSubcommand(validateSubcommand, 1)
	flaggy.AttachSubcommand(statusSubcommand, 1)
	//flaggy.AttachSubcommand(sshSubcommand, 1)
	//flaggy.AttachSubcommand(consoleSubcommand, 1)

	flaggy.SetVersion(version)
	flaggy.Parse()

	if dryRun {
		debug = true
	}

	cfgStat, err := os.Stat(configPath)

	if os.IsNotExist(err) {
		return fmt.Errorf("configuration file %s not found", configPath)
	}
	if cfgStat.IsDir() {
		return fmt.Errorf("configuration file %s is a directory", configPath)
	}

	var config config.Config
	if _, err := toml.DecodeFile(configPath, &config); err != nil {
		return err
	}

	if config.Path, err = filepath.Abs(configPath); err != nil {
		return err
	}

	config.Defaults()
	config.UpdateRelativePaths()

	if debug {
		fmt.Printf("Parsed config:\n\n%# v\n", pretty.Formatter(config))
	}

	switch {
	case startSubcommand.Used:
		if err := start.Run(&config, vmName, debug, dryRun); err != nil {
			return err
		}
	case stopSubcommand.Used:
		if err := stop.Run(&config, vmName, stopSignal); err != nil {
			return err
		}
	case destroySubcommand.Used:
		if err := destroy.Run(&config); err != nil {
			return err
		}
	case statusSubcommand.Used:
		if err := status.Current(&config, vmName, debug); err != nil {
			return err
		}
	case sshSubcommand.Used:
		if err := ssh.Run(&config); err != nil {
			return err
		}
	case consoleSubcommand.Used:
		if err := console.Run(&config); err != nil {
			return err
		}
	default:
		flaggy.ShowHelpAndExit("")
	}

	return nil
}
