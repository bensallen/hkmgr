package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/mitchellh/go-ps"
)

var hyperkitPath = "hyperkit"

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

// Status is the status of a VM process
type Status int

const (
	// Running status represents a running hyperkit process being found
	Running Status = 1
	// Stopped status represents that a pid file was found but a matching hyperkit process was not found
	Stopped Status = 2
	// NotFound status represents that a pid file was not found
	NotFound Status = 3
)

// Up starts a VM if its not already running.
func (v *VMConfig) Up() error {
	if v.Status() == Running {
		return nil
	}

	cmdArgs := v.Cli()

	fmt.Printf("cmd: %s %s\n", hyperkitPath, strings.Join(cmdArgs, " "))

	cmd := exec.Command(hyperkitPath, cmdArgs...)
	err := cmd.Start()
	if err != nil {
		return err
	}

	w, err := os.Create(v.RunDir + "/pid")
	defer w.Close()
	if err != nil {
		return err
	}

	if _, err := w.WriteString(strconv.Itoa(cmd.Process.Pid)); err != nil {
		return err
	}

	return nil
}

// Status attempts to find the pid file in the run dir of a VM and checks to see if its running or not.
func (v *VMConfig) Status() Status {
	pidFile := v.RunDir + "/pid"
	if _, err := os.Stat(pidFile); os.IsNotExist(err) {
		return NotFound
	}
	pidTxt, err := ioutil.ReadFile(pidFile)
	if err != nil {
		return NotFound
	}

	pid, err := strconv.Atoi(string(pidTxt))
	if err != nil {
		return NotFound
	}

	proc, err := ps.FindProcess(pid)
	if err != nil {
		return NotFound
	}
	if proc == nil {
		return Stopped
	} else {
		return Running
	}
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
