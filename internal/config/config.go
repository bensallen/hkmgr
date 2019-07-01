package config

import "fmt"

type Config struct {
	Network Network `toml:"network"`
	VM      VM      `toml:"vm"`
}

type Network map[string]NetTypes

type NetTypes struct {
	Vmnet  Vmnet  `toml:"vmnet"`
	Tap    Tap    `toml:"tap"`
	VPNKit VPNKit `toml:"vpnkit"`
}

type Vmnet struct {
	Bridge string `toml:"bridge"`
	IP     string `toml:"ip"`
}

type Tap struct {
	Bridge  string   `toml:"bridge"`
	IP      string   `toml:"ip"`
	Nat     bool     `toml:"nat"`
	NatIf   string   `toml:"nat_if"`
	PfRules []string `toml:"pf_rules"`
	DHCP    bool     `tool:"dhcp"`
}

type VPNKit struct {
}

// VMs
type VM map[string]VMConfig

type VMConfig struct {
	Memory        string             `toml:"memory"`
	Cores         int                `toml:"cores"`
	UUID          string             `toml:"uuid"`
	SSHKey        string             `toml:"ssh_key"`
	ProvisionPre  string             `toml:"provision_pre"`
	ProvisionPost string             `toml:"provision_post"`
	Before        []string           `toml:"before"`
	After         []string           `toml:"after"`
	Requires      []string           `toml:"requires"`
	RunDir        string             `toml:"run_dir"`
	Network       map[string]NetConf `toml:"network"`
	Boot          Boot               `toml:"boot"`
	HDD           []HDD              `toml:"hdd"`
	CDROM         []CDROM            `toml:"cdrom"`
}

func (v *VMConfig) Cli() []string {
	return v.Boot.Cli()
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
	IP  string `toml:"ip"`
	MAC string `toml:"mac"`
}

type HDD struct {
	Path   string
	Type   string
	Size   string
	Create bool
}

type CDROM struct {
	Path    string
	Extract bool
}
