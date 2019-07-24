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
	Discover() error
	Up() error
	Destroy() error
}

type Vmnet struct {
	Bridge string `toml:"bridge"`
	IP     string `toml:"ip"`
}

func (v *Vmnet) Discover() error {
	return nil
}

func (v *Vmnet) Up() error {
	return nil
}

func (v *Vmnet) Destroy() error {
	return nil
}

type Tap struct {
	Bridge    string   `toml:"bridge"`
	IP        string   `toml:"ip"`
	Nat       bool     `toml:"nat"`
	NatIf     string   `toml:"nat_if"`
	PfRules   []string `toml:"pf_rules"`
	DHCP      bool     `toml:"dhcp"`
	BridgeDev *network.Bridge
}

func (t *Tap) Discover() error {
	if t.BridgeDev == nil {
		err := t.toBridge()
		if err != nil {
			return err
		}
	}
	return nil
}

func (t *Tap) Up() error {
	if err := t.Discover(); err != nil {
		return err
	}

	return t.BridgeDev.Up()
}

func (t *Tap) toBridge() error {
	bridge := network.Bridge{}
	if t.IP != "" {
		var ip net.IP
		var mask []byte
		if strings.Contains(t.IP, "/") {
			cidr := &net.IPNet{}
			var err error
			ip, cidr, err = net.ParseCIDR(t.IP)
			if err != nil {
				return err
			}
			mask = cidr.Mask
		} else {
			ip = net.ParseIP(t.IP)
			if ip == nil {
				return fmt.Errorf("Could not parse IP address: %s", t.IP)
			}
			mask = ip.DefaultMask()
		}
		bridge = network.Bridge{Device: t.Bridge, IP: ip, Netmask: mask}
	} else {
		bridge = network.Bridge{Device: t.Bridge}
	}

	t.BridgeDev = &bridge

	return nil
}

func (t *Tap) Destroy() error {
	return nil
}

type VPNKit struct {
}

func (v *VPNKit) Discover() error {
	return nil
}

func (v *VPNKit) Up() error {
	return nil
}

func (v *VPNKit) Destroy() error {
	return nil
}
