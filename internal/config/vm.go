package config

import (
	"fmt"
	"strconv"
)

// VMs
type VM map[string]VMConfig

type VMConfig struct {
	Memory        string    `toml:"memory"`
	Cores         int       `toml:"cores"`
	UUID          string    `toml:"uuid"`
	SSHKey        string    `toml:"ssh_key"`
	ProvisionPre  string    `toml:"provision_pre"`
	ProvisionPost string    `toml:"provision_post"`
	Before        []string  `toml:"before"`
	After         []string  `toml:"after"`
	Requires      []string  `toml:"requires"`
	RunDir        string    `toml:"run_dir"`
	Network       []NetConf `toml:"network"`
	Boot          Boot      `toml:"boot"`
	HDD           []HDD     `toml:"hdd"`
	CDROM         []CDROM   `toml:"cdrom"`
}

func (v *VMConfig) Cli() []string {

	var args []string

	if v.UUID != "" {
		args = append(args, "-U", v.UUID)
	}

	if v.Cores != 0 {
		args = append(args, "-c", strconv.Itoa(v.Cores))
	}

	if v.Memory != "" {
		args = append(args, "-m", v.Memory)
	}

	args = append(args, "-A", "-s", "0:0,hostbridge", "-s", "31,lpc", "-s", "1,virtio-rnd")

	if v.RunDir != "" {
		args = append(args, "-l", fmt.Sprintf("com1,autopty=%s/tty,asl", v.RunDir))
	}

	for i, net := range v.Network {
		var opts string
		if net.MAC != "" {
			opts = fmt.Sprintf("mac=%s", net.MAC)
		}
		args = append(args, "-s", fmt.Sprintf("2:%d,%s,%s,%s", i, net.Driver, net.Device, opts))
	}

	for i, hdd := range v.HDD {
		args = append(args, "-s", fmt.Sprintf("3:%d,%s,file://%s,format=%s", i, hdd.Driver, hdd.Path, hdd.Format))
	}

	for i, cd := range v.CDROM {
		args = append(args, "-s", fmt.Sprintf("4:%d,%s,%s", i, cd.Driver, cd.Path))
	}

	args = append(args, v.Boot.Cli()...)

	return args
}

// Boot config
type Boot struct {
	Kexec    Kexec    `toml:"kexec"`
	Firmware Firmware `toml:"firmware"`
	FBSD     FBSD     `toml:"fbsd"`
}

func (b *Boot) Cli() []string {
	if (Kexec{}) != b.Kexec {
		return b.Kexec.Cli()
	}

	if (Firmware{}) != b.Firmware {
		return b.Firmware.Cli()
	}

	if (FBSD{}) != b.FBSD {
		return b.FBSD.Cli()
	}

	return []string{}
}

type Kexec struct {
	Kernel  string `toml:"kernel"`
	Initrd  string `toml:"initrd"`
	Cmdline string `toml:"cmdline"`
}

func (k *Kexec) Cli() []string {
	return []string{"-f", fmt.Sprintf("kexec,%s,%s,%s", k.Kernel, k.Initrd, k.Cmdline)}
}

type Firmware struct {
	Path string `toml:"path"`
}

func (f *Firmware) Cli() []string {
	return []string{"-f", fmt.Sprintf("bootrom,%s,,", f.Path)}
}

type FBSD struct {
	Userboot   string `toml:"userboot"`
	BootVolume string `toml:"userboot"`
	KernelEnv  string `toml:"kernelenv"`
}

func (f *FBSD) Cli() []string {
	return []string{"-f", fmt.Sprintf("fbsd,%s,%s,%s", f.Userboot, f.BootVolume, f.KernelEnv)}
}

// VM Network Config
type NetConf struct {
	IP       string `toml:"ip"`
	MAC      string `toml:"mac"`
	Device   string `toml:"device"`
	Driver   string `toml:"driver"`
	MemberOf string `toml:"memberOf"`
}

type HDD struct {
	Path   string
	Format string
	Driver string
	Size   string
	Create bool
}

func (h *HDD) create() error {
	return nil
}

type CDROM struct {
	Path    string
	Driver  string
	Extract bool
}
