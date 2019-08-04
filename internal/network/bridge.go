package network

import (
	"encoding/hex"
	"fmt"
	"net"
	"os/exec"
	"strings"
	"time"
)

//Bridge is a Darwin/BSD network bridge device
type Bridge struct {
	Device  string
	IP      net.IP
	Netmask net.IPMask
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
		if err := b.setMembers(nil, false); err != nil {
			return err
		}
		if err := b.setUp(); err != nil {
			return err
		}
	} else {
		if b.IP.String() != bridge.IP.String() {
			if err := b.setIP(); err != nil {
				return err
			}
		}
		if b.Netmask.String() != bridge.Netmask.String() {
			if err := b.setIP(); err != nil {
				return err
			}
		}
		if err := b.setMembers(bridge.Members, false); err != nil {
			return err
		}
		if err := b.setUp(); err != nil {
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
			bridge.IP = net.ParseIP(kv[1])
			netmaskHex := strings.TrimPrefix(kv[3], "0x")
			bridge.Netmask, _ = hex.DecodeString(netmaskHex)
		}
	}
	return &bridge
}

// create runs ifconfig bridge<N> create
func (b *Bridge) create() error {
	cmd := exec.Command("ifconfig", b.Device, "create")
	fmt.Printf("cmd: %s\n", strings.Join(cmd.Args, " "))
	return cmd.Run()
}

// Destroy runs ifconfig bridge<N> destroy
func (b *Bridge) Destroy() error {
	cmd := exec.Command("ifconfig", b.Device, "destroy")
	fmt.Printf("cmd: %s\n", strings.Join(cmd.Args, " "))
	return cmd.Run()
}

// setUp runs ifconfig bridge<N> up
func (b *Bridge) setUp() error {
	cmd := exec.Command("ifconfig", b.Device, "up")
	fmt.Printf("cmd: %s\n", strings.Join(cmd.Args, " "))
	return cmd.Run()
}

// setIP runs ifconfig bridge<N> <ipAddr> netmask <netmask>. Note, netmask is
// passed in its hex form, eg. 0xffffff00. Similar to ifconfig's default output.
func (b *Bridge) setIP() error {
	if b.IP == nil {
		return nil
	}
	cmd := exec.Command("ifconfig", b.Device, b.IP.String(), "netmask", fmt.Sprintf("0x%s", b.Netmask.String()))
	fmt.Printf("cmd: %s\n", strings.Join(cmd.Args, " "))
	return cmd.Run()
}

// setMembers idempotently adds member devices listed in the Members attribute to a bridge device
// and conditionally removes any additional members passed in cur. A list of current bridge members
// should be passsed as an agrument for idempotence.
func (b *Bridge) setMembers(cur []string, delete bool) error {
	add, del := sliceDiff(b.Members, cur)
	if add != nil {
		if err := addMembers(b.Device, add); err != nil {
			return err
		}
	}
	if delete && del != nil {
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

		// Ugly polling and timeout mechanism waiting for the hyperkit to bring up the tap interface.
		var count int
		for {
			if count > 10 {
				return fmt.Errorf("failed to find device %s after 10 seconds, attempting to add it to bridge %s", member, device)
			}
			if _, err := net.InterfaceByName(member); err != nil {
				time.Sleep(1 * time.Second)
				count++
				continue
			}
			break
		}
		args = append(args, "addm", member)
	}
	cmd := exec.Command("ifconfig", args...)
	fmt.Printf("cmd: %s\n", strings.Join(cmd.Args, " "))
	return cmd.Run()
}

// delMembers deletes member devices to a bridge device by running ifconfig bridge<N> deletem <dev1> deletem <dev2> ...
func delMembers(device string, members []string) error {
	if len(members) == 0 {
		return nil
	}

	args := []string{device}
	for _, member := range members {
		args = append(args, "deletem", member)
	}
	cmd := exec.Command("ifconfig", args...)
	fmt.Printf("cmd: %s\n", strings.Join(cmd.Args, " "))
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
