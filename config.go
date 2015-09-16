package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

type Config struct {
	VM      map[string]VmConfig
	LogPath string
	Pidfile PidConfig
}
type PidConfig struct {
	Path   string
	Prefix string
	Suffix string
}

func (pid PidConfig) getPath(name string) string {
	return path.Join(pid.Path, pid.Prefix+name+pid.Suffix)
}

func NewConfig() *Config {
	return &Config{
		Pidfile: PidConfig{
			Path:   "/tmp",
			Prefix: "vm.",
			Suffix: ".pid",
		},
	}
}

func (config Config) getConfigArgs(name string) []string {
	vmconfig, ok := config.VM[name]
	if !ok {
		L.ERR("config not exist!")
		os.Exit(-1)
	}
	args := vmconfig.makeArgs()

	args = append(args, "-pidfile", config.Pidfile.getPath(name))
	return args
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

	args = append(args, "-daemonize")

	args = append(args, vm.CPU.makeArgs()...)

	if vm.Memory == "" {
		L.ERR("Memory must be exist")
		os.Exit(-2)
	}
	args = append(args, "-m", vm.Memory)

	for _, disk := range vm.Disks {
		args = append(args, disk.makeArgs()...)
	}

	for _, cdrom := range vm.CdRoms {
		args = append(args, cdrom.makeArgs()...)
	}

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
		L.ERR("CPU must be >= 1")
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
		L.WARN("Disk must have path - skip disk")
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
		L.ERR("Connot read config file :", filename, err)
		os.Exit(-100)
	}
	config := NewConfig()
	err = yaml.Unmarshal(buf, config)
	if err != nil {
		L.ERR("Connot parse config file :", err)
		os.Exit(-100)
	}
	return config
}
