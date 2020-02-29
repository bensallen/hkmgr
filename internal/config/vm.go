package config

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"github.com/google/uuid"
	"github.com/mitchellh/go-ps"
)

var hyperkitPath = "hyperkit"

// VM is the map of VM names to VMConfig
type VM map[string]*VMConfig

type VMConfig struct {
	Memory        string     `toml:"memory"`
	Cores         int        `toml:"cores"`
	UUID          string     `toml:"uuid"`
	SSHKey        string     `toml:"ssh_key"`
	ProvisionPre  string     `toml:"provision_pre"`
	ProvisionPost string     `toml:"provision_post"`
	Before        []string   `toml:"before"`
	After         []string   `toml:"after"`
	Requires      []string   `toml:"requires"`
	RunDir        string     `toml:"run_dir"`
	Network       []*NetConf `toml:"network"`
	Boot          Boot       `toml:"boot"`
	HDD           []*HDD     `toml:"hdd"`
	CDROM         []*CDROM   `toml:"cdrom"`
	PID           int
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

func (s Status) String() string {
	switch s {
	case 1:
		return "running"
	case 2:
		return "stopped"
	case 3:
		return "PID file not found"
	default:
		return "unknown"
	}
}

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
	if err != nil {
		return err
	}
	defer w.Close()

	if _, err := w.WriteString(strconv.Itoa(cmd.Process.Pid)); err != nil {
		return err
	}

	return nil
}

// Down stops a VM if its running.
func (v *VMConfig) Down(signal string) error {
	if v.Status() != Running {
		return fmt.Errorf("not running")
	}

	return v.Kill(signal)
}

// Status attempts to find the pid file in the run dir of a VM and checks to see if its running or not.
func (v *VMConfig) Status() Status {
	pidFilePath := v.RunDir + "/pid"
	pid, err := pidFile(pidFilePath)
	if err != nil {
		return NotFound
	}

	v.PID = pid

	proc, err := ps.FindProcess(pid)
	if err != nil {
		return NotFound
	}
	if proc == nil {
		return Stopped
	}
	if proc.Executable() == "hyperkit" || proc.Executable() == "com.docker.hyper" {
		return Running
	}
	return Stopped
}

// Kill attempts to kill a VM via the pid file with the specified signal.
func (v *VMConfig) Kill(signal string) error {
	var sysSig syscall.Signal
	if signal == "" {
		sysSig = syscall.SIGTERM
	} else {
		var err error
		if sysSig, err = sigLookup(signal); err != nil {
			return err
		}
	}

	pidFilePath := v.RunDir + "/pid"
	pid, err := pidFile(pidFilePath)
	if err != nil {
		return err
	}

	proc, err := os.FindProcess(pid)
	if err != nil {
		return err
	}

	return proc.Signal(sysSig)
}

func sigLookup(s string) (syscall.Signal, error) {
	sigmap := map[string]syscall.Signal{
		"SIGINT":  syscall.Signal(2),
		"SIGKILL": syscall.Signal(9),
		"SIGUSR1": syscall.Signal(10),
		"SIGUSR2": syscall.Signal(12),
		"SIGTERM": syscall.Signal(15),
		"2":       syscall.Signal(2),
		"9":       syscall.Signal(9),
		"10":      syscall.Signal(10),
		"12":      syscall.Signal(12),
		"15":      syscall.Signal(15),
	}

	sig, ok := sigmap[s]
	if !ok {
		return 0, fmt.Errorf("%s is not a supported signal", s)
	}

	return sig, nil
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
		args = append(args, "-l", fmt.Sprintf("com1,autopty=%s/tty,log=%s/log", v.RunDir, v.RunDir))
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

func (v *VMConfig) Validate() error {
	if v.UUID == "" {
		return errors.New("UUID not specified")
	}

	if v.Cores == 0 {
		return errors.New("cores not specified")
	}

	if v.Memory == "" {
		return errors.New("memory not specified")
	}

	if v.RunDir == "" {
		return errors.New("RunDir not specified")
	}

	if err := v.Boot.validate(); err != nil {
		return err
	}

	// Return here if the VM is already running, as the remaining checks will fail with a running VM.
	if v.Status() == Running {
		return nil
	}

	for _, net := range v.Network {
		if err := net.validate(); err != nil {
			return err
		}
	}

	return nil
}

func (v *VMConfig) defaults(configDir string, name string) error {
	if v.RunDir == "" {
		v.RunDir = filepath.Join(configDir, ".run/vm/", name)
	}

	if v.UUID == "" {
		UUID, err := uuidFile(v.RunDir + "/uuid")

		if err != nil {
			UUID = uuid.New()
			w, err := os.Create(v.RunDir + "/uuid")
			if err != nil {
				return err
			}

			defer w.Close()

			if _, err := w.WriteString(UUID.String()); err != nil {
				return err
			}
		}
		v.UUID = UUID.String()
	}

	for _, net := range v.Network {
		if err := net.defaults(v.RunDir); err != nil {
			return err
		}
	}

	return nil
}

func (v *VMConfig) updateRelativePaths(configDir string, name string) {
	if v.RunDir[:1] != "/" {
		v.RunDir = filepath.Join(configDir, v.RunDir)
	}

	v.Boot.updateRelativePaths(configDir)
	for i := range v.HDD {
		v.HDD[i].updateRelativePaths(v.RunDir)
	}
	for i := range v.CDROM {
		v.CDROM[i].updateRelativePaths(v.RunDir)
	}
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

func (b *Boot) validate() error {
	if (Kexec{}) != b.Kexec {
		return b.Kexec.Validate()
	}

	if (Firmware{}) != b.Firmware {
		return b.Firmware.Validate()
	}

	if (FBSD{}) != b.FBSD {
		return b.FBSD.Validate()
	}

	return nil
}

func (b *Boot) updateRelativePaths(configDir string) {
	if (Kexec{}) != b.Kexec {
		b.Kexec.updateRelativePaths(configDir)
	}

	if (Firmware{}) != b.Firmware {
		b.Firmware.updateRelativePaths(configDir)
	}

	if (FBSD{}) != b.FBSD {
		b.FBSD.updateRelativePaths(configDir)
	}
}

type Kexec struct {
	Kernel  string `toml:"kernel"`
	Initrd  string `toml:"initrd"`
	Cmdline string `toml:"cmdline"`
}

func (k *Kexec) Cli() []string {
	return []string{"-f", fmt.Sprintf("kexec,%s,%s,%s", k.Kernel, k.Initrd, k.Cmdline)}
}

func (k *Kexec) Validate() error {
	if !fileExists(k.Kernel) {
		return fmt.Errorf("kernel not found: %s", k.Kernel)
	}

	if !fileExists(k.Initrd) {
		return fmt.Errorf("initrd not found: %s", k.Initrd)
	}

	return nil
}

func (k *Kexec) updateRelativePaths(configDir string) {
	if k.Kernel[:1] != "/" {
		k.Kernel = filepath.Join(configDir, k.Kernel)
	}

	if k.Initrd[:1] != "/" {
		k.Initrd = filepath.Join(configDir, k.Initrd)
	}
}

type Firmware struct {
	Path string `toml:"path"`
}

func (f *Firmware) Cli() []string {
	return []string{"-f", fmt.Sprintf("bootrom,%s,,", f.Path)}
}

//Validate for Firmware is currently a noop, TODO.
func (f *Firmware) Validate() error {
	return nil
}

func (f *Firmware) updateRelativePaths(configDir string) {
	if f.Path[:1] != "/" {
		f.Path = filepath.Join(configDir, f.Path)
	}
}

type FBSD struct {
	Userboot   string `toml:"userboot"`
	BootVolume string `toml:"userboot"`
	KernelEnv  string `toml:"kernelenv"`
}

func (f *FBSD) Cli() []string {
	return []string{"-f", fmt.Sprintf("fbsd,%s,%s,%s", f.Userboot, f.BootVolume, f.KernelEnv)}
}

//Validate for FBSD is currently a noop, TODO.
func (f *FBSD) Validate() error {
	return nil
}

func (f *FBSD) updateRelativePaths(configDir string) {
	if f.Userboot[:1] != "/" {
		f.Userboot = filepath.Join(configDir, f.Userboot)
	}

	if f.Userboot[:1] != "/" {
		f.BootVolume = filepath.Join(configDir, f.BootVolume)
	}
}

// NetConf is a VM network configuration
type NetConf struct {
	IP       string `toml:"ip"`
	MAC      string `toml:"mac"`
	Device   string `toml:"device"`
	Driver   string `toml:"driver"`
	MemberOf string `toml:"memberOf"`
}

func (n *NetConf) validate() error {

	switch n.Driver {

	case "virtio-tap":
		if n.MAC == "" {
			return errors.New("interface type tap requires a MAC address to be specified")
		}
		if n.Device == "" {
			return errors.New("interface type tap requires a Device to be specified")
		}
		dev, err := os.OpenFile(n.devicePath(), os.O_WRONLY, 0666)
		dev.Close()

		if err != nil {
			return fmt.Errorf("cannot open or write to tap interface %s, %v", n.Device, err)
		}

	case "virtio-net":
		if os.Geteuid() != 0 {
			return errors.New("virtio-net requires running as UID=0")
		}

	case "virtio-vpnkit":
		return errors.New("virtio-vpnkit support is not yet implemented")

	default:
		return fmt.Errorf("network driver %s not supported: drivers virtio-tap, virtio-net, and virtio-vpnkit are supported", n.Driver)
	}
	return nil
}

func (n *NetConf) defaults(runDir string) error {
	switch n.Driver {

	case "virtio-tap":
		if n.MAC == "" {
			MAC, err := hwaddrFile(runDir + "/" + n.MemberOf + "_mac")

			if err != nil {
				MAC, err = genMAC()
				if err != nil {
					return err
				}

				w, err := os.Create(runDir + "/" + n.MemberOf + "_mac")
				if err != nil {
					return err
				}

				defer w.Close()

				if _, err := w.WriteString(MAC.String()); err != nil {
					return err
				}
			}
			n.MAC = MAC.String()
		}
	}
	return nil
}

func (n *NetConf) devicePath() string {
	if n.Device[:1] == "/" {
		return n.Device
	}
	return "/dev/" + n.Device
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

func (h *HDD) updateRelativePaths(runDir string) {
	if h.Path[:1] != "/" {
		h.Path = filepath.Join(runDir, h.Path)
	}
}

type CDROM struct {
	Path    string
	Driver  string
	Extract bool
}

func (c *CDROM) updateRelativePaths(runDir string) {
	if c.Path[:1] != "/" {
		c.Path = filepath.Join(runDir, c.Path)
	}
}

// Check if a file exists and isn't a directory
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
