package config

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strconv"

	"github.com/google/uuid"
)

func pidFile(path string) (int, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return 0, fmt.Errorf("pid file not found")
	}
	pidTxt, err := ioutil.ReadFile(path)
	if err != nil {
		return 0, fmt.Errorf("pid file cannot be read")
	}

	pid, err := strconv.Atoi(string(pidTxt))
	if err != nil {
		return 0, fmt.Errorf("pid file does not contain an integer")
	}

	return pid, nil
}

func uuidFile(path string) (uuid.UUID, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return [16]byte{}, fmt.Errorf("uuid file not found")
	}
	uuidTxt, err := ioutil.ReadFile(path)
	if err != nil {
		return [16]byte{}, fmt.Errorf("uuid file cannot be read")
	}

	uuid, err := uuid.Parse(string(uuidTxt))
	if err != nil {
		return [16]byte{}, fmt.Errorf("uuid file does not contain an UUID")
	}

	return uuid, nil
}

func hwaddrFile(path string) (net.HardwareAddr, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("mac addr file not found")
	}
	hwaddrTxt, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("mac addr file cannot be read")
	}

	hwaddr, err := net.ParseMAC(string(hwaddrTxt))
	if err != nil {
		return nil, fmt.Errorf("mac addr file does not contain a hardware address")
	}

	return hwaddr, nil
}
