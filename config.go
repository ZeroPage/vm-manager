package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strconv"
)

type Config struct {
	VM map[string]VmConfig
}

type VmConfig struct {
	CPU     int
	Memory  string
	Network NetworkConfig
	Logfile string
	Disk    DiskConfig //TODO multiple Disk
	CdRom   string
	VGA     string
	VNC     string
}

func (vm VmConfig) makeArgs() []string {
	var args []string

	if vm.CPU == 0 {
		fmt.Println("CPU must be >= 1")
		os.Exit(-2)
	}
	args = append(args, "-smp")
	args = append(args, strconv.Itoa(vm.CPU))

	if vm.Memory == "" {
		fmt.Println("Memory must be exist")
		os.Exit(-2)
	}
	args = append(args, "-m")
	args = append(args, vm.Memory)

	if vm.Disk.Path == "" {
		fmt.Println("Disk Must Exist")
		os.Exit(-2)
	}
	args = append(args, "-drive")
	args = append(args, vm.Disk.makeArgs()...)

	args = append(args, "-daemonize")

	if vm.Network.Ifname == "" || vm.Network.MAC == "" {
		fmt.Println("Network Must Exist")
		os.Exit(-2)
	}

	args = append(args, vm.Network.makeArgs()...)

	if vm.VGA != "" {
		args = append(args, "-vga")
		args = append(args, vm.VGA)
	}

	if vm.VNC != "" {
		args = append(args, "-vnc")
		args = append(args, vm.VNC)
	}

	return args
}

type NetworkConfig struct {
	MAC    string
	Ifname string
}

func (nc NetworkConfig) makeArgs() (args []string) {
	mac := "nic,macaddr=" + nc.MAC

	ifname := "tap,ifname=" + nc.Ifname
	ifname += ",script=/etc/kvm-ifup,downscript=/etc/kvm-ifdown"

	args = append(args, "-net", mac, "-net", ifname)
	return
}

type DiskConfig struct {
	Path      string
	Interface string
}

func (dc DiskConfig) makeArgs() (args []string) {
	var arg string

	arg += "file=" + dc.Path
	if dc.Interface != "" {
		arg += ",if=" + dc.Interface
	}
	args = append(args, arg)
	return
}

func readConfig(filename string) *Config {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Connot read config file :", filename)
		os.Exit(-100)
	}
	config := &Config{}
	err = yaml.Unmarshal(buf, config)
	if err != nil {
		fmt.Println("Connot parse config file :", filename)
		os.Exit(-100)
	}
	return config
}
