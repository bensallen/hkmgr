package config

import (
	"fmt"
	"net"
	"strings"

	"github.com/bensallen/hkmgr/internal/network"
)

type Network map[string]NetTypes

type NetTypes struct {
	Vmnet  *Vmnet  `toml:"vmnet"`
	Tap    *Tap    `toml:"tap"`
	VPNKit *VPNKit `toml:"vpnkit"`
}

//NetType asdf
func (nt *NetTypes) NetType() NetType {
	if nt.Vmnet != nil {
		return nt.Vmnet
	}
	if nt.Tap != nil {
		return nt.Tap
	}
	if nt.VPNKit != nil {
		return nt.VPNKit
	}
	return nil
}

type NetType interface {
	Up() error
	Destroy() error
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

	bridge, err := t.toBridge()
	if err != nil {
		return err
	}

	return bridge.Up()
}

func (t *Tap) toBridge() (*network.Bridge, error) {
	bridge := network.Bridge{}
	if t.IP != "" {
		var ip net.IP
		var mask []byte
		if strings.Contains(t.IP, "/") {
			cidr := &net.IPNet{}
			var err error
			ip, cidr, err = net.ParseCIDR(t.IP)
			if err != nil {
				return nil, err
			}
			mask = cidr.Mask
		} else {
			ip = net.ParseIP(t.IP)
			if ip == nil {
				return nil, fmt.Errorf("Could not parse IP address: %s", t.IP)
			}
			mask = ip.DefaultMask()
		}
		bridge = network.Bridge{Device: t.Bridge, IP: ip, Netmask: mask}
	} else {
		bridge = network.Bridge{Device: t.Bridge}
	}

	return &bridge, nil
}

func (t *Tap) Destroy() error {
	return nil
}

type VPNKit struct {
}

func (v *VPNKit) Up() error {
	return nil
}

func (v *VPNKit) Destroy() error {
	return nil
}
