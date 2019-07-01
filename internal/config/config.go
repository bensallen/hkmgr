package config

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

// Boot config
type Boot struct {
	Kexec    Kexec    `toml:"kexec"`
	Firmware Firmware `toml:"firmware"`
}

type Kexec struct {
	Kernel  string `toml:"kernel"`
	Initrd  string `toml:"initrd"`
	Cmdline string `toml:"cmdline"`
}

type Firmware struct {
	Path string `toml:"path"`
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
