package config

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

func (v *Vmnet) Up() error {
	return nil
}

func (v *Vmnet) Destroy() error {
	return nil
}

type Tap struct {
	Bridge  string   `toml:"bridge"`
	IP      string   `toml:"ip"`
	Nat     bool     `toml:"nat"`
	NatIf   string   `toml:"nat_if"`
	PfRules []string `toml:"pf_rules"`
	DHCP    bool     `tool:"dhcp"`
}

func (t *Tap) Up() error {
	return nil
}

func (t *Tap) Destroy() error {
	return nil
}

type VPNKit struct {
}
