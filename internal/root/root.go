package root

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/bensallen/hkmgr/internal/config"
	"github.com/bensallen/hkmgr/internal/console"
	"github.com/bensallen/hkmgr/internal/destroy"
	"github.com/bensallen/hkmgr/internal/down"
	"github.com/bensallen/hkmgr/internal/ssh"
	"github.com/bensallen/hkmgr/internal/status"
	"github.com/bensallen/hkmgr/internal/up"
	"github.com/integrii/flaggy"
	"github.com/kr/pretty"
)

// Version is the CLI version or release number. To be overriden at build time.
var Version = "unknown"

// Run hkmgr
func Run() error {
	var upSubcommand *flaggy.Subcommand
	var downSubcommand *flaggy.Subcommand
	var destroySubcommand *flaggy.Subcommand
	var validateSubcommand *flaggy.Subcommand
	var sshSubcommand *flaggy.Subcommand
	var statusSubcommand *flaggy.Subcommand
	var consoleSubcommand *flaggy.Subcommand

	//
	var cliConfigPaths []string
	var debug bool
	var dryRun bool
	var vmName string

	flaggy.SetName("hkmgr")
	flaggy.SetDescription("VM manager for hyperkit")

	flaggy.DefaultParser.AdditionalHelpPrepend = "http://github.com/bensallen/hkmgr"
	flaggy.StringSlice(&cliConfigPaths, "c", "config", "Path to configuration TOML file or directory of TOML files")

	flaggy.Bool(&debug, "d", "debug", "Enable debug output")
	flaggy.Bool(&dryRun, "n", "dry-run", "Don't execute any commands that affect change, just show what will be run")

	upSubcommand = flaggy.NewSubcommand("up")
	upSubcommand.Description = "Start VMs"
	upSubcommand.ShortName = "start"
	upSubcommand.AddPositionalValue(&vmName, "name", 1, false, "Specify a VM, otherwise all VMs will be run")

	downSubcommand = flaggy.NewSubcommand("down")
	downSubcommand.Description = "Stop VMs"
	downSubcommand.ShortName = "stop"
	var downSignal string
	downSubcommand.String(&downSignal, "s", "signal", "Signal to send to VM")
	downSubcommand.AddPositionalValue(&vmName, "name", 1, false, "Specify a VM, otherwise all VMs will be stopped")

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

	flaggy.AttachSubcommand(upSubcommand, 1)
	flaggy.AttachSubcommand(downSubcommand, 1)
	//flaggy.AttachSubcommand(destroySubcommand, 1)
	//flaggy.AttachSubcommand(validateSubcommand, 1)
	flaggy.AttachSubcommand(statusSubcommand, 1)
	//flaggy.AttachSubcommand(sshSubcommand, 1)
	//flaggy.AttachSubcommand(consoleSubcommand, 1)

	flaggy.SetVersion(Version)
	flaggy.Parse()

	if dryRun {
		debug = true
	}

	// Default to look for configs hkmgr.toml and hkmgr.d/*.toml
	if len(cliConfigPaths) == 0 {
		cliConfigPaths = []string{"hkmgr.toml", "hkmgr.d"}
	}

	cfgPaths := []string{}
	var firstCfgPath string
	for _, cliConfigPath := range cliConfigPaths {
		cfgStat, err := os.Stat(cliConfigPath)

		if os.IsNotExist(err) {
			continue
		}
		if cfgStat.IsDir() {
			globPaths, err := filepath.Glob(cliConfigPath + "/*.toml")

			if err != nil {
				fmt.Printf("configuration path %s is a directory and failed with %v", cliConfigPath, err)
				continue
			}

			if len(globPaths) == 0 {
				continue
			}
			cfgPaths = append(cfgPaths, globPaths...)
		} else {
			cfgPaths = append(cfgPaths, cliConfigPath)
		}
		// Record the first valid path so that it can be used to resolve relative paths
		if firstCfgPath == "" {
			firstCfgPath = cliConfigPath
		}
	}
	fmt.Printf("%#v\n", cfgPaths)

	if len(cfgPaths) == 0 {
		return fmt.Errorf("no configuration files found")
	}

	var config config.Config
	for _, cfgPath := range cfgPaths {
		if _, err := toml.DecodeFile(cfgPath, &config); err != nil {
			return err
		}
	}

	absPath, err := filepath.Abs(firstCfgPath)
	if err != nil {
		return err
	}
	config.Path = absPath

	if err := config.Defaults(); err != nil {
		return err
	}
	config.UpdateRelativePaths()

	if debug {
		fmt.Printf("Parsed config:\n\n%# v\n", pretty.Formatter(config))
	}

	switch {
	case upSubcommand.Used:
		if err := up.Run(&config, vmName, debug, dryRun); err != nil {
			return err
		}
	case downSubcommand.Used:
		if err := down.Run(&config, vmName, downSignal); err != nil {
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
