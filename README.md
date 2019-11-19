# Hkmgr

Hkmgr is a virtual machine manager that makes use of the [hyperkit](https://github.com/moby/hyperkit) hypervisor. Hkmgr helps simplify running hyperkit and automates some of the network configuration on macOS.

## Install

```shell
brew tap bensallen/homebrew-hkmgr
brew install hkmgr
```

- For tap network support, [tuntaposx](http://tuntaposx.sourceforge.net) is required.

## Example: Boot a RancherOS VM

```shell
mkdir ros-vm1
cd ros-vm1
wget https://github.com/rancher/os/releases/download/v1.5.4/initrd
wget https://github.com/rancher/os/releases/download/v1.5.4/vmlinuz
```

- Create a hkmgr.toml

```toml
[network.net1.vmnet]
bridge = "bridge100"
ip = "192.168.64.1/24"

[network.net2.tap]
bridge = "bridge1"
ip = "192.168.99.1/24"

[vm.ros-vm1]
memory = "4GB"
cores = 1
uuid = "9445CA7C-F976-456E-9061-B932194D8166"

[[vm.ros-vm1.network]]
driver = "virtio-net"
memberOf = "net1"

[[vm.ros-vm1.network]]
ip = "192.168.99.11/24"
mac = "AA:BB:CC:DD:00:01"
device = "tap0"
driver = "virtio-tap"
memberOf = "net2"

[vm.ros-vm1.boot.kexec]
kernel = "vmlinuz"
initrd = "initrd"
cmdline = "console=tty0 console=ttyS0,115200 earlyprintk=serial rancher.autologin=ttyS0 rancher.network.interfaces.eth1.address=192.168.99.11/24"
```

- Boot the VM

```shell
$ sudo hkmgr up
Password:
Booting VM: 9445CA7C-F976-456E-9061-B932194D8166
cmd: hyperkit -U 9445CA7C-F976-456E-9061-B932194D8166 -c 4 -m 4GB -A -s 0:0,hostbridge -s 31,lpc -s 1,virtio-rnd -l com1,autopty=/Users/ballen/Demo/.run/vm/ros-vm1/tty,asl -s 2:0,virtio-net,, -s 2:1,virtio-tap,tap0,mac=AA:BB:CC:DD:00:01 -f kexec,/Users/ballen/Demo/vmlinuz,/Users/ballen/Demo/initrd,console=tty0 console=ttyS0,115200 earlyprintk=serial rancher.autologin=ttyS0 rancher.network.interfaces.eth1.address=192.168.99.11/24
adding member tap0 to network net2 for vm 9445CA7C-F976-456E-9061-B932194D8166
Configuring Network: "net1"
Configuring Network: "net2"
cmd: ifconfig bridge1 addm tap0
cmd: ifconfig bridge1 up
```

- Access the serial console

```shell
sudo screen .run/vm/ros-vm1/tty
```

## Hkmgr CLI

```
$ hkmgr -h
hkmgr - VM manager for hyperkit
http://github.com/bensallen/hkmgr

  Usage:
    hkmgr [up|down|status]

  Subcommands:
    up - Start VMs
    down - Stop VMs
    status - Display status of VMs

  Flags:
       --version  Displays the program version string.
    -h --help  Displays help with available flag, subcommand, and positional value parameters.
    -c --config  Path to configuration TOML file
    -d --debug  Enable debug output
    -n --dry-run  Don't execute any commands that affect change, just show what will be run
```

```
$ hkmgr up -h
up - Start VMs

  Usage:
    up [name]

  Positional Variables:
    name - Specify a VM, otherwise all VMs will be run
  Flags:
       --version  Displays the program version string.
    -h --help  Displays help with available flag, subcommand, and positional value parameters.
    -c --config  Path to configuration TOML file
    -d --debug  Enable debug output
    -n --dry-run  Don't execute any commands that affect change, just show what will be run
```

```
$ hkmgr status -h
status - Display status of VMs

  Usage:
    status [name]

  Positional Variables:
    name - Specify a VM
  Flags:
       --version  Displays the program version string.
    -h --help  Displays help with available flag, subcommand, and positional value parameters.
    -c --config  Path to configuration TOML file
    -d --debug  Enable debug output
    -n --dry-run  Don't execute any commands that affect change, just show what will be run
```

```
$ hkmgr down -h
down - Stop VMs

  Usage:
    down [name]

  Positional Variables:
    name - Specify a VM, otherwise all VMs will be stopped
  Flags:
       --version  Displays the program version string.
    -h --help  Displays help with available flag, subcommand, and positional value parameters.
    -s --signal  Signal to send to VM
    -c --config  Path to configuration TOML file
    -d --debug  Enable debug output
    -n --dry-run  Don't execute any commands that affect change, just show what will be run
```

## Build

```shell
git clone https://github.com/bensallen/hkmgr.git
cd hkmgr
go build -o hkmgr cmd/hkmgr/main.go
```

## Contribute

Issues and PRs as welcome at https://github.com/bensallen/hkmgr.

Check GitHub issues or [TODO.md](TODO.md) for features that are outstanding.