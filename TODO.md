# TODOs

## Features

- Add CI, CircleCI, clean code, etc.
- Add dependency graph ordering

### Config

- Reorganize internal/config to move VM action logic outside of config, eg. vm/vm.go
- Generate HDD if it doesn't already exist
- Generate UUID if not specified, store in .run/vm/<name>/uuid
- Generate MAC for tap interfaces if not specified, store in .run/vm/<name>/<net>_mac
- Check write privs on run_dir, pid, and tty in Validate()
- Check write privs on hdd, read privs on cdrom
- Add template support for kexec cmdline for IP, ssh public
- Add hyperkit multiboot as another boot option

### Logging

- Remove all fmt.Printf and use log or similar
- Add debug logging
- Add cmd stdout/stderr to debug logging

### Network

- Add VPNKit support

### Up

- Change permissions on tap and tty, then drop privs to run hyperkit
- Automatically pick a unused tap interface if not specified

### Init

- Generate a hkmgr.toml

### Down

- Prompt when running without a specific VM, asking if all VMs should be stopped to avoid annoyances. Add -y flag to answer via CLI.

### Destroy

### Status

- Add pid in status output

### SSH

### Console

- Test and consider integrating https://github.com/ishuah/bifrost/, otherwise just exec screen

### Host Config Automation

- Support for adding/removing routes on the host
- Move pf rules to a new host sections of config
- Support for adding/removing pf rules on the host
- Add sysctl enable forwarding

## Bugs

- Fix double call to setIP

Configuring Network: "net0"
cmd: ifconfig bridge1 192.168.99.1 netmask 0xffffff00
cmd: ifconfig bridge1 192.168.99.1 netmask 0xffffff00


## Completed

- Add arg to each relevant command to work on specific VM
- Default find hkmgr.toml in current path if not specified
- Change addm/deletem of bridge members when working with a single VM to only addm, so we don't superfluously deletem members of VMs not being considered.
- Check if the config file exists
- Check UID == 0 if using vmnet
- Check status of a VM before starting in Up()
- Check that kernel and vmlinux exist in Up() of boot.kexec
- If one of the tap interfaces doesn't come up, still add the ones that do come up to the bridge
- Make relative paths be relative to the config file.
- Make relative paths in a VM config relative to the run_dir
