package network

import (
	"encoding/hex"
	"fmt"
	"os/exec"
	"strings"
)

//Bridge is a Darwin/BSD network bridge device
type Bridge struct {
	Device  string
	IP      string
	Netmask []byte
	Members []string
}

//Up brings the defined bridge interface up in an idempotent fashion
func (b *Bridge) Up() error {
	bridge, err := findBridge(b.Device)
	if err != nil {
		if err := b.create(); err != nil {
			return err
		}
		if err := b.setIP(); err != nil {
			return err
		}
		if err := b.setMembers(nil); err != nil {
			return err
		}
		if err := b.setUp(); err != nil {
			return err
		}
	} else {
		if b.IP != bridge.IP {
			if err := b.setIP(); err != nil {
				return err
			}
		}
		if err := b.setMembers(bridge.Members); err != nil {
			return err
		}
	}

	return nil
}

//findBridge runs ifconfig <device> and parses output into a Bridge
func findBridge(device string) (*Bridge, error) {

	cmd := exec.Command("ifconfig", device)
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	bridge := parseIfconfig(string(out))
	bridge.Device = device

	return bridge, nil
}

// parseIfconfig parses output from ifconfig returning a *Bridge with members,
// IP addr, and netmask populated if found.
func parseIfconfig(ifconfig string) *Bridge {
	bridge := Bridge{}

	for _, line := range strings.Split(ifconfig, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "member:") {
			kv := strings.Split(line, " ")
			bridge.Members = append(bridge.Members, kv[1])
		} else if strings.HasPrefix(line, "inet ") {
			kv := strings.Split(line, " ")
			bridge.IP = kv[1]
			netmaskHex := strings.TrimPrefix(kv[3], "0x")
			bridge.Netmask, _ = hex.DecodeString(netmaskHex)
		}
	}
	return &bridge
}

// create runs ifconfig bridge<N> create
func (b *Bridge) create() error {
	cmd := exec.Command("ifconfig", b.Device, "create")
	return cmd.Run()
}

// Destroy runs ifconfig bridge<N> destroy
func (b *Bridge) Destroy() error {
	cmd := exec.Command("ifconfig", b.Device, "destroy")
	return cmd.Run()
}

// setUp runs ifconfig bridge<N> up
func (b *Bridge) setUp() error {
	cmd := exec.Command("ifconfig", b.Device, "up")
	return cmd.Run()
}

// setIP runs ifconfig bridge<N> <ipAddr> netmask <netmask>. Note, netmask is
// passed in its hex form, eg. 0xffffff00. Similar to ifconfig's default output.
func (b *Bridge) setIP() error {
	cmd := exec.Command("ifconfig", b.Device, b.IP, "netmask", fmt.Sprintf("0x%x", b.Netmask))
	return cmd.Run()
}

// setMembers adds member devices to a bridge device
func (b *Bridge) setMembers(cur []string) error {
	add, del := sliceDiff(b.Members, cur)
	if add != nil {
		if err := addMembers(b.Device, add); err != nil {
			return err
		}
	}
	if del != nil {
		if err := delMembers(b.Device, del); err != nil {
			return err
		}
	}
	return nil
}

// addMembers added member devices to a bridge device by running ifconfig bridge<N> addm <dev1> addm <dev2> ...
func addMembers(device string, members []string) error {
	if len(members) == 0 {
		return nil
	}

	args := []string{device}
	for _, member := range members {
		args = append(args, "addm", member)
	}
	cmd := exec.Command("ifconfig", args...)
	return cmd.Run()
}

// delMembers deletes member devices to a bridge device by running ifconfig bridge<N> delm <dev1> delm <dev2> ...
func delMembers(device string, members []string) error {
	if len(members) == 0 {
		return nil
	}

	args := []string{device}
	for _, member := range members {
		args = append(args, "delm", member)
	}
	cmd := exec.Command("ifconfig", args...)
	return cmd.Run()
}

// sliceDiff compares two slices of strings. The first returned slice are
// elements in x that aren't represented in y. The second returned slice are
// elements in y that aren't represented in x. Duplicate elements are ignored.
func sliceDiff(x, y []string) ([]string, []string) {
	if len(x) == 0 {
		return nil, y
	}

	if len(y) == 0 {
		return x, nil
	}

	xMap := make(map[string]bool)
	yMap := make(map[string]bool)

	for _, xElem := range x {
		xMap[xElem] = true
	}
	for _, yElem := range y {
		yMap[yElem] = true
	}

	var diffX []string
	var diffY []string

	for xMapKey := range xMap {
		if _, ok := yMap[xMapKey]; !ok {
			diffX = append(diffX, xMapKey)
		}
	}
	for yMapKey := range yMap {
		if _, ok := xMap[yMapKey]; !ok {
			diffY = append(diffY, yMapKey)
		}
	}
	return diffX, diffY
}
