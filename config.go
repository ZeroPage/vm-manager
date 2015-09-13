package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strconv"
)

type Config struct {
	VM      map[string]VmConfig
	LogPath string
}

type VmConfig struct {
	CPU     CPUConfig
	Memory  string
	Network NetworkConfig
	LogFile string
	Disks   []DiskConfig
	CdRoms  []CDRomConfig
	VGA     string
	VNC     string
	Others  []string
}

func (vm VmConfig) makeArgs() []string {
	var args []string

	args = append(args, vm.CPU.makeArgs()...)

	if vm.Memory == "" {
		fmt.Println("Memory must be exist")
		os.Exit(-2)
	}
	args = append(args, "-m", vm.Memory)

	for _, disk := range vm.Disks {
		args = append(args, disk.makeArgs()...)
	}

	for _, cdrom := range vm.CdRoms {
		args = append(args, cdrom.makeArgs()...)
	}
	args = append(args, "-daemonize")

	if vm.Network.Ifname != "" || vm.Network.MAC != "" {
		args = append(args, vm.Network.makeArgs()...)
	}

	if vm.VGA != "" {
		args = append(args, "-vga", vm.VGA)
	}

	if vm.VNC != "" {
		args = append(args, "-vnc", vm.VNC)
	}

	args = append(args, vm.Others...)

	return args
}

type CPUConfig struct {
	Type        string
	Core        int
	Sockets     int
	VirtualCore int
	Threads     int
}

func (cc CPUConfig) makeArgs() (args []string) {
	if cc.Type != "" {
		args = append(args, "-cpu", cc.Type)
	}

	if cc.Core <= 0 {
		fmt.Println("CPU must be >= 1")
		os.Exit(-2)
	}

	var smp string
	smp += strconv.Itoa(cc.Core)

	if cc.Sockets != 0 {
		smp += ",sockets=" + strconv.Itoa(cc.Sockets)
	}
	if cc.VirtualCore != 0 {
		smp += ",cores=" + strconv.Itoa(cc.VirtualCore)
	}
	if cc.Threads != 0 {
		smp += ",threads=" + strconv.Itoa(cc.Threads)
	}

	args = append(args, "-smp", smp)
	return
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
	Index     int
}

func (dc DiskConfig) makeArgs() (args []string) {
	var drive string
	if dc.Path == "" {
		fmt.Println("Disk must have path")
		return
	}
	drive += "file=" + dc.Path

	if dc.Interface != "" {
		drive += ",if=" + dc.Interface
	}
	if dc.Index != 0 {
		drive += ",index=" + strconv.Itoa(dc.Index)
	}
	drive += ",media=disk"
	args = append(args, "-drive", drive)
	return
}

type CDRomConfig string

func (cdrom CDRomConfig) makeArgs() (args []string) {
	var drive string
	drive += "file=" + string(cdrom)
	drive += ",media=cdrom"
	drive += ",if=ide"
	args = append(args, "-drive", drive)
	return
}

func readConfig(filename string) *Config {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Connot read config file :", filename, err)
		os.Exit(-100)
	}
	config := &Config{}
	err = yaml.Unmarshal(buf, config)
	if err != nil {
		fmt.Println("Connot parse config file :", err)
		os.Exit(-100)
	}
	return config
}
