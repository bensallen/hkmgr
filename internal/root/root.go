package root

import (
	"github.com/BurntSushi/toml"
	"github.com/bensallen/hkmgr/internal/config"
	"github.com/bensallen/hkmgr/internal/console"
	"github.com/bensallen/hkmgr/internal/destroy"
	"github.com/bensallen/hkmgr/internal/down"
	"github.com/bensallen/hkmgr/internal/ssh"
	"github.com/bensallen/hkmgr/internal/up"
	"github.com/integrii/flaggy"
)

var version = "unknown"

func Run() error {
	var upSubcommand *flaggy.Subcommand
	var downSubcommand *flaggy.Subcommand
	var destroySubcommand *flaggy.Subcommand
	var validateSubcommand *flaggy.Subcommand
	var sshSubcommand *flaggy.Subcommand
	var consoleSubcommand *flaggy.Subcommand

	var configPath string

	flaggy.SetName("hkmgr")
	flaggy.SetDescription("VM manager for hyperkit")

	flaggy.DefaultParser.AdditionalHelpPrepend = "http://github.com/bensallen/hkmgr"

	flaggy.String(&configPath, "c", "config", "Path to configuration TOML file")

	upSubcommand = flaggy.NewSubcommand("up")
	upSubcommand.Description = "Start VMs"

	downSubcommand = flaggy.NewSubcommand("down")
	downSubcommand.Description = "Stop VMs"

	destroySubcommand = flaggy.NewSubcommand("destroy")
	destroySubcommand.Description = "Destroy VMs"

	validateSubcommand = flaggy.NewSubcommand("validate")
	validateSubcommand.Description = "Validate configuration"

	sshSubcommand = flaggy.NewSubcommand("ssh")
	sshSubcommand.Description = "SSH to VM"

	consoleSubcommand = flaggy.NewSubcommand("console")
	consoleSubcommand.Description = "Open Console of VM"

	flaggy.AttachSubcommand(upSubcommand, 1)
	flaggy.AttachSubcommand(downSubcommand, 1)
	flaggy.AttachSubcommand(destroySubcommand, 1)
	flaggy.AttachSubcommand(validateSubcommand, 1)
	flaggy.AttachSubcommand(sshSubcommand, 1)
	flaggy.AttachSubcommand(consoleSubcommand, 1)

	flaggy.SetVersion(version)
	flaggy.Parse()

	var config config.Config
	if _, err := toml.DecodeFile(configPath, &config); err != nil {
		return err
	}

	if upSubcommand.Used {
		if err := up.Run(&config); err != nil {
			return err
		}
	} else if downSubcommand.Used {
		if err := down.Run(&config); err != nil {
			return err
		}
	} else if destroySubcommand.Used {
		if err := destroy.Run(&config); err != nil {
			return err
		}
	} else if sshSubcommand.Used {
		if err := ssh.Run(&config); err != nil {
			return err
		}
	} else if consoleSubcommand.Used {
		if err := console.Run(&config); err != nil {
			return err
		}
	}

	return nil
}